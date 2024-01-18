// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/entity"
	"google.golang.org/grpc"
)

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

func newRPCWithClient(
	addr string,
	client api.NeonClient,
) *RPC {
	return &RPC{addr: addr, client: client}
}

func (r *RPC) GetStatsF(ctx context.Context) func() (*entity.Stats, error) {
	return func() (*entity.Stats, error) {
		rsp, err := r.client.GetStats(ctx, &api.GetStatsRequest{})
		if err != nil {
			return nil, err
		}
		stats := entity.FromStatsPb(rsp.GetGlobal())
		return stats, nil
	}
}

func (r *RPC) ListFeedsF(ctx context.Context) func() ([]*entity.Feed, error) {
	return func() ([]*entity.Feed, error) {
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
}

func (r *RPC) PullFeedsF(
	ctx context.Context,
	_ []entity.ID,
) func() (<-chan entity.PullResult, error) {
	return func() (<-chan entity.PullResult, error) {
		stream, err := r.client.PullFeeds(ctx, &api.PullFeedsRequest{})
		if err != nil {
			return nil, err
		}

		ch := make(chan entity.PullResult)
		go func() {
			defer close(ch)
			for {
				rsp, serr := stream.Recv()
				if serr != nil {
					if serr != io.EOF {
						ch <- entity.NewPullResultFromError(nil, serr)
					}
					return
				}
				if perr := rsp.Error; perr != nil {
					ch <- entity.NewPullResultFromError(&rsp.Url, fmt.Errorf("%s", *perr))
					continue
				}
				ch <- entity.NewPullResultFromFeed(&rsp.Url, entity.FromFeedPb(rsp.GetFeed()))
			}
		}()
		return ch, nil
	}
}

func (r *RPC) String() string {
	return fmt.Sprintf("grpc://%s", r.addr)
}
