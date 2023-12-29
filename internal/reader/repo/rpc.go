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
type rpcRepo struct {
	addr   string
	client api.NeonClient

	statsCache *entity.Stats
}

// Ensure rpcRepo implements Repo.
var _ Repo = new(rpcRepo)

func NewRPCRepo(ctx context.Context, addr string, dialOpts ...grpc.DialOption) (*rpcRepo, error) {
	conn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("timeout when connecting to server %q", addr)
		}
		return nil, err
	}
	return newRPCRepoWithClient(addr, api.NewNeonClient(conn)), nil
}

func newRPCRepoWithClient(addr string, client api.NeonClient) *rpcRepo {
	return &rpcRepo{addr: addr, client: client}
}

//nolint:unused
func (m *rpcRepo) GetStats(ctx context.Context) (<-chan *entity.Stats, error) {
	panic("GetStats is unimplemented")
}

//nolint:unused
func (m *rpcRepo) ListFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("ListFeeds is unimplemented")
}

//nolint:unused
func (m *rpcRepo) PullFeeds(ctx context.Context) (<-chan *entity.Feed, error) {
	panic("PullFeeds is unimplemented")
}
