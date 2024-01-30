// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/chanutil"
	"github.com/bow/neon/internal/entity"
)

func (db *SQLite) PullFeeds(
	ctx context.Context,
	ids []entity.ID,
	entryReadStatus *bool,
	maxEntriesPerFeed *uint32,
) <-chan entity.PullResult {

	var (
		fail = failF("SQLite.PullFeeds")
		c    = make(chan entity.PullResult)
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
			c <- entity.NewPullResultFromError(nil, fail(err))
			return nil
		}
		if len(pks) == 0 {
			return nil
		}

		chs := make([]<-chan entity.PullResult, len(pks))
		for i, pk := range pks {
			chs[i] = pullFeedEntries(
				ctx,
				tx,
				pk,
				db.parser,
				entryReadStatus,
				maxEntriesPerFeed,
			)
		}

		for pr := range chanutil.Merge(chs) {
			pr := pr
			if e := pr.Error(); e != nil {
				pr.SetError(fail(e))
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

		db.mu.Lock()
		defer db.mu.Unlock()

		err := db.withTx(ctx, dbFunc)
		if err != nil {
			c <- entity.NewPullResultFromError(nil, fail(err))
		}
	}()

	return c
}

type pullKey struct {
	feedID  ID
	feedURL string
}

func (pk pullKey) ok(feed *entity.Feed) entity.PullResult {
	pr := entity.NewPullResultFromFeed(&pk.feedURL, feed)
	pr.SetStatus(entity.PullSuccess)
	return pr
}

func (pk pullKey) err(e error) entity.PullResult {
	pr := entity.NewPullResultFromError(&pk.feedURL, e)
	pr.SetStatus(entity.PullFail)
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

func pullFeedEntries(
	ctx context.Context,
	tx *sql.Tx,
	pk pullKey,
	parser Parser,
	entryReadStatus *bool,
	maxEntriesPerFeed *uint32,
) chan entity.PullResult {

	pullTime := time.Now().UTC()
	pullf := func() entity.PullResult {

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

		entries, err := getEntries(
			ctx,
			tx,
			[]ID{pk.feedID},
			maxEntriesPerFeed,
			entryReadStatus,
			nil,
		)
		if err != nil {
			return pk.err(err)
		}
		if len(entries) == 0 && maxEntriesPerFeed == nil {
			return pk.ok(nil)
		}

		rec, err := getFeed(ctx, tx, pk.feedID)
		if err != nil {
			return pk.err(err)
		}

		rec.entries = entries

		return pk.ok(rec.feed())
	}

	ic := make(chan entity.PullResult)
	go func() {
		defer close(ic)
		ic <- pullf()
	}()

	oc := make(chan entity.PullResult)
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
