// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"database/sql"

	"github.com/bow/neon/internal/entity"
)

// EditEntries updates fields of an entry.
func (db *SQLite) EditEntries(
	ctx context.Context,
	ops []*entity.EntryEditOp,
) ([]*entity.Entry, error) {

	updateFunc := func(
		ctx context.Context,
		tx *sql.Tx, op *entity.EntryEditOp,
	) (*entryRecord, error) {
		if err := setEntryIsRead(ctx, tx, op.ID, op.IsRead); err != nil {
			return nil, err
		}
		if err := setEntryIsBookmarked(ctx, tx, op.ID, op.IsBookmarked); err != nil {
			return nil, err
		}
		return getEntry(ctx, tx, op.ID)
	}

	recs := make([]*entryRecord, len(ops))
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		for i, op := range ops {
			rec, err := updateFunc(ctx, tx, op)
			if err != nil {
				return err
			}
			recs[i] = rec
		}
		return nil
	}

	fail := failF("SQLite.EditEntries")

	db.mu.Lock()
	defer db.mu.Unlock()

	err := db.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}

	return entryRecords(recs).entriesSlice(), nil
}

var (
	setEntryIsRead       = tableFieldSetter[bool](entriesTable, "is_read")
	setEntryIsBookmarked = tableFieldSetter[bool](entriesTable, "is_bookmarked")
)
