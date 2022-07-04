package store

import (
	"context"
	"database/sql"
)

func (s *SQLite) ExportOPML(ctx context.Context) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("SQLite.ExportOPML")

	var payload []byte
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		feeds, err := getAllFeeds(ctx, tx)
		if err != nil {
			return fail(err)
		}
		if payload, err = Subscription(feeds).Export(); err != nil {
			return fail(err)
		}
		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
