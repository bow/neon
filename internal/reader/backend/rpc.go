// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/chanutil"
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

func (r *RPC) GetAllFeedsF(ctx context.Context) func() ([]*entity.Feed, error) {
	return func() ([]*entity.Feed, error) {
		feeds, err := r.listEmptyFeeds(ctx)
		if err != nil {
			return nil, err
		}
		return r.fillEmptyFeeds(ctx, feeds)
	}
}

func (r *RPC) PullFeedsF(
	ctx context.Context,
	_ []entity.ID,
) func() (<-chan entity.PullResult, error) {
	return func() (<-chan entity.PullResult, error) {
		max := uint32(0)
		req := api.PullFeedsRequest{MaxEntriesPerFeed: &max}
		stream, err := r.client.PullFeeds(ctx, &req)
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

func (r *RPC) listEmptyFeeds(ctx context.Context) ([]*entity.Feed, error) {
	nmax := uint32(0)
	rsp, err := r.client.ListFeeds(
		ctx,
		&api.ListFeedsRequest{MaxEntriesPerFeed: &nmax},
	)
	if err != nil {
		return nil, err
	}

	return entity.FromFeedPbs(rsp.GetFeeds()), nil
}

func (r *RPC) fillEmptyFeeds(
	ctx context.Context,
	feeds []*entity.Feed,
) ([]*entity.Feed, error) {

	chs := make([]<-chan result[*entity.Feed], len(feeds))
	for i, feed := range feeds {
		feed := feed
		ch := make(chan result[*entity.Feed])
		chs[i] = ch
		go func() {
			defer close(ch)
			stream, err := r.client.StreamEntries(
				ctx,
				&api.StreamEntriesRequest{FeedId: feed.ID},
			)
			if err != nil {
				ch <- errResult[*entity.Feed](err)
				return
			}
			for {
				srsp, serr := stream.Recv()
				if serr != nil {
					if serr != io.EOF {
						ch <- errResult[*entity.Feed](serr)
					} else {
						ch <- okResult(feed)
					}
					return
				}
				entry := entity.FromEntryPb(srsp.GetEntry())
				feed.Entries[entry.ID] = entry
			}
		}()
	}

	filled := make([]*entity.Feed, 0)
	for res := range chanutil.Merge(chs) {
		res := res
		if err := res.err; err != nil {
			return nil, err
		}
		filled = append(filled, res.value)
	}
	return filled, nil
}
