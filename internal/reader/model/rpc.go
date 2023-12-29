// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/entity"
	"google.golang.org/grpc"
)

//nolint:unused
type rpcModel struct {
	addr   string
	client api.NeonClient

	statsCache *entity.Stats
}

// Ensure rpcModel implements Model.
var _ Model = new(rpcModel)

func NewRPCModel(ctx context.Context, addr string, dialOpts ...grpc.DialOption) (*rpcModel, error) {
	conn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("timeout when connecting to server %q", addr)
		}
		return nil, err
	}
	return newRPCModelWithClient(addr, api.NewNeonClient(conn)), nil
}

func newRPCModelWithClient(addr string, client api.NeonClient) *rpcModel {
	return &rpcModel{addr: addr, client: client}
}

//nolint:unused
func (m *rpcModel) GetStats(ctx context.Context) (<-chan *entity.Stats, error) {
	panic("GetStats is unimplemented")
}

//nolint:unused
func (m *rpcModel) ListFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("ListFeeds is unimplemented")
}

//nolint:unused
func (m *rpcModel) PullFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("PullFeeds is unimplemented")
}
