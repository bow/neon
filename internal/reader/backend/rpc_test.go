// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bow/neon/api"
)

func TestGetStatsOk(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	rpc, client := newBackendRPCTest(t)

	client.EXPECT().
		GetStats(gomock.Any(), gomock.Any()).
		Return(
			&api.GetStatsResponse{Global: &api.GetStatsResponse_Stats{NumFeeds: 5}},
			nil,
		)

	stats, err := rpc.GetStatsF(context.Background())()
	r.NoError(err)
	a.Equal(uint32(5), stats.NumFeeds)
}

func TestGetStatsErr(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	rpc, client := newBackendRPCTest(t)

	client.EXPECT().
		GetStats(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("nope"))

	stats, err := rpc.GetStatsF(context.Background())()
	r.Nil(stats)
	a.EqualError(err, "nope")
}

func newBackendRPCTest(t *testing.T) (*RPC, *MockNeonClient) {
	t.Helper()
	client := NewMockNeonClient(gomock.NewController(t))
	return newRPCWithClient("", client), client
}
