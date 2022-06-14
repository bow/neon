package internal

import (
	"context"
	"testing"

	"github.com/bow/courier/api"
	gomock "github.com/golang/mock/gomock"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddFeedOkMinimal(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	url := "http://bar.com/feed.xml"
	feed := gofeed.Feed{
		Title:       "feed-title",
		Description: "feed-description",
		Link:        "https://bar.com",
		FeedLink:    "https://bar.com/feed.xml",
	}
	parser := NewMockFeedParser(gomock.NewController(t))
	parser.
		EXPECT().
		ParseURL(url).
		MaxTimes(1).
		Return(&feed, nil)

	client, db := setupOfflineTest(t, parser)

	existf := func() bool {
		sql := `SELECT * FROM feeds WHERE xml_url = ?`
		return db.rowExists(sql, feed.FeedLink)
	}

	a.Equal(0, db.countFeeds())
	a.False(existf())

	req := api.AddFeedRequest{Url: url}
	rsp, err := client.AddFeed(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	a.Equal(1, db.countFeeds())
	a.True(existf())
}
