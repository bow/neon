// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"database/sql"

	"github.com/bow/neon/internal"
)

func (db *SQLite) DeleteFeeds(ctx context.Context, ids []internal.ID) error {

	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		sql1 := `DELETE FROM feeds WHERE id = ?`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return err
		}

		deleteFunc := func(ctx context.Context, id ID) error {
			res, err := stmt1.ExecContext(ctx, id)
			if err != nil {
				return err
			}
			n, err := res.RowsAffected()
			if err != nil {
				return err
			}
			if n != int64(1) {
				return internal.FeedNotFoundError{ID: id}
			}
			return nil
		}

		for _, id := range internal.Dedup(ids) {
			if err := deleteFunc(ctx, id); err != nil {
				return err
			}
		}

		return nil
	}

	fail := failF("SQLite.DeleteFeeds")

	db.mu.Lock()
	defer db.mu.Unlock()

	err := db.withTx(ctx, dbFunc)
	if err != nil {
		return fail(err)
	}

	return nil
}
