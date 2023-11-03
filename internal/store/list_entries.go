// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"errors"
)

func (s *SQLite) ListEntries(ctx context.Context, feedID ID) ([]*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries := make([]*Entry, 0)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		_, err := getFeed(ctx, tx, feedID)
		if errors.Is(err, sql.ErrNoRows) {
			return FeedNotFoundError{feedID}
		}
		ientries, err := getAllFeedEntries(ctx, tx, feedID, nil)
		entries = ientries
		return err
	}

	fail := failF("SQLite.ListEntries")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return entries, nil
}
