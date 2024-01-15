// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bow/neon/api"
)

func TestGetStatsFOk(t *testing.T) {
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

func TestGetStatsFErr(t *testing.T) {
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

func TestListFeedsFOk(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	rpc, client := newBackendRPCTest(t)

	client.EXPECT().
		ListFeeds(gomock.Any(), gomock.Any()).
		Return(
			&api.ListFeedsResponse{
				Feeds: []*api.Feed{
					{
						Id:           uint32(5),
						Title:        "F1",
						FeedUrl:      "https://f1.com/feed.xml",
						SubTime:      timestamppb.New(time.Now()),
						LastPullTime: timestamppb.New(time.Now()),
					},
					{
						Id:           uint32(8),
						Title:        "F3",
						FeedUrl:      "https://f3.com/feed.xml",
						SubTime:      timestamppb.New(time.Now()),
						LastPullTime: timestamppb.New(time.Now()),
					},
				},
			},
			nil,
		)

	feeds, err := rpc.ListFeedsF(context.Background())()
	r.NoError(err)
	a.Len(feeds, 2)
	a.Equal(uint32(5), feeds[0].ID)
	a.Equal("F1", feeds[0].Title)
	a.Equal("https://f1.com/feed.xml", feeds[0].FeedURL)
	a.Equal(uint32(8), feeds[1].ID)
	a.Equal("F3", feeds[1].Title)
	a.Equal("https://f3.com/feed.xml", feeds[1].FeedURL)
}

func TestListFeedsFErr(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	rpc, client := newBackendRPCTest(t)

	client.EXPECT().
		ListFeeds(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("nope"))

	feeds, err := rpc.ListFeedsF(context.Background())()
	r.Nil(feeds)
	a.EqualError(err, "nope")
}

func newBackendRPCTest(t *testing.T) (*RPC, *MockNeonClient) {
	t.Helper()
	client := NewMockNeonClient(gomock.NewController(t))
	return newRPCWithClient("", client), client
}
