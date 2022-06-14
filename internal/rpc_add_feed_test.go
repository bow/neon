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

	req := api.AddFeedRequest{Url: "http://bar.com/feed.xml"}
	feed := gofeed.Feed{
		Title:       "feed-title",
		Description: "feed-description",
		Link:        "https://bar.com",
		FeedLink:    "https://bar.com/feed.xml",
	}
	parser := NewMockFeedParser(gomock.NewController(t))
	parser.
		EXPECT().
		ParseURL(req.Url).
		MaxTimes(1).
		Return(&feed, nil)

	client, db := setupOfflineTest(t, parser)

	existf := func() bool {
		return db.rowExists(feedExistSQL, feed.Title, feed.Description, feed.FeedLink, feed.Link)
	}

	a.Equal(0, db.countFeeds())
	a.Equal(0, db.countFeedCategories())
	a.False(existf())

	rsp, err := client.AddFeed(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	a.Equal(1, db.countFeeds())
	a.Equal(0, db.countFeedCategories())
	a.True(existf())
}

// Query for checking that a feed row with the given columns exist.
const feedExistSQL = `
	SELECT
		*
	FROM
		feeds
	WHERE
		title = ?
		AND description = ?
		AND xml_url = ?
		AND html_url = ?
`
