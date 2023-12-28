// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"database/sql"
	"time"

	"github.com/bow/neon/internal/entity"
)

func (db *SQLite) ImportSubscription(
	ctx context.Context,
	sub *entity.Subscription,
) (processed int, imported int, err error) {

	if len(sub.Feeds) == 0 {
		return 0, 0, nil
	}

	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		now := time.Now()

		for _, feed := range sub.Feeds {
			f := feed
			feedID, isAdded, ierr := upsertFeed(
				ctx,
				tx,
				f.FeedURL,
				pointerOrNil(f.Title),
				f.Description,
				f.SiteURL,
				&f.IsStarred,
				nil,
				&now,
			)
			if ierr != nil {
				return ierr
			}

			if ierr = addFeedTags(ctx, tx, feedID, f.Tags); ierr != nil {
				return ierr
			}
			processed++
			if isAdded {
				imported++
			}
		}

		return nil
	}

	fail := failF("SQLite.ImportSubscription")

	db.mu.Lock()
	defer db.mu.Unlock()

	err = db.withTx(ctx, dbFunc)
	if err != nil {
		return 0, 0, fail(err)
	}

	return processed, imported, nil
}
