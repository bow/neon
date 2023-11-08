// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"

	"github.com/bow/iris/internal"
)

func (s *SQLite) ExportSubscription(
	ctx context.Context,
	title *string,
) (*internal.Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sub internal.Subscription
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		recs, err := getAllFeeds(ctx, tx)
		if err != nil {
			return err
		}
		isub := internal.Subscription{
			Title: title,
			Feeds: feedRecords(recs).feeds(),
		}
		sub = isub
		return nil
	}

	fail := failF("SQLite.ExportSubscription")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return &sub, nil
}
