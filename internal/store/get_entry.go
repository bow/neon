// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"errors"
)

func (s *SQLite) GetEntry(ctx context.Context, entryID ID) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var entry *Entry
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		ientry, err := getEntry(ctx, tx, entryID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return EntryNotFoundError{entryID}
			}
			return err
		}
		entry = ientry
		return nil
	}

	fail := failF("SQLite.ViewFeed")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return entry, nil
}

func getEntry(ctx context.Context, tx *sql.Tx, entryID ID) (*Entry, error) {

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
			e.id = $1
		ORDER BY
			COALESCE(e.update_time, e.pub_time) DESC
`
	scanRow := func(row *sql.Row) (*Entry, error) {
		var entry Entry
		if err := row.Scan(
			&entry.ID,
			&entry.FeedID,
			&entry.Title,
			&entry.IsRead,
			&entry.ExtID,
			&entry.Description,
			&entry.Content,
			&entry.URL,
			&entry.Updated,
			&entry.Published,
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
