package store

import (
	"context"
	"database/sql"
)

func (s *SQLite) ExportOPML(ctx context.Context, title *string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var payload []byte
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		feeds, err := getAllFeeds(ctx, tx)
		if err != nil {
			return err
		}
		if payload, err = Subscription(feeds).Export(title); err != nil {
			return err
		}
		return nil
	}

	fail := failF("SQLite.ExportOPML")

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, fail(err)
	}
	return payload, nil
}
