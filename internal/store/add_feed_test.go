package store

import (
	"context"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddFeedOkMinimal(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	feed := gofeed.Feed{
		Title:       "feed-title",
		Description: "feed-description",
		Link:        "https://bar.com",
		FeedLink:    "https://bar.com/feed.xml",
	}

	existf := func() bool {
		return st.rowExists(feedExistSQL, feed.Title, feed.Description, feed.FeedLink, feed.Link)
	}

	a.Equal(0, st.countFeeds())
	a.Equal(0, st.countEntries(feed.FeedLink))
	a.Equal(0, st.countFeedCategories())
	a.False(existf())

	created, err := st.AddFeed(context.Background(), &feed, nil, nil, nil)
	r.NoError(err)

	a.Equal(feed.Title, created.Title)
	a.Equal(feed.Description, created.Description.String)
	a.Equal(feed.Link, created.SiteURL.String)
	a.Equal(feed.FeedLink, created.FeedURL)
	a.Empty([]string(created.Categories))

	a.Equal(1, st.countFeeds())
	a.Equal(0, st.countEntries(feed.FeedLink))
	a.Equal(0, st.countFeedCategories())
	a.True(existf())
}

func TestAddFeedOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

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
	var (
		title       = "user-title"
		description = "user-description"
		categories  = []string{"cat-1", "cat-2", "cat-3"}
	)

	existf1 := func() bool {
		return st.rowExists(feedExistSQL, feed.Title, feed.Description, feed.FeedLink, feed.Link)
	}
	existf2 := func() bool {
		return st.rowExists(feedExistSQL, title, description, feed.FeedLink, feed.Link)
	}
	existe := func(item *gofeed.Item) bool {
		return st.rowExists(feedEntryExistSQL, feed.FeedLink, item.Title, item.Link)
	}

	a.Equal(0, st.countFeeds())
	a.Equal(0, st.countEntries(feed.FeedLink))
	a.Equal(0, st.countFeedCategories())
	a.False(existf1())
	a.False(existf2())
	a.False(existe(feed.Items[0]))
	a.False(existe(feed.Items[1]))

	created, err := st.AddFeed(context.Background(), &feed, &title, &description, categories)
	r.NoError(err)

	a.Equal(title, created.Title)
	a.Equal(description, created.Description.String)
	a.Equal(feed.Link, created.SiteURL.String)
	a.Equal(feed.FeedLink, created.FeedURL)
	a.Equal(categories, []string(created.Categories))

	a.Equal(1, st.countFeeds())
	a.Equal(2, st.countEntries(feed.FeedLink))
	a.Equal(3, st.countFeedCategories())
	a.False(existf1())
	a.True(existf2())
	a.True(existe(feed.Items[0]))
	a.True(existe(feed.Items[1]))
}

func TestAddFeedOkURLExists(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

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
	var categories = []string{"cat-0"}
	st.addFeedWithURL(feed.FeedLink)

	existf := func() bool {
		return st.rowExists(feedExistSQL, feed.Title, feed.Description, feed.FeedLink, feed.Link)
	}
	existe := func(item *gofeed.Item) bool {
		return st.rowExists(feedEntryExistSQL, feed.FeedLink, item.Title, item.Link)
	}

	a.Equal(1, st.countFeeds())
	a.Equal(0, st.countEntries(feed.FeedLink))
	a.Equal(0, st.countFeedCategories())
	a.False(existf())
	a.False(existe(feed.Items[0]))
	a.False(existe(feed.Items[1]))

	created, err := st.AddFeed(context.Background(), &feed, nil, nil, categories)
	r.NoError(err)

	a.Equal(t.Name(), created.Title)
	a.Equal("", created.Description.String)
	a.Equal("", created.SiteURL.String)
	a.Equal(feed.FeedLink, created.FeedURL)
	a.Equal(categories, []string(created.Categories))

	a.Equal(1, st.countFeeds())
	a.Equal(2, st.countEntries(feed.FeedLink))
	a.Equal(1, st.countFeedCategories())
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

// Query for checking that an entry linked to a given feed URL exists.
const feedEntryExistSQL = `
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
