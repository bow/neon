// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/bow/neon/internal/entity"
)

func (db *SQLite) ListEntries(
	ctx context.Context,
	feedIDs []entity.ID,
	isBookmarked *bool,
) ([]*entity.Entry, error) {

	recs := make([]*entryRecord, 0)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		// TODO: Expose num arg in function signature.
		irecs, err := getEntries(ctx, tx, feedIDs, nil, nil, isBookmarked)
		if err != nil {
			return err
		}
		recs = irecs
		return nil
	}

	fail := failF("SQLite.ListEntries")

	db.mu.RLock()
	defer db.mu.RUnlock()

	err := db.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}

	return entryRecords(recs).entriesSlice(), nil
}

func getEntries(
	ctx context.Context,
	tx *sql.Tx,
	feedIDs []ID,
	numMaxEntries *uint32,
	isRead *bool,
	isBookmarked *bool,
) ([]*entryRecord, error) {

	sql1 := `
		SELECT
			e.id AS id
			, e.feed_id AS feed_id
			, e.title AS title
			, e.is_read AS is_read
			, e.is_bookmarked AS is_bookmarked
			, e.external_id AS ext_id
			, e.description AS description
			, e.content AS content
			, e.url AS url
			, e.update_time AS update_time
			, e.pub_time AS pub_time
		FROM
			entries e
		WHERE
			COALESCE(e.feed_id IN (SELECT value FROM json_each($1)), true)
			AND COALESCE(e.is_read = $2, true)
			AND COALESCE(e.is_bookmarked = $3, true)
		ORDER BY
			COALESCE(e.update_time, e.pub_time) DESC
`

	if numMaxEntries != nil {
		nmax := *numMaxEntries
		if nmax == 0 {
			return nil, nil
		}
		sql1 += fmt.Sprintf("\nLIMIT %d", nmax)
	}

	scanRow := func(rows *sql.Rows) (*entryRecord, error) {
		var entry entryRecord
		if err := rows.Scan(
			&entry.id,
			&entry.feedID,
			&entry.title,
			&entry.isRead,
			&entry.isBookmarked,
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

	feedIDsJSON := "null"
	if len(feedIDs) > 0 {
		s, merr := json.Marshal(feedIDs)
		if merr != nil {
			return nil, merr
		}
		feedIDsJSON = string(s)
	}

	rows, err := stmt1.QueryContext(ctx, feedIDsJSON, isRead, isBookmarked)
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
