// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"database/sql"

	"github.com/bow/neon/internal/entity"
)

// EditFeed updates fields of an feed.
func (db *SQLite) EditFeeds(
	ctx context.Context,
	ops []*entity.FeedEditOp,
) ([]*entity.Feed, error) {

	updateFunc := func(
		ctx context.Context,
		tx *sql.Tx, op *entity.FeedEditOp,
	) (*feedRecord, error) {
		if err := setFeedTitle(ctx, tx, op.ID, op.Title); err != nil {
			return nil, err
		}
		if err := setFeedDescription(ctx, tx, op.ID, op.Description); err != nil {
			return nil, err
		}
		if err := setFeedTags(ctx, tx, op.ID, op.Tags); err != nil {
			return nil, err
		}
		if err := setFeedIsStarred(ctx, tx, op.ID, op.IsStarred); err != nil {
			return nil, err
		}
		return getFeed(ctx, tx, op.ID)
	}

	var feeds = make([]*feedRecord, len(ops))
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		for i, op := range ops {
			feed, err := updateFunc(ctx, tx, op)
			if err != nil {
				return err
			}
			feeds[i] = feed
		}
		return nil
	}

	fail := failF("SQLite.EditFeed")

	db.mu.Lock()
	defer db.mu.Unlock()

	err := db.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}

	return feedRecords(feeds).feeds(), nil
}

func getFeed(ctx context.Context, tx *sql.Tx, feedID ID) (*feedRecord, error) {

	sql1 := `
		SELECT
			f.id AS id
			, f.title AS title
			, f.description AS description
			, f.feed_url AS feed_url
			, f.site_url AS site_url
			, f.is_starred AS is_starred
			, f.sub_time AS sub_time
			, f.update_time AS update_time
			, f.last_pull_time AS last_pull_time
			, json_group_array(fc.name) FILTER (WHERE fc.name IS NOT NULL) AS tags
		FROM
			feeds f
			LEFT JOIN feeds_x_feed_tags fxfc ON fxfc.feed_id = f.id
			LEFT JOIN feed_tags fc ON fxfc.feed_tag_id = fc.id
		WHERE
			f.id = ?
		GROUP BY
			f.id
		ORDER BY
			COALESCE(f.update_time, f.sub_time) DESC
`
	scanRow := func(row *sql.Row) (*feedRecord, error) {
		var feed feedRecord
		if err := row.Scan(
			&feed.id,
			&feed.title,
			&feed.description,
			&feed.feedURL,
			&feed.siteURL,
			&feed.isStarred,
			&feed.subscribed,
			&feed.updated,
			&feed.lastPulled,
			&feed.tags,
		); err != nil {
			return nil, err
		}
		if len(feed.tags) == 0 {
			feed.tags = nil
		}
		return &feed, nil
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	return scanRow(stmt1.QueryRowContext(ctx, feedID))
}

var (
	setFeedTitle       = tableFieldSetter[string](feedsTable, "title")
	setFeedDescription = tableFieldSetter[string](feedsTable, "description")
	setFeedIsStarred   = tableFieldSetter[bool](feedsTable, "is_starred")
	setFeedSiteURL     = tableFieldSetter[string](feedsTable, "site_url")
)

func setFeedTags(
	ctx context.Context,
	tx *sql.Tx,
	feedID ID,
	tags *[]string,
) error {

	if tags == nil {
		return nil
	}

	sql1 := `DELETE FROM feeds_x_feed_tags WHERE feed_id = ?`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()

	if _, err = stmt1.ExecContext(ctx); err != nil {
		return err
	}

	if err = addFeedTags(ctx, tx, feedID, *tags); err != nil {
		return err
	}

	sql2 := `
		DELETE
			feed_tags
		WHERE
			id IN (
				SELECT
					fc.id
				FROM
					feed_tags fc
					LEFT JOIN feeds_x_feed_tags fxfc ON fxfc.feed_tag_id = fc.id
				WHERE
					fxfc.feed_id IS NULL
			)
	`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()

	return nil
}
