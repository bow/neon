// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"database/sql"

	"github.com/bow/neon/internal/entity"
)

func (db *SQLite) ExportSubscription(
	ctx context.Context,
	title *string,
) (*entity.Subscription, error) {

	var sub entity.Subscription
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		recs, err := getAllFeeds(ctx, tx)
		if err != nil {
			return err
		}
		isub := entity.Subscription{
			Title: title,
			Feeds: feedRecords(recs).feeds(),
		}
		sub = isub
		return nil
	}

	fail := failF("SQLite.ExportSubscription")

	db.mu.RLock()
	defer db.mu.RUnlock()

	err := db.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return &sub, nil
}
