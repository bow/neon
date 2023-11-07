// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"

	"github.com/bow/iris/internal"
)

// EditEntries updates fields of an entry.
func (s *SQLite) EditEntries(
	ctx context.Context,
	ops []*internal.EntryEditOp,
) ([]*internal.Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updateFunc := func(
		ctx context.Context,
		tx *sql.Tx, op *internal.EntryEditOp,
	) (*entryRecord, error) {
		if err := setEntryIsRead(ctx, tx, op.ID, op.IsRead); err != nil {
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

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}

	return toEntries(recs)
}

var setEntryIsRead = tableFieldSetter[bool](entriesTable, "is_read")
