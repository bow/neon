// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package repo

import (
	"context"

	"github.com/bow/neon/internal/entity"
)

// Repo describes the console data layer.
type Repo interface {
	GetStats(context.Context) (<-chan *entity.Stats, error)
	ListFeeds(context.Context) (<-chan *entity.Feed, error)
	PullFeeds(context.Context) (<-chan *entity.Feed, error)
	Backend() string
}
