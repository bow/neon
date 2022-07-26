package server

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/bow/courier/api"
	"github.com/bow/courier/internal"
	"github.com/bow/courier/internal/store"
)

func TestAddFeedOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, str := setupServerTest(t)

	req := api.AddFeedRequest{
		Url:         "http://foo.com/feed.xml",
		Title:       pointer("user-title"),
		Description: pointer("user-description"),
		Tags:        []string{"tag-1", "tag-2", "tag-3"},
		IsStarred:   pointer(true),
	}
	created := store.Feed{
		Title:       "feed-title-original",
		Description: store.WrapNullString("feed-description-original"),
		SiteURL:     store.WrapNullString("https://foo.com"),
		FeedURL:     "https://foo.com/feed.xml",
		Subscribed:  "2021-07-01T23:33:06.156+02:00",
		IsStarred:   true,
	}

	str.EXPECT().
		AddFeed(
			gomock.Any(),
			req.GetUrl(),
			req.Title,
			req.Description,
			req.Tags,
			req.IsStarred,
		).
		Return(&created, nil)

	rsp, err := client.AddFeed(context.Background(), &req)
	r.NoError(err)

	a.Equal(created.Title, rsp.Feed.Title)
	a.Equal(created.Description.String, *rsp.Feed.Description)
	a.Equal(created.SiteURL.String, *rsp.Feed.SiteUrl)
	a.Equal(created.FeedURL, rsp.Feed.FeedUrl)
	a.Equal(created.IsStarred, rsp.Feed.IsStarred)
}

func TestListFeedsOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, str := setupServerTest(t)

	req := api.ListFeedsRequest{}
	feeds := []*store.Feed{
		{
			Title:      "Feed A",
			FeedURL:    "http://a.com/feed.xml",
			Subscribed: "2022-06-22T19:39:38.964+02:00",
			Updated:    store.WrapNullString("2022-03-19T16:23:18.600+02:00"),
		},
		{
			Title:      "Feed X",
			FeedURL:    "http://x.com/feed.xml",
			Subscribed: "2022-06-22T19:39:44.037+02:00",
			Updated:    store.WrapNullString("2022-04-20T16:32:30.760+02:00"),
		},
	}
	str.EXPECT().
		ListFeeds(gomock.Any()).
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
	client, str := setupServerTest(t)

	ops := []*store.FeedEditOp{
		{DBID: 14, Title: pointer("newer")},
		{DBID: 58, Tags: pointer([]string{"x", "y"})},
		{DBID: 77, IsStarred: pointer(true)},
	}
	feeds := []*store.Feed{
		{DBID: 14, Title: "newer", Subscribed: "2022-06-30T00:53:50.200+02:00"},
		{DBID: 58, Tags: []string{"x", "y"}, Subscribed: "2022-06-30T00:53:58.135+02:00"},
		{DBID: 77, IsStarred: true, Subscribed: "2022-06-30T00:53:59.812+02:00"},
	}

	str.EXPECT().
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
	a.Equal(int32(feeds[0].DBID), feed0.Id)
	a.Equal(feeds[0].Title, feed0.Title)
	feed1 := rsp.Feeds[1]
	a.Equal(int32(feeds[1].DBID), feed1.Id)
	a.Equal([]string(feeds[1].Tags), feed1.Tags)
	feed2 := rsp.Feeds[2]
	a.Equal(int32(feeds[2].DBID), feed2.Id)
	a.Equal(feeds[2].IsStarred, feed2.IsStarred)
}

func TestDeleteFeedsOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, str := setupServerTest(t)

	str.EXPECT().
		DeleteFeeds(gomock.Any(), []store.DBID{1, 9}).
		Return(nil)

	req := api.DeleteFeedsRequest{FeedIds: []int32{1, 9}}
	rsp, err := client.DeleteFeeds(context.Background(), &req)
	r.NoError(err)

	a.True(proto.Equal(&api.DeleteFeedsResponse{}, rsp))
}

func TestDeleteFeedsErrNotFound(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, str := setupServerTest(t)

	str.EXPECT().
		DeleteFeeds(gomock.Any(), []store.DBID{1, 9}).
		Return(fmt.Errorf("wrapped: %w", store.FeedNotFoundError{ID: 9}))

	req := api.DeleteFeedsRequest{FeedIds: []int32{1, 9}}
	rsp, err := client.DeleteFeeds(context.Background(), &req)

	r.Nil(rsp)
	a.EqualError(err, "rpc error: code = NotFound desc = feed with ID=9 not found")
}

func TestPullFeedsAllOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, str := setupServerTest(t)

	prs := []store.PullResult{
		store.NewPullResultFromFeed(
			pointer("http://a.com/feed.xml"),
			&store.Feed{
				Title:      "feed-A",
				FeedURL:    "https://a.com/feed.xml",
				Subscribed: "2021-07-23T17:20:29.499+02:00",
				IsStarred:  true,
				Entries: []*store.Entry{
					{Title: "Entry A1", IsRead: false},
					{Title: "Entry A2", IsRead: false},
				},
			},
		),
		store.NewPullResultFromFeed(pointer("http://z.com/feed.xml"), nil),
		store.NewPullResultFromFeed(
			pointer("http://c.com/feed.xml"),
			&store.Feed{
				Title:      "feed-C",
				FeedURL:    "https://c.com/feed.xml",
				Subscribed: "2021-07-23T17:21:11.489+02:00",
				IsStarred:  false,
				Entries: []*store.Entry{
					{Title: "Entry C3", IsRead: false},
				},
			},
		),
	}

	ch := make(chan store.PullResult)
	go func() {
		defer close(ch)

		// Randomize ordering, to simulate actual URL pulls.
		shufres := make([]store.PullResult, len(prs))
		copy(shufres, prs)
		rand.Seed(time.Now().UnixNano())
		shf := func(i, j int) { shufres[i], shufres[j] = shufres[j], shufres[i] }
		rand.Shuffle(len(shufres), shf)

		for i := 0; i < len(shufres); i++ {
			ch <- shufres[i]
		}
	}()

	str.EXPECT().
		PullFeeds(gomock.Any()).
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

func TestPullFeedsErrSomeFeed(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, str := setupServerTest(t)

	prs := []store.PullResult{
		store.NewPullResultFromFeed(
			pointer("https://a.com/feed.xml"),
			&store.Feed{
				Title:      "feed-A",
				FeedURL:    "https://a.com/feed.xml",
				Subscribed: "2021-07-23T17:20:29.499+02:00",
				IsStarred:  true,
				Entries: []*store.Entry{
					{Title: "Entry A1", IsRead: false},
					{Title: "Entry A2", IsRead: false},
				},
			},
		),
		store.NewPullResultFromError(pointer("https://x.com/feed.xml"), fmt.Errorf("timed out")),
		store.NewPullResultFromFeed(pointer("https://z.com/feed.xml"), nil),
		store.NewPullResultFromFeed(
			pointer("https://c.com/feed.xml"),
			&store.Feed{
				Title:      "feed-C",
				FeedURL:    "https://c.com/feed.xml",
				Subscribed: "2021-07-23T17:21:11.489+02:00",
				IsStarred:  false,
				Entries: []*store.Entry{
					{Title: "Entry C3", IsRead: false},
				},
			},
		),
	}

	ch := make(chan store.PullResult)
	go func() {
		defer close(ch)

		// Randomize ordering, to simulate actual URL pulls.
		shufres := make([]store.PullResult, len(prs))
		copy(shufres, prs)
		rand.Seed(time.Now().UnixNano())
		shf := func(i, j int) { shufres[i], shufres[j] = shufres[j], shufres[i] }
		rand.Shuffle(len(shufres), shf)

		for i := 0; i < len(shufres); i++ {
			ch <- shufres[i]
		}
	}()

	str.EXPECT().
		PullFeeds(gomock.Any()).
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

func TestEditEntriesOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, str := setupServerTest(t)

	ops := []*store.EntryEditOp{
		{DBID: 37, IsRead: pointer(true)},
		{DBID: 49, IsRead: pointer(false)},
	}
	entries := []*store.Entry{
		{DBID: 37, IsRead: true},
		{DBID: 49, IsRead: false},
	}

	str.EXPECT().
		EditEntries(gomock.Any(), ops).
		Return(entries, nil)

	req := api.EditEntriesRequest{
		Ops: []*api.EditEntriesRequest_Op{
			{
				Id: 37,
				Fields: &api.EditEntriesRequest_Op_Fields{
					IsRead: pointer(true),
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
	a.Equal(int32(entries[0].DBID), entry0.Id)
	a.Equal(entries[0].IsRead, entry0.IsRead)
	entry1 := rsp.Entries[1]
	a.Equal(int32(entries[1].DBID), entry1.Id)
	a.Equal(entries[1].IsRead, entry1.IsRead)
}

func TestExportOPMLOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, str := setupServerTest(t)

	payload := `<\?xml version="1.0" encoding="UTF-8"\?>
<opml version="2.0">
  <head>
    <title>Courier export</title>
	<dateCreated>Thu, 17 Feb 2022 16:37:19 CET</dateCreated>
  </head>
  <body>
    <outline text="Feed X" type="rss" xmlUrl="http://x.com/feed.xml" category="foo,baz"></outline>
  </body>
</opml>`
	str.EXPECT().
		ExportOPML(gomock.Any(), nil).
		Return([]byte(payload), nil)

	req := api.ExportOPMLRequest{}
	rsp, err := client.ExportOPML(context.Background(), &req)
	r.NoError(err)

	a.Equal([]byte(payload), rsp.GetPayload())
}

func TestImportOPMLOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	client, str := setupServerTest(t)

	payload := []byte("payload")
	str.EXPECT().
		ImportOPML(gomock.Any(), payload).
		Return(3, nil)

	req := api.ImportOPMLRequest{Payload: payload}
	rsp, err := client.ImportOPML(context.Background(), &req)
	r.NoError(err)

	a.Equal(int32(3), rsp.GetNumImported())
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
		BuildTime: internal.BuildTime(),
	}
	a.Equal(want.Name, rsp.Name)
	a.Equal(want.Version, rsp.Version)
	a.Equal(want.GitCommit, rsp.GitCommit)
	a.Equal(want.BuildTime, rsp.BuildTime)
}

func pointer[T any](value T) *T { return &value }
