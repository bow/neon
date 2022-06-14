package internal

import (
	"context"
	"path/filepath"
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

	storePath := filepath.Join(t.TempDir(), t.Name()+".db")
	r.NoFileExists(storePath)

	server := defaultTestServerBuilder(t).Parser(parser).StorePath(storePath)
	client := newTestClientBuilder().ServerBuilder(server).Build(t)
	r.FileExists(storePath)

	ctx := context.Background()
	db := newTestDB(t, storePath)

	existf := func() bool {
		sql := `SELECT * FROM feeds WHERE xml_url = ?`
		return db.rowExists(sql, feed.FeedLink)
	}

	preFeedCount := db.countFeeds()
	a.Equal(0, preFeedCount)
	a.False(existf())

	req := api.AddFeedRequest{Url: url}
	rsp, err := client.AddFeed(ctx, &req)
	r.NoError(err)
	r.NotNil(rsp)

	postFeedCount := db.countFeeds()
	a.Equal(1, postFeedCount)
	a.True(existf())
}
