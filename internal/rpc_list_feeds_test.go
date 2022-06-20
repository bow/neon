package internal

import (
	"context"
	"testing"

	"github.com/bow/courier/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFeedsOkEmpty(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	req := api.ListFeedsRequest{}

	client := newTestClientBuilder(t).Build()

	rsp, err := client.ListFeeds(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	a.Empty(rsp.GetFeeds())
}

func TestListFeedsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	req := api.ListFeedsRequest{}

	cbuilder := newTestClientBuilder(t)
	db := newTestDB(t, cbuilder.serverBuilder.storePath)
	client := cbuilder.Build()

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: WrapNullString("2022-03-19T16:23:18.600+02:00"),
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: WrapNullString("2022-04-20T16:32:30.760+02:00"),
		},
	}
	db.addFeeds(dbFeeds)

	r.Equal(2, db.countFeeds())

	rsp, err := client.ListFeeds(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	feeds := rsp.GetFeeds()
	r.Len(feeds, 2)

	feed0 := feeds[0]
	a.Equal(feed0.GetFeedUrl(), dbFeeds[1].FeedURL)

	feed1 := feeds[1]
	a.Equal(feed1.GetFeedUrl(), dbFeeds[0].FeedURL)
}
