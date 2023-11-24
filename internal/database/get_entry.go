// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bow/iris/internal"
)

func (db *SQLite) GetEntry(
	ctx context.Context,
	id internal.ID,
) (*internal.Entry, error) {

	var rec *entryRecord
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		irec, err := getEntry(ctx, tx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return internal.EntryNotFoundError{ID: id}
			}
			return err
		}
		rec = irec
		return nil
	}

	fail := failF("SQLite.GetEntry")

	db.mu.Lock()
	defer db.mu.Unlock()

	err := db.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}

	return rec.entry(), nil
}

func getEntry(ctx context.Context, tx *sql.Tx, entryID ID) (*entryRecord, error) {

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
			e.id = $1
		ORDER BY
			COALESCE(e.update_time, e.pub_time) DESC
`
	scanRow := func(row *sql.Row) (*entryRecord, error) {
		var entry entryRecord
		if err := row.Scan(
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

	return scanRow(stmt1.QueryRowContext(ctx, entryID))
}
