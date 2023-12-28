// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package nreader

import (
	"context"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/entity"
)

//nolint:unused
type model struct {
	addr   string
	client api.NeonClient

	statsCache *entity.Stats
}

//nolint:unused
func (m *model) ListFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("ListFeeds is unimplemented")
}

//nolint:unused
func (m *model) PullFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("PullFeeds is unimplemented")
}

//nolint:unused
func (m *model) GetStats(ctx context.Context) (<-chan *entity.Stats, error) {
	panic("GetStats is unimplemented")
}
