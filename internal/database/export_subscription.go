// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"database/sql"

	"github.com/bow/neon/internal"
)

func (db *SQLite) ExportSubscription(
	ctx context.Context,
	title *string,
) (*internal.Subscription, error) {

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

	db.mu.Lock()
	defer db.mu.Unlock()

	err := db.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return &sub, nil
}
