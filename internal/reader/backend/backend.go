// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"

	"github.com/bow/neon/internal/entity"
)

// Backend describes the console backend.
type Backend interface {
	GetStats(context.Context) <-chan Result[*entity.Stats]
	ListFeeds(context.Context) <-chan Result[*entity.Feed]
	PullFeeds(context.Context) <-chan Result[*entity.Feed]
	String() string
}

type Result[T any] struct {
	Value T
	Err   error
}

func OkResult[T any](value T) Result[T] {
	return Result[T]{Value: value, Err: nil}
}

func ErrResult[T any](err error) Result[T] {
	return Result[T]{Err: err}
}
