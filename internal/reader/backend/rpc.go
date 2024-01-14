// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/entity"
	"google.golang.org/grpc"
)

//nolint:unused
type RPC struct {
	addr   string
	client api.NeonClient

	ctx         context.Context
	callTimeout time.Duration
}

// Ensure rpcRepo implements Repo.
var _ Backend = new(RPC)

func NewRPC(
	ctx context.Context,
	callTimeout time.Duration,
	addr string,
	dialOpts ...grpc.DialOption,
) (*RPC, error) {
	conn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("timeout when connecting to server %q", addr)
		}
		return nil, err
	}
	return newRPCWithClient(ctx, callTimeout, addr, api.NewNeonClient(conn)), nil
}

func newRPCWithClient(
	ctx context.Context,
	callTimeout time.Duration,
	addr string,
	client api.NeonClient,
) *RPC {
	return &RPC{ctx: ctx, addr: addr, client: client, callTimeout: callTimeout}
}

//nolint:unused
func (r *RPC) GetStatsF() func() (*entity.Stats, error) {
	return func() (*entity.Stats, error) {
		ctx, cancel := r.callCtx()
		defer cancel()

		rsp, err := r.client.GetStats(ctx, &api.GetStatsRequest{})
		if err != nil {
			return nil, err
		}
		stats := entity.FromStatsPb(rsp.GetGlobal())
		return stats, nil
	}
}

func (r *RPC) ListFeeds(ctx context.Context) ([]*entity.Feed, error) {
	rsp, err := r.client.ListFeeds(ctx, &api.ListFeedsRequest{})
	if err != nil {
		return nil, err
	}
	rfeeds := rsp.GetFeeds()
	feeds := make([]*entity.Feed, len(rfeeds))
	for i, rfeed := range rfeeds {
		feeds[i] = entity.FromFeedPb(rfeed)
	}
	return feeds, nil
}

//nolint:unused
func (r *RPC) PullFeeds(ctx context.Context, _ []entity.ID) <-chan entity.PullResult {
	panic("PullFeeds is unimplemented")
}

func (r *RPC) StringF() func() string {
	return func() string { return fmt.Sprintf("grpc://%s", r.addr) }
}

func (r *RPC) callCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.ctx, r.callTimeout)
}
