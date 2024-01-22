// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"

	"github.com/bow/neon/internal/entity"
)

// Backend describes the console backend.
type Backend interface {
	GetStatsF(context.Context) func() (*entity.Stats, error)
	ListFeedsF(context.Context) func() ([]*entity.Feed, error)
	PullFeedsF(context.Context, []entity.ID, bool) func() (<-chan entity.PullResult, error)
	String() string
}
