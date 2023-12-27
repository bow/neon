// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/bow/neon/internal"
)

// AddFeed adds the given feed into the database.
func (db *SQLite) AddFeed(
	ctx context.Context,
	feedURL string,
	title *string,
	desc *string,
	tags []string,
	isStarred *bool,
) (*internal.Feed, bool, error) {

	fail := failF("SQLite.AddFeed")

	feed, err := db.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, false, err
	}
	// Handle possible specs deviations.
	if feed.FeedLink == "" {
		feed.FeedLink = feedURL
	}

	var (
		record *feedRecord
		added  = pointer(false)
	)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		now := time.Now()

		feedID, feedAdded, ierr := upsertFeed(
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

		if ierr = upsertEntries(ctx, tx, feedID, feed.Items); ierr != nil {
			return ierr
		}

		if len(tags) > 0 {
			if ierr = addFeedTags(ctx, tx, feedID, tags); ierr != nil {
				return ierr
			}
		} else {
			if ierr = removeFeedTags(ctx, tx, feedID); ierr != nil {
				return ierr
			}
		}

		if record, ierr = getFeed(ctx, tx, feedID); ierr != nil {
			return ierr
		}
		added = &feedAdded

		return nil
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	err = db.withTx(ctx, dbFunc)
	if err != nil {
		return nil, *added, fail(err)
	}

	return record.feed(), *added, nil
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
) (feedID ID, added bool, err error) {

	sql1 := `
		INSERT INTO
			feeds(
				feed_url
				, title
				, description
				, site_url
				, is_starred
				, update_time
				, sub_time
				, last_pull_time
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return feedID, added, err
	}
	defer stmt1.Close()

	res, err := stmt1.ExecContext(
		ctx,
		feedURL,
		title,
		desc,
		siteURL,
		deref(isStarred, false),
		updateTime,
		subTime,
		subTime, // last_pull_time defaults to sub_time.
	)

	if err == nil {
		lid, ierr := res.LastInsertId()
		feedID = ID(lid)
		if ierr != nil {
			return feedID, added, ierr
		}
		added = true
	} else {
		if !isUniqueErr(err, "UNIQUE constraint failed: feeds.feed_url") {
			return feedID, added, err
		}
		var ierr error
		if feedID, ierr = editFeedWithFeedURL(
			ctx,
			tx,
			feedURL,
			title,
			desc,
			siteURL,
			isStarred,
		); ierr != nil {
			return feedID, added, err
		}
		added = false
	}

	return feedID, added, nil
}

func editFeedWithFeedURL(
	ctx context.Context,
	tx *sql.Tx,
	feedURL string,
	title *string,
	desc *string,
	siteURL *string,
	isStarred *bool,
) (ID, error) {

	var feedID ID

	sql1 := `SELECT id FROM feeds WHERE feed_url = ?`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return 0, err
	}

	if err := stmt1.QueryRowContext(ctx, feedURL).Scan(&feedID); err != nil {
		return 0, err
	}
	if err := setFeedTitle(ctx, tx, feedID, title); err != nil {
		return 0, err
	}
	if err := setFeedDescription(ctx, tx, feedID, desc); err != nil {
		return 0, err
	}
	if err := setFeedIsStarred(ctx, tx, feedID, isStarred); err != nil {
		return 0, err
	}
	if err := setFeedSiteURL(ctx, tx, feedID, siteURL); err != nil {
		return 0, err
	}
	return feedID, nil
}

func upsertEntries(
	ctx context.Context,
	tx *sql.Tx,
	feedID ID,
	entries []*gofeed.Item,
) error {

	sql1 := `
		INSERT INTO
			entries(
				feed_id
				, external_id
				, url
				, title
				, description
				, content
				, pub_time
				, update_time
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
		updateTime := resolveEntryUpdateTime(entry)
		_, err := insertStmt.ExecContext(
			ctx,
			feedID,
			entry.GUID,
			entry.Link,
			entry.Title,
			pointerOrNil(entry.Description),
			pointerOrNil(entry.Content),
			resolveEntryPublishedTime(entry),
			updateTime,
		)
		if err != nil {
			if !isUniqueErr(err, "UNIQUE constraint failed: entries.feed_id, entries.external_id") {
				return err
			}
			if _, ierr := updateStmt.ExecContext(
				ctx,
				updateTime,
				feedID,
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
	feedID ID,
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
	ids := make(map[string]ID)
	for _, tag := range tags {
		if _, exists := ids[tag]; exists {
			continue
		}
		var id ID
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

	for _, catID := range ids {
		if _, err := stmt3.ExecContext(ctx, feedID, catID); err != nil {
			return err
		}
	}

	return nil
}

func removeFeedTags(
	ctx context.Context,
	tx *sql.Tx,
	feedID ID,
) error {
	sql1 := `DELETE FROM feeds_x_feed_tags WHERE feed_id = ?`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()

	_, err = stmt1.ExecContext(ctx, feedID)
	if err != nil {
		return err
	}
	return nil
}
