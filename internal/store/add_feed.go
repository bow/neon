// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/mmcdole/gofeed"
)

// AddFeed adds the given feed into the database.
func (s *SQLite) AddFeed(
	ctx context.Context,
	feedURL string,
	title *string,
	desc *string,
	tags []string,
	isStarred *bool,
) (*Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("SQLite.AddFeed")

	feed, err := s.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fail(err)
	}
	// Handle possible specs deviations.
	if feed.FeedLink == "" {
		feed.FeedLink = feedURL
	}

	var created *Feed
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		now := time.Now()

		feedDBID, _, ierr := upsertFeed(
			ctx,
			tx,
			feed.FeedLink,
			pointerOrNil(deref(title, feed.Title)),
			pointerOrNil(deref(desc, feed.Description)),
			pointerOrNil(feed.Link),
			isStarred,
			resolveFeedUpdateTime(feed),
			&now,
		)
		if ierr != nil {
			return ierr
		}

		if ierr = upsertEntries(ctx, tx, feedDBID, feed.Items); ierr != nil {
			return ierr
		}

		if ierr = addFeedTags(ctx, tx, feedDBID, tags); ierr != nil {
			return ierr
		}

		if created, ierr = getFeed(ctx, tx, feedDBID); ierr != nil {
			return ierr
		}

		return nil
	}

	err = s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return created, nil
}

func upsertFeed(
	ctx context.Context,
	tx *sql.Tx,
	feedURL string,
	title *string,
	desc *string,
	siteURL *string,
	isStarred *bool,
	updateTime *time.Time,
	subTime *time.Time,
) (feedDBID DBID, added bool, err error) {

	sql1 := `
		INSERT INTO
			feeds(
				feed_url,
				title,
				description,
				site_url,
				is_starred,
				update_time,
				sub_time,
				last_pull_time
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return feedDBID, added, err
	}
	defer stmt1.Close()

	sst := serializeTime(subTime)
	res, err := stmt1.ExecContext(
		ctx,
		feedURL,
		title,
		desc,
		siteURL,
		deref(isStarred, false),
		serializeTime(updateTime),
		sst,
		sst, // last_pull_time defaults to sub_time.
	)

	if err == nil {
		lid, ierr := res.LastInsertId()
		feedDBID = DBID(lid)
		if ierr != nil {
			return feedDBID, added, ierr
		}
		added = true
	} else {
		if !isUniqueErr(err, "UNIQUE constraint failed: feeds.feed_url") {
			return feedDBID, added, err
		}
		var ierr error
		if feedDBID, ierr = editFeedWithFeedURL(
			ctx,
			tx,
			feedURL,
			title,
			desc,
			siteURL,
			isStarred,
		); ierr != nil {
			return feedDBID, added, err
		}
		added = false
	}

	return feedDBID, added, nil
}

func editFeedWithFeedURL(
	ctx context.Context,
	tx *sql.Tx,
	feedURL string,
	title *string,
	desc *string,
	siteURL *string,
	isStarred *bool,
) (DBID, error) {

	var feedDBID DBID

	sql1 := `SELECT id FROM feeds WHERE feed_url = ?`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return 0, err
	}

	if err := stmt1.QueryRowContext(ctx, feedURL).Scan(&feedDBID); err != nil {
		return 0, err
	}
	if err := setFeedTitle(ctx, tx, feedDBID, title); err != nil {
		return 0, err
	}
	if err := setFeedDescription(ctx, tx, feedDBID, desc); err != nil {
		return 0, err
	}
	if err := setFeedIsStarred(ctx, tx, feedDBID, isStarred); err != nil {
		return 0, err
	}
	if err := setFeedSiteURL(ctx, tx, feedDBID, siteURL); err != nil {
		return 0, err
	}
	return feedDBID, nil
}

func upsertEntries(
	ctx context.Context,
	tx *sql.Tx,
	feedDBID DBID,
	entries []*gofeed.Item,
) error {

	sql1 := `
		INSERT INTO
			entries(
				feed_id,
				external_id,
				url,
				title,
				description,
				content,
				pub_time,
				update_time
			)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?)
`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()

	// TODO: Also update other entry columns.
	sql2 := `
		UPDATE
			entries
		SET
			is_read = false,
			update_time = $1
		WHERE
			feed_id = $2
			AND external_id = $3
			AND (
				update_time IS NULL AND $1 IS NOT NULL
				OR update_time IS NOT NULL AND $1 IS NULL
				OR update_time != $1
			)
`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()

	upsert := func(entry *gofeed.Item, insertStmt, updateStmt *sql.Stmt) error {
		updateTime := serializeTime(resolveEntryUpdateTime(entry))
		_, err := insertStmt.ExecContext(
			ctx,
			feedDBID,
			entry.GUID,
			entry.Link,
			entry.Title,
			pointerOrNil(entry.Description),
			pointerOrNil(entry.Content),
			serializeTime(resolveEntryPublishedTime(entry)),
			updateTime,
		)
		if err != nil {
			if !isUniqueErr(err, "UNIQUE constraint failed: entries.feed_id, entries.external_id") {
				return err
			}
			if _, ierr := updateStmt.ExecContext(
				ctx,
				updateTime,
				feedDBID,
				entry.GUID,
			); ierr != nil {
				return ierr
			}
		}
		return nil
	}

	for _, entry := range entries {
		if err := upsert(entry, stmt1, stmt2); err != nil {
			return err
		}
	}
	return nil
}

func addFeedTags(
	ctx context.Context,
	tx *sql.Tx,
	feedDBID DBID,
	tags []string,
) error {

	sql1 := `INSERT OR IGNORE INTO feed_tags(name) VALUES (?)`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()
	for _, tag := range tags {
		_, err = stmt1.ExecContext(ctx, tag)
		if err != nil {
			return err
		}
	}

	sql2 := `SELECT id FROM feed_tags WHERE name = ?`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()
	ids := make(map[string]DBID)
	for _, tag := range tags {
		if _, exists := ids[tag]; exists {
			continue
		}
		var id DBID
		row := stmt2.QueryRowContext(ctx, tag)
		if err = row.Scan(&id); err != nil {
			return err
		}
		ids[tag] = id
	}

	sql3 := `INSERT OR IGNORE INTO feeds_x_feed_tags(feed_id, feed_tag_id) VALUES (?, ?)`
	stmt3, err := tx.PrepareContext(ctx, sql3)
	if err != nil {
		return err
	}
	defer stmt3.Close()

	for _, catDBID := range ids {
		if _, err := stmt3.ExecContext(ctx, feedDBID, catDBID); err != nil {
			return err
		}
	}

	return nil
}
