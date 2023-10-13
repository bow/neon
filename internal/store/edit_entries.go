// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
)

// EditEntries updates fields of an entry.
func (s *SQLite) EditEntries(
	ctx context.Context,
	ops []*EntryEditOp,
) ([]*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updateFunc := func(ctx context.Context, tx *sql.Tx, op *EntryEditOp) (*Entry, error) {
		if err := setEntryIsRead(ctx, tx, op.DBID, op.IsRead); err != nil {
			return nil, err
		}
		return getEntry(ctx, tx, op.DBID)
	}

	var entries = make([]*Entry, len(ops))
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		for i, op := range ops {
			entry, err := updateFunc(ctx, tx, op)
			if err != nil {
				return err
			}
			entries[i] = entry
		}
		return nil
	}

	fail := failF("SQLite.EditEntries")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return entries, nil
}

var setEntryIsRead = setTableField[bool]("entries", "is_read")
