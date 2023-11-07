// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"

	"github.com/bow/iris/internal"
)

func (s *SQLite) ListFeeds(ctx context.Context) ([]*internal.Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	feeds := make([]*feedRecord, 0)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		ifeeds, err := getAllFeeds(ctx, tx)
		if err != nil {
			return err
		}
		for _, ifeed := range ifeeds {
			ifeed := ifeed
			entries, err := getAllFeedEntries(ctx, tx, ifeed.id, nil)
			if err != nil {
				return err
			}
			ifeed.entries = entries
		}
		feeds = ifeeds

		return nil
	}

	fail := failF("SQLite.ListFeeds")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return toFeeds(feeds)
}

func getAllFeeds(ctx context.Context, tx *sql.Tx) ([]*feedRecord, error) {

	sql1 := `
		SELECT
			f.id AS id,
			f.title AS title,
			f.description AS description,
			f.feed_url AS feed_url,
			f.site_url AS site_url,
			f.is_starred AS is_starred,
			f.sub_time AS sub_time,
			f.last_pull_time AS last_pull_time,
			f.update_time AS update_time,
			json_group_array(fc.name) FILTER (WHERE fc.name IS NOT NULL) AS tags
		FROM
			feeds f
			LEFT JOIN feeds_x_feed_tags fxfc ON fxfc.feed_id = f.id
			LEFT JOIN feed_tags fc ON fxfc.feed_tag_id = fc.id
		GROUP BY
			f.id
		ORDER BY
			COALESCE(f.update_time, f.sub_time) DESC
`
	scanRow := func(rows *sql.Rows) (*feedRecord, error) {
		var feed feedRecord
		if err := rows.Scan(
			&feed.id,
			&feed.title,
			&feed.description,
			&feed.feedURL,
			&feed.siteURL,
			&feed.isStarred,
			&feed.subscribed,
			&feed.lastPulled,
			&feed.updated,
			&feed.tags,
		); err != nil {
			return nil, err
		}
		return &feed, nil
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	rows, err := stmt1.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	feeds := make([]*feedRecord, 0)
	for rows.Next() {
		feed, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

func getAllFeedEntries(
	ctx context.Context,
	tx *sql.Tx,
	feedID ID,
	isRead *bool,
) ([]*entryRecord, error) {

	sql1 := `
		SELECT
			e.id AS id,
			e.feed_id AS feed_id,
			e.title AS title,
			e.is_read AS is_read,
			e.external_id AS ext_id,
			e.description AS description,
			e.content AS content,
			e.url AS url,
			e.update_time AS update_time,
			e.pub_time AS pub_time
		FROM
			entries e
		WHERE
			e.feed_id = $1
			AND COALESCE(e.is_read = $2, true)
		ORDER BY
			COALESCE(e.update_time, e.pub_time) DESC
`
	scanRow := func(rows *sql.Rows) (*entryRecord, error) {
		var entry entryRecord
		if err := rows.Scan(
			&entry.id,
			&entry.feedID,
			&entry.title,
			&entry.isRead,
			&entry.extID,
			&entry.description,
			&entry.content,
			&entry.url,
			&entry.updated,
			&entry.published,
		); err != nil {
			return nil, err
		}
		return &entry, nil
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	rows, err := stmt1.QueryContext(ctx, feedID, isRead)
	if err != nil {
		return nil, err
	}

	entries := make([]*entryRecord, 0)
	for rows.Next() {
		entry, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
