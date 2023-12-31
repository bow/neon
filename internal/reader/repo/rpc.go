// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/entity"
	"google.golang.org/grpc"
)

//nolint:unused
type RPC struct {
	addr   string
	client api.NeonClient

	statsCache *entity.Stats
}

// Ensure rpcRepo implements Repo.
var _ Repo = new(RPC)

func NewRPC(ctx context.Context, addr string, dialOpts ...grpc.DialOption) (*RPC, error) {
	conn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("timeout when connecting to server %q", addr)
		}
		return nil, err
	}
	return newRPCWithClient(addr, api.NewNeonClient(conn)), nil
}

func newRPCWithClient(addr string, client api.NeonClient) *RPC {
	return &RPC{addr: addr, client: client}
}

//nolint:unused
func (r *RPC) GetStats(ctx context.Context) (<-chan *entity.Stats, error) {
	panic("GetStats is unimplemented")
}

//nolint:unused
func (r *RPC) ListFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("ListFeeds is unimplemented")
}

//nolint:unused
func (r *RPC) PullFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("PullFeeds is unimplemented")
}
