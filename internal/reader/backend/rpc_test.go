// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package backend

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal/entity"
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
		ListFeeds(gomock.Any(), gomock.Any(), gomock.Any()).
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
		ListFeeds(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("nope"))

	feeds, err := rpc.ListFeedsF(context.Background())()
	r.Nil(feeds)
	a.EqualError(err, "nope")
}

func TestPullFeedsFExtended(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	rpc, client := newBackendRPCTest(t)
	streamClient := NewMockNeon_PullFeedsClient(gomock.NewController(t))

	client.EXPECT().
		PullFeeds(gomock.Any(), gomock.Any()).
		Return(streamClient, nil)
	streamClient.EXPECT().
		Recv().
		Return(
			&api.PullFeedsResponse{
				Url:   "https://err.com/feed.xml",
				Feed:  nil,
				Error: pointer("http 404"),
			},
			nil,
		)
	streamClient.EXPECT().
		Recv().
		Return(
			&api.PullFeedsResponse{
				Url: "https://ok.com/feed.xml",
				Feed: &api.Feed{
					Id:           uint32(5),
					Title:        "OK",
					FeedUrl:      "https://ok.com/feed.xml",
					SubTime:      timestamppb.New(time.Now()),
					LastPullTime: timestamppb.New(time.Now()),
				},
				Error: nil,
			},
			nil,
		)
	streamClient.EXPECT().
		Recv().
		Return(nil, io.EOF)

	ch, err := rpc.PullFeedsF(context.Background(), nil)()
	r.NoError(err)
	a.NotNil(ch)

	prs := make([]entity.PullResult, 0)
	for pr := range ch {
		prs = append(prs, pr)
	}
	r.Len(prs, 2)

	pr0 := prs[0] // #nosec: G602
	a.Equal("https://err.com/feed.xml", pr0.URL())
	a.Nil(pr0.Feed())
	a.EqualError(pr0.Error(), "http 404")

	pr1 := prs[1] // #nosec: G602
	a.Equal("https://ok.com/feed.xml", pr1.URL())
	a.NotNil(pr1.Feed())
	a.Equal("OK", pr1.Feed().Title)
	a.Nil(pr1.Error())
}

func TestPullFeedsFErr(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	rpc, client := newBackendRPCTest(t)

	client.EXPECT().
		PullFeeds(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("call fail"))

	ch, err := rpc.PullFeedsF(context.Background(), nil)()
	r.Nil(ch)
	a.EqualError(err, "call fail")
}

func TestPullFeedsFErrStream(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	rpc, client := newBackendRPCTest(t)
	streamClient := NewMockNeon_PullFeedsClient(gomock.NewController(t))

	client.EXPECT().
		PullFeeds(gomock.Any(), gomock.Any()).
		Return(streamClient, nil)
	streamClient.EXPECT().
		Recv().
		Return(nil, fmt.Errorf("stream fail"))

	ch, err := rpc.PullFeedsF(context.Background(), nil)()
	r.NoError(err)
	a.NotNil(ch)

	prs := make([]entity.PullResult, 0)
	for pr := range ch {
		prs = append(prs, pr)
	}
	r.Len(prs, 1)

	pr := prs[0] // #nosec: G602
	a.Empty(pr.URL())
	a.Nil(pr.Feed())
	a.EqualError(pr.Error(), "stream fail")
}

func newBackendRPCTest(t *testing.T) (*RPC, *MockNeonClient) {
	t.Helper()
	client := NewMockNeonClient(gomock.NewController(t))
	return newRPCWithClient("", client), client
}

func pointer[T any](value T) *T { return &value }
