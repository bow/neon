package store

import (
	"context"
	"fmt"
)

func (s *SQLite) ExportOPML(ctx context.Context) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("SQLite.ExportOPML")

	return nil, fail(fmt.Errorf("unimplemented"))
}
