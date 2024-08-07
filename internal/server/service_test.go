// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package server

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"regexp"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/entity"
)

func TestAddFeedOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, ds := setupServerTest(t)

	req := api.AddFeedRequest{
		Url:         "http://foo.com/feed.xml",
		Title:       pointer("user-title"),
		Description: pointer("user-description"),
		Tags:        []string{"tag-1", "tag-2", "tag-3"},
		IsStarred:   pointer(true),
	}
	record := &entity.Feed{
		ID:          entity.ID(5),
		Title:       "feed-title-original",
		Description: pointer("feed-description-original"),
		FeedURL:     "https://foo.com/feed.xml",
		SiteURL:     pointer("https://foo.com"),
		Subscribed:  mustTimeVV(t, "2021-07-01T23:33:06.156+02:00"),
		LastPulled:  mustTimeVV(t, "2021-07-01T23:33:06.156+02:00"),
		IsStarred:   true,
	}

	ds.EXPECT().
		AddFeed(
			gomock.Any(),
			req.GetUrl(),
			req.Title,
			req.Description,
			req.Tags,
			req.IsStarred,
			nil,
		).
		Return(record, true, nil)

	rsp, err := client.AddFeed(context.Background(), &req)
	r.NoError(err)

	a.True(rsp.IsAdded)
	a.Equal(record.Title, rsp.Feed.GetTitle())
	a.Equal(record.Description, rsp.Feed.Description)
	a.Equal(record.SiteURL, rsp.Feed.SiteUrl)
	a.Equal(record.FeedURL, rsp.Feed.FeedUrl)
	a.Equal(record.IsStarred, rsp.Feed.IsStarred)
}

func TestListFeedsOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, ds := setupServerTest(t)

	req := api.ListFeedsRequest{}
	feeds := []*entity.Feed{
		{
			ID:         entity.ID(2),
			Title:      "Feed A",
			FeedURL:    "http://a.com/feed.xml",
			Subscribed: mustTimeVV(t, "2022-06-22T19:39:38.964+02:00"),
			LastPulled: mustTimeVV(t, "2022-06-22T19:39:38.964+02:00"),
			Updated:    pointer(mustTimeVV(t, "2022-03-19T16:23:18.600+02:00")),
		},
		{
			ID:         entity.ID(3),
			Title:      "Feed X",
			FeedURL:    "http://x.com/feed.xml",
			Subscribed: mustTimeVV(t, "2022-06-22T19:39:44.037+02:00"),
			LastPulled: mustTimeVV(t, "2022-06-22T19:39:44.037+02:00"),
			Updated:    pointer(mustTimeVV(t, "2022-04-20T16:32:30.760+02:00")),
		},
	}

	ds.EXPECT().
		ListFeeds(gomock.Any(), req.MaxEntriesPerFeed).
		Return(feeds, nil)

	rsp, err := client.ListFeeds(context.Background(), &req)
	r.NoError(err)

	// TODO: Expand test.
	a.Len(rsp.GetFeeds(), 2)
}

func TestEditFeedsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, ds := setupServerTest(t)

	ops := []*entity.FeedEditOp{
		{ID: 14, Title: pointer("newer")},
		{ID: 58, Tags: pointer([]string{"x", "y"})},
		{ID: 77, IsStarred: pointer(true)},
	}
	feeds := []*entity.Feed{
		{
			ID:         14,
			Title:      "newer",
			Subscribed: mustTimeVV(t, "2022-06-30T00:53:50.200+02:00"),
			LastPulled: mustTimeVV(t, "2022-06-30T00:53:50.200+02:00"),
		},
		{
			ID:         58,
			Tags:       []string{"x", "y"},
			Subscribed: mustTimeVV(t, "2022-06-30T00:53:58.135+02:00"),
			LastPulled: mustTimeVV(t, "2022-06-30T00:53:58.135+02:00"),
		},
		{
			ID:         77,
			IsStarred:  true,
			Subscribed: mustTimeVV(t, "2022-06-30T00:53:59.812+02:00"),
			LastPulled: mustTimeVV(t, "2022-06-30T00:53:59.812+02:00"),
		},
	}

	ds.EXPECT().
		EditFeeds(gomock.Any(), gomock.AssignableToTypeOf(ops)).
		Return(feeds, nil)

	req := api.EditFeedsRequest{
		Ops: []*api.EditFeedsRequest_Op{
			{
				Id: 14,
				Fields: &api.EditFeedsRequest_Op_Fields{
					Title: pointer("newer"),
				},
			},
			{
				Id: 58,
				Fields: &api.EditFeedsRequest_Op_Fields{
					Tags: []string{"x", "y"},
				},
			},
			{
				Id: 77,
				Fields: &api.EditFeedsRequest_Op_Fields{
					IsStarred: pointer(true),
				},
			},
		},
	}
	rsp, err := client.EditFeeds(context.Background(), &req)
	r.NoError(err)

	r.Len(rsp.Feeds, 3)
	feed0 := rsp.Feeds[0]
	a.Equal(feeds[0].ID, feed0.GetId())
	a.Equal(feeds[0].Title, feed0.GetTitle())
	feed1 := rsp.Feeds[1]
	a.Equal(feeds[1].ID, feed1.GetId())
	a.Equal(feeds[1].Tags, feed1.GetTags())
	feed2 := rsp.Feeds[2]
	a.Equal(feeds[2].ID, feed2.GetId())
	a.Equal(feeds[2].IsStarred, feed2.GetIsStarred())
}

func TestDeleteFeedsOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, ds := setupServerTest(t)

	ds.EXPECT().
		DeleteFeeds(gomock.Any(), []entity.ID{1, 9}).
		Return(nil)

	req := api.DeleteFeedsRequest{FeedIds: []uint32{1, 9}}
	rsp, err := client.DeleteFeeds(context.Background(), &req)
	r.NoError(err)

	a.True(proto.Equal(&api.DeleteFeedsResponse{}, rsp))
}

func TestDeleteFeedsErrNotFound(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, ds := setupServerTest(t)

	ds.EXPECT().
		DeleteFeeds(gomock.Any(), []entity.ID{1, 9}).
		Return(fmt.Errorf("wrapped: %w", entity.FeedNotFoundError{ID: 9}))

	req := api.DeleteFeedsRequest{FeedIds: []uint32{1, 9}}
	rsp, err := client.DeleteFeeds(context.Background(), &req)

	r.Nil(rsp)
	a.EqualError(err, "rpc error: code = NotFound desc = feed with ID=9 not found")
}

func TestPullFeedsAllOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, ds := setupServerTest(t)

	prs := []entity.PullResult{
		entity.NewPullResultFromFeed(
			pointer("http://a.com/feed.xml"),
			&entity.Feed{
				Title:      "feed-A",
				FeedURL:    "https://a.com/feed.xml",
				Subscribed: mustTimeVV(t, "2021-07-23T17:20:29.499+02:00"),
				LastPulled: mustTimeVV(t, "2021-07-23T17:20:29.499+02:00"),
				IsStarred:  true,
				Entries: map[entity.ID]*entity.Entry{
					1: {Title: "Entry A1", IsRead: false},
					2: {Title: "Entry A2", IsRead: false},
				},
			},
		),
		entity.NewPullResultFromFeed(pointer("http://z.com/feed.xml"), nil),
		entity.NewPullResultFromFeed(
			pointer("http://c.com/feed.xml"),
			&entity.Feed{
				Title:      "feed-C",
				FeedURL:    "https://c.com/feed.xml",
				Subscribed: mustTimeVV(t, "2021-07-23T17:21:11.489+02:00"),
				LastPulled: mustTimeVV(t, "2021-07-23T17:21:11.489+02:00"),
				IsStarred:  false,
				Entries: map[entity.ID]*entity.Entry{
					3: {Title: "Entry C3", IsRead: false},
				},
			},
		),
	}

	ch := make(chan entity.PullResult)
	go func() {
		defer close(ch)

		// Randomize ordering, to simulate actual URL pulls.
		shufres := make([]entity.PullResult, len(prs))
		copy(shufres, prs)
		r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec: G404
		shf := func(i, j int) { shufres[i], shufres[j] = shufres[j], shufres[i] }
		r.Shuffle(len(shufres), shf)

		for i := 0; i < len(shufres); i++ {
			ch <- shufres[i]
		}
	}()

	ds.EXPECT().
		PullFeeds(gomock.Any(), []entity.ID{}, gomock.Any(), gomock.Any(), nil).
		Return(ch)

	req := api.PullFeedsRequest{}
	stream, err := client.PullFeeds(context.Background(), &req)
	r.NoError(err)

	var (
		rsp       *api.PullFeedsResponse
		errStream error
		rsps      = make([]*api.PullFeedsResponse, 2)
	)

	for i := 0; i < len(rsps); i++ {
		rsp, errStream = stream.Recv()
		a.NoError(errStream)
		a.NotNil(rsp)
		rsps[i] = rsp
	}

	rsp, errStream = stream.Recv()
	a.ErrorIs(errStream, io.EOF)
	a.Nil(rsp)

	// Sort responses so tests are insensitive to input order.
	sort.SliceStable(rsps, func(i, j int) bool { return rsps[i].GetUrl() < rsps[j].GetUrl() })

	rsp0 := rsps[0]
	r.Equal(prs[0].URL(), rsp0.GetUrl())
	r.Nil(rsp0.Error)
	r.NotNil(rsp0.Feed)
	a.Len(rsp0.GetFeed().GetEntries(), 2)

	rsp1 := rsps[1]
	r.Equal(prs[2].URL(), rsp1.GetUrl())
	r.Nil(rsp1.Error)
	r.NotNil(rsp0.Feed)
	a.Len(rsp1.GetFeed().GetEntries(), 1)
}

func TestPullFeedsSelectedAllOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, ds := setupServerTest(t)

	prs := []entity.PullResult{
		entity.NewPullResultFromFeed(pointer("http://z.com/feed.xml"), nil),
		entity.NewPullResultFromFeed(
			pointer("http://c.com/feed.xml"),
			&entity.Feed{
				Title:      "feed-C",
				FeedURL:    "https://c.com/feed.xml",
				Subscribed: mustTimeVV(t, "2021-07-23T17:21:11.489+02:00"),
				LastPulled: mustTimeVV(t, "2021-07-23T17:21:11.489+02:00"),
				IsStarred:  false,
				Entries: map[entity.ID]*entity.Entry{
					3: {Title: "Entry C3", IsRead: false},
				},
			},
		),
	}

	ch := make(chan entity.PullResult)
	go func() {
		defer close(ch)

		// Randomize ordering, to simulate actual URL pulls.
		shufres := make([]entity.PullResult, len(prs))
		copy(shufres, prs)
		r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec: G404
		shf := func(i, j int) { shufres[i], shufres[j] = shufres[j], shufres[i] }
		r.Shuffle(len(shufres), shf)

		for i := 0; i < len(shufres); i++ {
			ch <- shufres[i]
		}
	}()

	ds.EXPECT().
		PullFeeds(gomock.Any(), []entity.ID{2, 3}, gomock.Any(), gomock.Any(), nil).
		Return(ch)

	req := api.PullFeedsRequest{FeedIds: []uint32{2, 3}}
	stream, err := client.PullFeeds(context.Background(), &req)
	r.NoError(err)

	var (
		rsp       *api.PullFeedsResponse
		errStream error
		rsps      = make([]*api.PullFeedsResponse, 1)
	)

	for i := 0; i < len(rsps); i++ {
		rsp, errStream = stream.Recv()
		a.NoError(errStream)
		a.NotNil(rsp)
		rsps[i] = rsp
	}

	rsp, errStream = stream.Recv()
	a.ErrorIs(errStream, io.EOF)
	a.Nil(rsp)

	// Sort responses so tests are insensitive to input order.
	sort.SliceStable(rsps, func(i, j int) bool { return rsps[i].GetUrl() < rsps[j].GetUrl() })

	rsp0 := rsps[0]
	r.Equal(prs[1].URL(), rsp0.GetUrl())
	r.Nil(rsp0.Error)
	r.NotNil(rsp0.Feed)
	a.Len(rsp0.GetFeed().GetEntries(), 1)
}

