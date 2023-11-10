// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"

	"github.com/bow/iris/internal"
)

func (s *SQLite) DeleteFeeds(ctx context.Context, ids []internal.ID) error {

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

	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return fail(err)
	}

	return nil
}
