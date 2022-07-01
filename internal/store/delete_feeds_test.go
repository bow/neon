package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteFeedsEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

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
	st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	err := st.DeleteFeeds(context.Background(), []DBID{})
	r.NoError(err)

	a.Equal(2, st.countFeeds())
}

func TestDeleteFeedsSingle(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: WrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{Title: "Entry A1", IsRead: false},
				{Title: "Entry A2", IsRead: false},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: WrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{Title: "Entry X1", IsRead: false},
			},
		},
	}
	keys := st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())
	r.Equal(2, st.countEntries(dbFeeds[0].FeedURL))
	r.Equal(1, st.countEntries(dbFeeds[1].FeedURL))

	existf := func(title string) bool {
		return st.rowExists(`SELECT * FROM feeds WHERE title = ?`, title)
	}

	a.True(existf("Feed A"))
	a.True(existf("Feed X"))

	err := st.DeleteFeeds(context.Background(), []DBID{keys["Feed X"].DBID})
	r.NoError(err)
	r.Equal(1, st.countFeeds())
	r.Equal(2, st.countEntries(dbFeeds[0].FeedURL))
	r.Equal(0, st.countEntries(dbFeeds[1].FeedURL))

	a.True(existf("Feed A"))
	a.False(existf("Feed X"))
}
