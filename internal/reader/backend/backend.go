// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"

	"github.com/bow/neon/internal/entity"
)

// Backend describes the console backend.
type Backend interface {
	GetStatsF() func() (*entity.Stats, error)
	ListFeedsF() func() ([]*entity.Feed, error)
	PullFeeds(context.Context, []entity.ID) <-chan entity.PullResult
	StringF() func() string
}
