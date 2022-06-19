package internal

import (
	"context"
	"testing"
	"time"

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
		ParseURLWithContext(req.Url, gomock.Any()).
		MaxTimes(1).
		Return(&feed, nil)

	client, db := setupOfflineTest(t, parser)

	existf := func() bool {
		return db.rowExists(feedExistSQL, feed.Title, feed.Description, feed.FeedLink, feed.Link)
	}

	a.Equal(0, db.countFeeds())
	a.Equal(0, db.countEntries(feed.FeedLink))
	a.Equal(0, db.countFeedCategories())
	a.False(existf())

	rsp, err := client.AddFeed(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	a.Equal(1, db.countFeeds())
	a.Equal(0, db.countEntries(feed.FeedLink))
	a.Equal(0, db.countFeedCategories())
	a.True(existf())
}

func TestAddFeedOkExtended(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

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
	parser := NewMockFeedParser(gomock.NewController(t))
	parser.
		EXPECT().
		ParseURLWithContext(req.Url, gomock.Any()).
		MaxTimes(1).
		Return(&feed, nil)

	client, db := setupOfflineTest(t, parser)

	existf1 := func() bool {
		return db.rowExists(feedExistSQL, feed.Title, feed.Description, feed.FeedLink, feed.Link)
	}
	existf2 := func() bool {
		return db.rowExists(feedExistSQL, req.Title, req.Description, feed.FeedLink, feed.Link)
	}
	existe := func(item *gofeed.Item) bool {
		return db.rowExists(entryExistSQL, feed.FeedLink, item.Title, item.Link)
	}

	a.Equal(0, db.countFeeds())
	a.Equal(0, db.countEntries(feed.FeedLink))
	a.Equal(0, db.countFeedCategories())
	a.False(existf1())
	a.False(existf2())
	a.False(existe(feed.Items[0]))
	a.False(existe(feed.Items[1]))

	rsp, err := client.AddFeed(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	a.Equal(1, db.countFeeds())
	a.Equal(2, db.countEntries(feed.FeedLink))
	a.Equal(3, db.countFeedCategories())
	a.False(existf1())
	a.True(existf2())
	a.True(existe(feed.Items[0]))
	a.True(existe(feed.Items[1]))
}

func TestAddFeedOkFeedExists(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	req := api.AddFeedRequest{Url: "http://qux.net/feed.xml", Categories: []string{"cat-0"}}
	feed := gofeed.Feed{
		Title:       "feed-title",
		Description: "feed-description",
		Link:        "https://bar.com",
		FeedLink:    "https://bar.com/feed.xml",
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
	parser := NewMockFeedParser(gomock.NewController(t))
	parser.
		EXPECT().
		ParseURLWithContext(req.Url, gomock.Any()).
		MaxTimes(1).
		Return(&feed, nil)

	client, db := setupOfflineTest(t, parser)
	db.addFeedWithURL(feed.FeedLink)

	existf := func() bool {
		return db.rowExists(feedExistSQL, feed.Title, feed.Description, feed.FeedLink, feed.Link)
	}
	existe := func(item *gofeed.Item) bool {
		return db.rowExists(entryExistSQL, feed.FeedLink, item.Title, item.Link)
	}

	a.Equal(1, db.countFeeds())
	a.Equal(0, db.countEntries(feed.FeedLink))
	a.Equal(0, db.countFeedCategories())
	a.False(existf())
	a.False(existe(feed.Items[0]))
	a.False(existe(feed.Items[1]))

	rsp, err := client.AddFeed(context.Background(), &req)
	r.NoError(err)
	r.NotNil(rsp)

	a.Equal(1, db.countFeeds())
	a.Equal(2, db.countEntries(feed.FeedLink))
	a.Equal(1, db.countFeedCategories())
	a.False(existf())
	a.True(existe(feed.Items[0]))
	a.True(existe(feed.Items[1]))
}

// Query for checking that a feed exists.
const feedExistSQL = `
	SELECT
		*
	FROM
		feeds
	WHERE
		title = ?
		AND description = ?
		AND feed_url = ?
		AND site_url = ?
`

// Query for checking that an entry exists.
const entryExistSQL = `
	SELECT
		*
	FROM
		entries e
		INNER JOIN feeds f ON e.feed_id = f.id
	WHERE
		f.feed_url = ?
		AND e.title = ?
		AND e.url = ?
`

func ts(t *testing.T, value string) *time.Time {
	t.Helper()
	rtv, err := time.Parse(time.RFC3339Nano, value)
	require.NoError(t, err)
	tv := rtv.UTC()
	return &tv
}

func stringp(value string) *string { return &value }
