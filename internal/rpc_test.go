package internal

import (
	"context"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/bow/courier/api"
	"github.com/bow/courier/internal/store"
)

func TestAddFeedOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, parser, st := setupServerTest(t)

	req := api.AddFeedRequest{
		Url:         "http://foo.com/feed.xml",
		Title:       stringp("user-title"),
		Description: stringp("user-description"),
		Categories:  []string{"cat-1", "cat-2", "cat-3"},
	}
	feed := gofeed.Feed{
		Title:       "feed-title-original",
		Description: "feed-description-original",
		Link:        "https://foo.com",
		FeedLink:    "https://foo.com/feed.xml",
		Items: []*gofeed.Item{
			{
				GUID:            "entry1",
				Link:            "https://bar.com/entry1.html",
				Title:           "First Entry",
				Content:         "This is the first entry.",
				PublishedParsed: ts(t, "2021-06-18T21:45:26.794+02:00"),
			},
			{
				GUID:            "entry2",
				Link:            "https://bar.com/entry2.html",
				Title:           "Second Entry",
				Content:         "This is the second entry.",
				PublishedParsed: ts(t, "2021-06-18T22:08:16.526+02:00"),
				UpdatedParsed:   ts(t, "2021-06-18T22:11:49.094+02:00"),
			},
		},
	}
	parser.
		EXPECT().
		ParseURLWithContext(req.Url, gomock.Any()).
		MaxTimes(1).
		Return(&feed, nil)

	st.
		EXPECT().
		AddFeed(gomock.Any(), &feed, req.Title, req.Description, req.Categories).
		MaxTimes(1).
		Return(nil)

	rsp, err := client.AddFeed(context.Background(), &req)
	r.NoError(err)

	a.True(proto.Equal(rsp, &api.AddFeedResponse{}))
}

func TestEditFeedOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.EditFeedRequest{}
	rsp, err := client.EditFeed(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestListFeedsOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	client, _, st := setupServerTest(t)

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
	st.
		EXPECT().
		ListFeeds(gomock.Any()).
		MaxTimes(1).
		Return(feeds, nil)

	rsp, err := client.ListFeeds(context.Background(), &req)
	r.NoError(err)

	// TODO: Expand test.
	a.Len(rsp.GetFeeds(), 2)
}

func TestDeleteFeedsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.DeleteFeedsRequest{}
	rsp, err := client.DeleteFeeds(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestPollFeedsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	stream, err := client.PollFeeds(context.Background())
	r.NoError(err)
	waitc := make(chan struct{})

	go func() {
		for {
			rsp, errStream := stream.Recv()
			r.Nil(rsp)
			r.EqualError(errStream, status.New(codes.Unimplemented, "unimplemented").String())
			close(waitc)
			return
		}
	}()

	req := api.PollFeedsRequest{}
	err = stream.Send(&req)
	r.NoError(err)

	err = stream.CloseSend()
	r.NoError(err)
	<-waitc
}

func TestSetEntryFieldsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.SetEntryFieldsRequest{}
	rsp, err := client.SetEntryFields(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestExportOPMLOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.ExportOPMLRequest{}
	rsp, err := client.ExportOPML(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
}

func TestImportOPMLOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	client := newTestClientBuilder(t).Build()

	req := api.ImportOPMLRequest{}
	rsp, err := client.ImportOPML(context.Background(), &req)

	r.Nil(rsp)
	r.EqualError(err, status.New(codes.Unimplemented, "unimplemented").String())
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
		Name:      AppName(),
		Version:   Version(),
		GitCommit: GitCommit(),
		BuildTime: BuildTime(),
	}
	a.Equal(want.Name, rsp.Name)
	a.Equal(want.Version, rsp.Version)
	a.Equal(want.GitCommit, rsp.GitCommit)
	a.Equal(want.BuildTime, rsp.BuildTime)
}

func ts(t *testing.T, value string) *time.Time {
	t.Helper()
	tv, err := store.DeserializeTime(&value)
	require.NoError(t, err)
	return tv
}

func stringp(value string) *string { return &value }
