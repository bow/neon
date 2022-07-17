package store

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullFeedsOkEmptyDB(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	st.parser.EXPECT().
		ParseURLWithContext(gomock.Any(), gomock.Any()).
		MaxTimes(0)

	c := st.PullFeeds(context.Background())
	a.Empty(c)
}

func TestPullFeedsOkEmptyEntries(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: WrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: WrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{},
		},
	}

	keys := st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[0].FeedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, dbFeeds[0]), nil)

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[1].FeedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, dbFeeds[1]), nil)

	c := st.PullFeeds(context.Background())

	got := make([]PullResult, 0)
	for res := range c {
		got = append(got, res)
	}

	want := []PullResult{
		{
			pk:     pullKey{feedDBID: keys["Feed A"].DBID, feedURL: dbFeeds[0].FeedURL},
			status: pullSuccess,
			ok:     nil,
			err:    nil,
		},
		{
			pk:     pullKey{feedDBID: keys["Feed X"].DBID, feedURL: dbFeeds[1].FeedURL},
			status: pullSuccess,
			ok:     nil,
			err:    nil,
		},
	}

	a.ElementsMatch(want, got)
}

func toGFeed(t *testing.T, feed *Feed) *gofeed.Feed {
	t.Helper()
	gfeed := gofeed.Feed{
		Title:    feed.Title,
		FeedLink: feed.FeedURL,
		Updated:  feed.Updated.String,
	}
	for _, entry := range feed.Entries {
		item := gofeed.Item{
			GUID:      entry.ExtID,
			Link:      entry.URL.String,
			Title:     entry.Title,
			Content:   entry.Content.String,
			Published: entry.Published.String,
			Updated:   entry.Updated.String,
		}
		gfeed.Items = append(gfeed.Items, &item)
	}
	return &gfeed
}