func TestPullFeedsErrSomeFeed(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, ds := setupServerTest(t)

	prs := []entity.PullResult{
		entity.NewPullResultFromFeed(
			pointer("https://a.com/feed.xml"),
			&entity.Feed{
				Title:      "feed-A",
				FeedURL:    "https://a.com/feed.xml",
				Subscribed: mustTimeVV(t, "2021-07-23T17:20:29.499+02:00"),
				LastPulled: mustTimeVV(t, "2021-07-23T17:20:29.499+02:00"),
				IsStarred:  true,
				Entries: map[entity.ID]*entity.Entry{
					1: {Title: "Entry A1", IsRead: false},
					2: {Title: "Entry A2", IsRead: false},
				},
			},
		),
		entity.NewPullResultFromError(pointer("https://x.com/feed.xml"), fmt.Errorf("timed out")),
		entity.NewPullResultFromFeed(pointer("https://z.com/feed.xml"), nil),
		entity.NewPullResultFromFeed(
			pointer("https://c.com/feed.xml"),
			&entity.Feed{
				Title:      "feed-C",
				FeedURL:    "https://c.com/feed.xml",
				Subscribed: mustTimeVV(t, "2021-07-23T17:21:11.489+02:00"),
				LastPulled: mustTimeVV(t, "2021-07-23T17:21:11.489+02:00"),
				IsStarred:  false,
				Entries: map[entity.ID]*entity.Entry{
					3: {Title: "Entry C3", IsRead: false},
				},
			},
		),
	}

	ch := make(chan entity.PullResult)
	go func() {
		defer close(ch)

		// Randomize ordering, to simulate actual URL pulls.
		shufres := make([]entity.PullResult, len(prs))
		copy(shufres, prs)
		r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec: G404
		shf := func(i, j int) { shufres[i], shufres[j] = shufres[j], shufres[i] }
		r.Shuffle(len(shufres), shf)

		for i := 0; i < len(shufres); i++ {
			ch <- shufres[i]
		}
	}()

	ds.EXPECT().
		PullFeeds(gomock.Any(), []entity.ID{}, gomock.Any(), gomock.Any(), nil).
		Return(ch)

	req := api.PullFeedsRequest{}
	stream, err := client.PullFeeds(context.Background(), &req)
	r.NoError(err)

	var (
		rsp       *api.PullFeedsResponse
		errStream error
		rsps      = make([]*api.PullFeedsResponse, 3)
	)

	for i := 0; i < len(rsps); i++ {
		rsp, errStream = stream.Recv()
		a.NoError(errStream)
		a.NotNil(rsp)
		rsps[i] = rsp
	}

	rsp, errStream = stream.Recv()
	a.ErrorIs(errStream, io.EOF)
	a.Nil(rsp)

	// Sort responses so tests are insensitive to input order.
	sort.SliceStable(rsps, func(i, j int) bool { return rsps[i].GetUrl() < rsps[j].GetUrl() })

	rsp0 := rsps[0]
	r.Equal(prs[0].URL(), rsp0.GetUrl())
	r.NotNil(rsp0.Feed)
	a.Len(rsp0.GetFeed().GetEntries(), 2)

	rsp1 := rsps[1]
	r.Equal(prs[3].URL(), rsp1.GetUrl())
	a.Len(rsp1.GetFeed().GetEntries(), 1)

	rsp2 := rsps[2]
	r.Equal(prs[1].URL(), rsp2.GetUrl())
	a.Nil(rsp2.GetFeed())
	a.Equal("timed out", rsp2.GetError())
}

