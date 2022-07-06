package store

import (
	"context"
	"database/sql"
)

func (s *SQLite) DeleteFeeds(ctx context.Context, ids []DBID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		sql1 := `DELETE FROM feeds WHERE id = ?`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return err
		}

		deleteFunc := func(ctx context.Context, tx *sql.Tx, id DBID) error {
			res, err := stmt1.ExecContext(ctx, id)
			if err != nil {
				return err
			}
			n, err := res.RowsAffected()
			if err != nil {
				return err
			}
			if n != int64(1) {
				return FeedNotFoundError{id}
			}
			return nil
		}

		for _, id := range ids {
			if err := deleteFunc(ctx, tx, id); err != nil {
				return err
			}
		}

		return nil
	}

	fail := failF("SQLite.DeleteFeeds")

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return fail(err)
	}
	return nil
}
