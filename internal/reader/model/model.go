// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package model

import (
	"context"

	"github.com/bow/neon/internal/entity"
)

// Model describes the console data layer.
type Model interface {
	GetStats(context.Context) (<-chan *entity.Stats, error)
	ListFeeds(context.Context) (<-chan *entity.Feed, error)
	PullFeeds(context.Context) (<-chan *entity.Feed, error)
}
