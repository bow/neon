// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *SQLite) GetGlobalStats(ctx context.Context) (*Stats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := &Stats{}
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		return fmt.Errorf("unimplemented")
	}

	fail := failF("SQLite.GetGlobalStats")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return stats, nil
}
