// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

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
}

// Ensure rpcRepo implements Repo.
var _ Backend = new(RPC)

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
func (r *RPC) GetStats(ctx context.Context) (*entity.Stats, error) {
	rsp, err := r.client.GetStats(ctx, &api.GetStatsRequest{})
	if err != nil {
		return nil, err
	}
	stats := entity.FromStatsPb(rsp.GetGlobal())
	return stats, nil
}

//nolint:unused
func (r *RPC) ListFeeds(ctx context.Context) <-chan Result[*entity.Feed] {
	panic("ListFeeds is unimplemented")
}

//nolint:unused
func (r *RPC) PullFeeds(ctx context.Context) <-chan Result[*entity.Feed] {
	panic("PullFeeds is unimplemented")
}

func (r *RPC) String() string {
	return fmt.Sprintf("grpc://%s", r.addr)
}
