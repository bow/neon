// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"fmt"
)

func (s *SQLite) ListEntries(ctx context.Context) ([]*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return nil, fmt.Errorf("unimplemented")
}
