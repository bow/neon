// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *SQLite) ViewEntry(ctx context.Context, _ DBID) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var entry *Entry
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		return fmt.Errorf("unimplemented")
	}

	fail := failF("SQLite.ViewFeed")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return entry, nil
}