func TestPullFeedsErrNonFeed(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, ds := setupServerTest(t)

	prs := []entity.PullResult{
		entity.NewPullResultFromFeed(
			pointer("https://a.com/feed.xml"),
			&entity.Feed{
				Title:      "feed-A",
				FeedURL:    "https://a.com/feed.xml",
				Subscribed: mustTimeVV(t, "2021-07-23T17:20:29.499+02:00"),
				LastPulled: mustTimeVV(t, "2021-07-23T17:20:29.499+02:00"),
				IsStarred:  true,
				Entries: map[entity.ID]*entity.Entry{
					1: {Title: "Entry A1", IsRead: false},
					2: {Title: "Entry A2", IsRead: false},
				},
			},
		),
		entity.NewPullResultFromFeed(pointer("https://z.com/feed.xml"), nil),
		entity.NewPullResultFromFeed(
			pointer("https://c.com/feed.xml"),
			&entity.Feed{
				Title:      "feed-C",
				FeedURL:    "https://c.com/feed.xml",
				Subscribed: mustTimeVV(t, "2021-07-23T17:21:11.489+02:00"),
				LastPulled: mustTimeVV(t, "2021-07-23T17:21:11.489+02:00"),
				IsStarred:  false,
				Entries: map[entity.ID]*entity.Entry{
					3: {Title: "Entry C3", IsRead: false},
				},
			},
		),
		entity.NewPullResultFromError(nil, fmt.Errorf("tx error")),
	}

	ch := make(chan entity.PullResult)
	go func() {
		defer close(ch)

		// Randomize ordering, to simulate actual URL pulls.
		shufres := make([]entity.PullResult, len(prs))
		copy(shufres, prs)
		r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec: G404
		shf := func(i, j int) { shufres[i], shufres[j] = shufres[j], shufres[i] }
		r.Shuffle(len(shufres), shf)

		for i := 0; i < len(shufres); i++ {
			ch <- shufres[i]
		}
	}()

	ds.EXPECT().
		PullFeeds(gomock.Any(), []entity.ID{}, gomock.Any(), gomock.Any(), nil).
		Return(ch)

	req := api.PullFeedsRequest{}
	stream, err := client.PullFeeds(context.Background(), &req)
	r.NoError(err)

	var (
		rsp       *api.PullFeedsResponse
		errStream error
		rsps      = make([]*api.PullFeedsResponse, 0)
	)

	for {
		rsp, errStream = stream.Recv()
		if errStream != nil {
			err = errStream
			break
		}
		if errStream == io.EOF {
			break
		}
		rsps = append(rsps, rsp)
	}

	a.LessOrEqual(len(rsps), 3)
	a.EqualError(err, "rpc error: code = Unknown desc = tx error")
}

func TestListEntriesOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, ds := setupServerTest(t)

	req := api.ListEntriesRequest{FeedIds: []entity.ID{2}, IsBookmarked: pointer(true)}
	entries := []*entity.Entry{
		{
			Title:        "Entry 1",
			IsRead:       false,
			IsBookmarked: true,
			Content:      pointer("Contents 1."),
		},
		{
			Title:        "Entry 2",
			IsRead:       false,
			IsBookmarked: true,
			Content:      pointer("Contents 2."),
		},
		{
			Title:        "Entry 3",
			IsRead:       true,
			IsBookmarked: true,
			Content:      pointer("Contents 3."),
		},
	}

	ds.EXPECT().
		ListEntries(gomock.Any(), req.GetFeedIds(), req.IsBookmarked).
		Return(entries, nil)

	rsp, err := client.ListEntries(context.Background(), &req)
	r.NoError(err)

	// TODO: Expand test.
	a.Len(rsp.GetEntries(), 3)
}

func TestEditEntriesOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, ds := setupServerTest(t)

	ops := []*entity.EntryEditOp{
		{ID: 37, IsRead: pointer(true), IsBookmarked: pointer(true)},
		{ID: 49, IsRead: pointer(false)},
	}
	entries := []*entity.Entry{
		{ID: 37, IsRead: true, IsBookmarked: true},
		{ID: 49, IsRead: false},
	}

	ds.EXPECT().
		EditEntries(gomock.Any(), ops).
		Return(entries, nil)

	req := api.EditEntriesRequest{
		Ops: []*api.EditEntriesRequest_Op{
			{
				Id: 37,
				Fields: &api.EditEntriesRequest_Op_Fields{
					IsRead:       pointer(true),
					IsBookmarked: pointer(true),
				},
			},
			{
				Id: 49,
				Fields: &api.EditEntriesRequest_Op_Fields{
					IsRead: pointer(false),
				},
			},
		},
	}
	rsp, err := client.EditEntries(context.Background(), &req)
	r.NoError(err)

	r.Len(rsp.Entries, 2)
	entry0 := rsp.Entries[0]
	a.Equal(entries[0].ID, entry0.Id)
	a.Equal(entries[0].IsRead, entry0.IsRead)
	a.Equal(entries[0].IsBookmarked, entry0.IsBookmarked)
	entry1 := rsp.Entries[1]
	a.Equal(entries[1].ID, entry1.Id)
	a.Equal(entries[1].IsRead, entry1.IsRead)
	a.Equal(entries[1].IsBookmarked, entry1.IsBookmarked)
}

func TestStreamEntriesAllOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, ds := setupServerTest(t)

	req := api.StreamEntriesRequest{FeedId: uint32(8)}
	entries := []*entity.Entry{
		{Title: "Entry 1"},
		{Title: "Entry 2"},
		{Title: "Entry 3"},
	}

	ds.EXPECT().
		ListEntries(gomock.Any(), []entity.ID{8}, nil).
		Return(entries, nil)

	stream, err := client.StreamEntries(context.Background(), &req)
	r.NoError(err)

	var (
		rsp       *api.StreamEntriesResponse
		errStream error
		rsps      = make([]*api.StreamEntriesResponse, 3)
	)
	for i := 0; i < len(rsps); i++ {
		rsp, errStream = stream.Recv()
		a.NoError(errStream)
		a.NotNil(rsp)
		rsps[i] = rsp
	}

	rsp, errStream = stream.Recv()
	a.ErrorIs(errStream, io.EOF)
	a.Nil(rsp)

	// Sort responses so tests are insensitive to input order.
	sort.SliceStable(rsps, func(i, j int) bool {
		return rsps[i].GetEntry().GetTitle() < rsps[j].GetEntry().GetTitle()
	})

	a.Equal("Entry 1", rsps[0].GetEntry().GetTitle())
	a.Equal("Entry 2", rsps[1].GetEntry().GetTitle())
	a.Equal("Entry 3", rsps[2].GetEntry().GetTitle())
}

func TestGetEntryOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	client, ds := setupServerTest(t)

	entry := entity.Entry{
		ID:        2,
		FeedID:    3,
		Title:     "Test Feed Entry",
		IsRead:    false,
		ExtID:     "4abaed90-3435-426f-bf95-05c700a503bf",
		Updated:   pointer(mustTimeVV(t, "2023-07-12T05:02:23.764+02:00")),
		Published: pointer(mustTimeVV(t, "2023-07-12T05:02:23.764+02:00")),
		Content:   pointer("Hello"),
		URL:       pointer("http://x.com/posts/test-feed-entry.html"),
	}

	ds.EXPECT().
		GetEntry(gomock.Any(), entity.ID(2)).
		Return(&entry, nil)

	req := api.GetEntryRequest{Id: 2}

	rsp, err := client.GetEntry(context.Background(), &req)
	r.NoError(err)

	r.NotNil(rsp)
	r.NotNil(rsp.Entry)
	re := rsp.Entry
	a.Equal(entry.ID, re.Id)
	a.Equal(entry.FeedID, re.FeedId)
	a.Equal(entry.Title, re.Title)
	a.Equal(entry.IsRead, re.IsRead)
	a.Equal(entry.ExtID, re.ExtId)
	a.Empty(re.GetDescription())
	a.Equal(*entry.Content, re.GetContent())
	a.Equal(*entry.URL, re.GetUrl())
	// TODO: Also test timestamps.
}

func TestExportOPMLOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, ds := setupServerTest(t)

	ds.EXPECT().
		ExportSubscription(gomock.Any(), nil).
		Return(
			&entity.Subscription{
				Title: pointer("neon export"),
				Feeds: []*entity.Feed{
					{
						Title:     "Feed Q",
						FeedURL:   "http://q.com/feed.xml",
						IsStarred: true,
					},
					{
						Title:   "Feed X",
						FeedURL: "http://x.com/feed.xml",
						Tags:    []string{"foo", "baz"},
					},
					{
						Title:   "Feed A",
						FeedURL: "http://a.com/feed.xml",
					},
				},
			},
			nil,
		)

	req := api.ExportOPMLRequest{}
	rsp, err := client.ExportOPML(context.Background(), &req)
	r.NoError(err)

	a.Regexp(
		regexp.MustCompile(`<\?xml version="1.0" encoding="UTF-8"\?>
<opml version="2.0">
  <head>
    <title>neon export</title>
    <dateCreated>\d+ [A-Z][a-z]+ \d+ \d+:\d+ .+</dateCreated>
  </head>
  <body>
    <outline text="Feed Q" type="rss" xmlUrl="http://q.com/feed.xml" xmlns:neon="https://github.com/bow/neon" neon:isStarred="true"></outline>
    <outline text="Feed X" type="rss" xmlUrl="http://x.com/feed.xml" category="foo,baz"></outline>
    <outline text="Feed A" type="rss" xmlUrl="http://a.com/feed.xml"></outline>
  </body>
</opml>`),
		string(rsp.GetPayload()),
	)
}

func TestImportOPMLOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, ds := setupServerTest(t)

	payload := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head>
    <title>neon export</title>
    <dateCreated>Wed, 09 Nov 2023 04:55:19 CET</dateCreated>
  </head>
  <body>
    <outline text="Feed Q" type="rss" xmlUrl="http://q.com/feed.xml" xmlns:neon="https://github.com/bow/neon" neon:isStarred="true"></outline>
    <outline text="Feed X" type="rss" xmlUrl="http://x.com/feed.xml" category="foo,baz"></outline>
    <outline text="Feed A" type="rss" xmlUrl="http://a.com/feed.xml"></outline>
  </body>
</opml>`)

	sub := entity.Subscription{
		Feeds: []*entity.Feed{
			{
				Title:     "Feed Q",
				FeedURL:   "http://q.com/feed.xml",
				IsStarred: true,
			},
			{
				Title:   "Feed X",
				FeedURL: "http://x.com/feed.xml",
				Tags:    []string{"foo", "baz"},
			},
			{
				Title:   "Feed A",
				FeedURL: "http://a.com/feed.xml",
			},
		},
	}

	ds.EXPECT().
		ImportSubscription(gomock.Any(), &sub).
		Return(3, 2, nil)

	req := api.ImportOPMLRequest{Payload: payload}
	rsp, err := client.ImportOPML(context.Background(), &req)
	r.NoError(err)

	a.Equal(uint32(3), rsp.GetNumProcessed())
	a.Equal(uint32(2), rsp.GetNumImported())
}

func TestGetStatsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, ds := setupServerTest(t)

	stats := entity.Stats{
		NumFeeds:             45,
		NumEntries:           5311,
		NumEntriesUnread:     8,
		LastPullTime:         pointer(mustTimeVV(t, "2023-11-04T05:13:12.805Z")),
		MostRecentUpdateTime: pointer(mustTimeVV(t, "2023-11-04T05:13:12.805Z")),
	}

	ds.EXPECT().
		GetGlobalStats(gomock.Any()).
		Return(&stats, nil)

	req := api.GetStatsRequest{}
	rsp, err := client.GetStats(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	gs := rsp.GetGlobal()
	a.Equal(uint32(45), gs.GetNumFeeds())
	a.Equal(uint32(5311), gs.GetNumEntries())
	a.Equal(uint32(8), gs.GetNumEntriesUnread())
	a.Equal(stats.LastPullTime.Unix(), gs.GetLastPullTime().Seconds)
	a.Equal(stats.MostRecentUpdateTime.Unix(), gs.GetMostRecentUpdateTime().Seconds)
}

func TestGetInfoOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.GetInfoRequest{}
	rsp, err := client.GetInfo(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	want := &api.GetInfoResponse{
		Name:      internal.AppName(),
		Version:   internal.Version(),
		GitCommit: internal.GitCommit(),
	}
	a.Equal(want.Name, rsp.Name)
	a.Equal(want.Version, rsp.Version)
	a.Equal(want.GitCommit, rsp.GitCommit)
}

func pointer[T any](value T) *T { return &value }

func mustTimeVV(t *testing.T, v string) time.Time {
	t.Helper()
	pv, err := time.Parse(time.RFC3339, v)
	require.NoError(t, err)
	return pv
}
