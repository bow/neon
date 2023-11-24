// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/bow/iris/internal"
)

func (db *SQLite) PullFeeds(
	ctx context.Context,
	ids []internal.ID,
) <-chan internal.PullResult {

	var (
		fail = failF("SQLite.PullFeeds")
		c    = make(chan internal.PullResult)
		wg   sync.WaitGroup
	)

	// nolint: unparam
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		defer wg.Done()

		var (
			pks []pullKey
			err error
		)
		if dedups := internal.Dedup(ids); len(dedups) == 0 {
			pks, err = getAllPullKeys(ctx, tx)
		} else {
			pks, err = getPullKeys(ctx, tx, dedups)
		}
		if err != nil {
			c <- internal.NewPullResultFromError(nil, fail(err))
			return nil
		}
		if len(pks) == 0 {
			c <- internal.NewPullResultFromFeed(nil, nil)
			return nil
		}

		chs := make([]<-chan internal.PullResult, len(pks))
		for i, pk := range pks {
			chs[i] = pullNewFeedEntries(ctx, tx, pk, db.parser)
		}

		for pr := range merge(chs) {
			pr := pr
			if e := pr.Error(); e != nil {
				pr.SetError(fail(e))
			}
			c <- pr
		}

		return nil
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	go func() {
		defer func() {
			wg.Wait()
			close(c)
		}()
		wg.Add(1)
		err := db.withTx(ctx, dbFunc)
		if err != nil {
			c <- internal.NewPullResultFromError(nil, fail(err))
		}
	}()

	return c
}

type pullKey struct {
	feedID  ID
	feedURL string
}

func (pk pullKey) ok(feed *internal.Feed) internal.PullResult {
	pr := internal.NewPullResultFromFeed(&pk.feedURL, feed)
	pr.SetStatus(internal.PullSuccess)
	return pr
}

func (pk pullKey) err(e error) internal.PullResult {
	pr := internal.NewPullResultFromError(&pk.feedURL, e)
	pr.SetStatus(internal.PullFail)
	return pr
}

var (
	setFeedUpdateTime   = tableFieldSetter[time.Time](feedsTable, "update_time")
	setFeedLastPullTime = tableFieldSetter[time.Time](feedsTable, "last_pull_time")
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
	parser internal.FeedParser,
) chan internal.PullResult {

	pullTime := time.Now().UTC()
	pullf := func() internal.PullResult {

		gfeed, err := parser.ParseURLWithContext(pk.feedURL, ctx)
		if err != nil {
			return pk.err(err)
		}

		updateTime := resolveFeedUpdateTime(gfeed)
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

		defaultReadStatus := false
		unreadEntries, err := getEntries(ctx, tx, []ID{pk.feedID}, &defaultReadStatus, nil)
		if err != nil {
			return pk.err(err)
		}
		if len(unreadEntries) == 0 {
			return pk.ok(nil)
		}

		rec, err := getFeed(ctx, tx, pk.feedID)
		if err != nil {
			return pk.err(err)
		}

		rec.entries = unreadEntries

		return pk.ok(rec.feed())
	}

	ic := make(chan internal.PullResult)
	go func() {
		defer close(ic)
		ic <- pullf()
	}()

	oc := make(chan internal.PullResult)
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
