// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"

	"github.com/bow/neon/internal/entity"
)

// Backend describes the console backend.
type Backend interface {
	GetStats(context.Context) (<-chan *entity.Stats, error)
	ListFeeds(context.Context) (<-chan *entity.Feed, error)
	PullFeeds(context.Context) (<-chan *entity.Feed, error)
	String() string
}
