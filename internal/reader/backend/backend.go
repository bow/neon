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
	GetAllFeedsF(context.Context) func() ([]*entity.Feed, error)
	PullFeedsF(context.Context, []entity.ID) func() (<-chan entity.PullResult, error)
	String() string
}

type result[T any] struct {
	value T
	err   error
}

func okResult[T any](value T) result[T] {
	return result[T]{value: value, err: nil}
}

func errResult[T any](err error) result[T] {
	return result[T]{err: err}
}
