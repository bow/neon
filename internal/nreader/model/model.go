// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package model

import (
	"context"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/entity"
)

type Model interface {
	ListFeeds(context.Context) (<-chan *entity.Feed, error)
	PullFeeds(context.Context) (<-chan *entity.Feed, error)
	GetStats(context.Context) (<-chan *entity.Feed, error)
}

//nolint:unused
type rpcModel struct {
	addr   string
	client api.NeonClient

	statsCache *entity.Stats
}

//nolint:unused
func (m *rpcModel) ListFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("ListFeeds is unimplemented")
}

//nolint:unused
func (m *rpcModel) PullFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("PullFeeds is unimplemented")
}

//nolint:unused
func (m *rpcModel) GetStats(ctx context.Context) (<-chan *entity.Stats, error) {
	panic("GetStats is unimplemented")
}
