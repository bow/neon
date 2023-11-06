// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

func (s *SQLite) PullFeeds(ctx context.Context, ids []ID) <-chan PullResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		fail = failF("SQLite.PullFeeds")
		c    = make(chan PullResult)
		wg   sync.WaitGroup
	)
	ids = dedup(ids)

	// nolint: unparam
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		defer wg.Done()

		var (
			pks []pullKey
			err error
		)
		if len(ids) == 0 {
			pks, err = getAllPullKeys(ctx, tx)
		} else {
			pks, err = getPullKeys(ctx, tx, ids)
		}
		if err != nil {
			c <- NewPullResultFromError(nil, fail(err))
			return nil
		}
		if len(pks) == 0 {
			c <- NewPullResultFromFeed(nil, nil)
			return nil
		}

		chs := make([]<-chan PullResult, len(pks))
		for i, pk := range pks {
			chs[i] = pullNewFeedEntries(ctx, tx, pk, s.parser)
		}

		for pr := range merge(chs) {
			pr := pr
			if pr.Error() != nil {
				pr.err = fail(pr.err)
			}
			c <- pr
		}

		return nil
	}

	go func() {
		defer func() {
			wg.Wait()
			close(c)
		}()
		wg.Add(1)
		err := s.withTx(ctx, dbFunc)
		if err != nil {
			c <- NewPullResultFromError(nil, fail(err))
		}
	}()

	return c
}

// PullResult is a container for a pull operation.
type PullResult struct {
	url    *string
	status pullStatus
	feed   *FeedRecord
	err    error
}

func NewPullResultFromFeed(url *string, feed *FeedRecord) PullResult {
	return PullResult{status: pullSuccess, url: url, feed: feed}
}

func NewPullResultFromError(url *string, err error) PullResult {
	return PullResult{status: pullFail, url: url, err: err}
}

func (msg PullResult) Feed() *FeedRecord {
	if msg.status == pullSuccess {
		return msg.feed
	}
	return nil
}

func (msg PullResult) Error() error {
	if msg.status == pullFail {
		return msg.err
	}
	return nil
}

func (msg PullResult) URL() string {
	if msg.url != nil {
		return *msg.url
	}
	return ""
}

type pullStatus int

const (
	pullSuccess pullStatus = iota
	pullFail
)

type pullKey struct {
	feedID  ID
	feedURL string
}

func (pk pullKey) ok(feed *FeedRecord) PullResult {
	return PullResult{url: &pk.feedURL, status: pullSuccess, feed: feed, err: nil}
}

func (pk pullKey) err(e error) PullResult {
	return PullResult{url: &pk.feedURL, status: pullFail, feed: nil, err: e}
}

var (
	setFeedUpdateTime   = tableFieldSetter[string](feedsTable, "update_time")
	setFeedLastPullTime = tableFieldSetter[string](feedsTable, "last_pull_time")
)

func getPullKeys(ctx context.Context, tx *sql.Tx, feedIDs []ID) ([]pullKey, error) {
	// FIXME: Find a cleaner way to check for array membership using database/sql.
	//        Until then, we just loop through all IDs.
	stmt1, err := tx.PrepareContext(ctx, `SELECT feed_url FROM feeds WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	pks := make([]pullKey, len(feedIDs))
	for i, id := range feedIDs {
		pk := pullKey{feedID: id}
		if err := stmt1.QueryRowContext(ctx, pk.feedID).Scan(&pk.feedURL); err != nil {
			return nil, err
		}
		pks[i] = pk
	}

	return pks, nil
}

func getAllPullKeys(ctx context.Context, tx *sql.Tx) ([]pullKey, error) {

	sql1 := `SELECT id, feed_url FROM feeds`

	scanRow := func(rows *sql.Rows) (pullKey, error) {
		var pk pullKey
		err := rows.Scan(&pk.feedID, &pk.feedURL)
		return pk, err
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return nil, err
	}

	rows, err := stmt1.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	pks := make([]pullKey, 0)
	for rows.Next() {
		pk, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		pks = append(pks, pk)
	}

	return pks, nil
}

func pullNewFeedEntries(
	ctx context.Context,
	tx *sql.Tx,
	pk pullKey,
	parser FeedParser,
) chan PullResult {

	pullTime := time.Now().UTC().Format(time.RFC3339)
	pullf := func() PullResult {

		gfeed, err := parser.ParseURLWithContext(pk.feedURL, ctx)
		if err != nil {
			return pk.err(err)
		}

		updateTime := serializeTime(resolveFeedUpdateTime(gfeed))
		if err = setFeedUpdateTime(ctx, tx, pk.feedID, updateTime); err != nil {
			return pk.err(err)
		}
		if err = setFeedLastPullTime(ctx, tx, pk.feedID, &pullTime); err != nil {
			return pk.err(err)
		}

		if len(gfeed.Items) == 0 {
			return pk.ok(nil)
		}

		if err = upsertEntries(ctx, tx, pk.feedID, gfeed.Items); err != nil {
			return pk.err(err)
		}

		unreadEntries, err := getAllFeedEntries(ctx, tx, pk.feedID, pointer(false))
		if err != nil {
			return pk.err(err)
		}
		if len(unreadEntries) == 0 {
			return pk.ok(nil)
		}

		feed, err := getFeed(ctx, tx, pk.feedID)
		if err != nil {
			return pk.err(err)
		}

		feed.entries = unreadEntries

		return pk.ok(feed)
	}

	ic := make(chan PullResult)
	go func() {
		defer close(ic)
		ic <- pullf()
	}()

	oc := make(chan PullResult)
	go func() {
		defer close(oc)
		select {
		case <-ctx.Done():
			oc <- pk.err(ctx.Err())
		case msg := <-ic:
			oc <- msg
		}
	}()

	return oc
}

func merge[T any](chs []<-chan T) chan T {
	var (
		wg     sync.WaitGroup
		merged = make(chan T, len(chs))
	)

	forward := func(ch <-chan T) {
		for msg := range ch {
			merged <- msg
		}
		wg.Done()
	}

	wg.Add(len(chs))
	for _, ch := range chs {
		go forward(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
