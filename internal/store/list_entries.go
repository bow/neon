// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bow/iris/internal"
)

func (s *SQLite) ListEntries(
	ctx context.Context,
	feedID internal.ID,
) ([]*internal.Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	recs := make([]*entryRecord, 0)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		_, err := getFeed(ctx, tx, feedID)
		if errors.Is(err, sql.ErrNoRows) {
			return FeedNotFoundError{feedID}
		}
		irecs, err := getAllFeedEntries(ctx, tx, feedID, nil)
		recs = irecs
		return err
	}

	fail := failF("SQLite.ListEntries")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}

	return entryRecords(recs).entries(), nil
}
