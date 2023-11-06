// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteFeedsOkEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: wrapNullString("2022-03-19T16:23:18.600+02:00"),
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: wrapNullString("2022-04-20T16:32:30.760+02:00"),
		},
	}
	st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	err := st.DeleteFeeds(context.Background(), []ID{})
	r.NoError(err)

	a.Equal(2, st.countFeeds())
}

func TestDeleteFeedsOkSingle(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: wrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{Title: "Entry A1"},
				{Title: "Entry A2"},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: wrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{Title: "Entry X1"},
			},
		},
	}
	keys := st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())
	a.Equal(2, st.countEntries(dbFeeds[0].FeedURL))
	a.Equal(1, st.countEntries(dbFeeds[1].FeedURL))

	existf := func(title string) bool {
		return st.rowExists(`SELECT * FROM feeds WHERE title = ?`, title)
	}

	a.True(existf("Feed A"))
	a.True(existf("Feed X"))

	err := st.DeleteFeeds(context.Background(), []ID{keys["Feed X"].ID})
	r.NoError(err)
	a.Equal(1, st.countFeeds())
	a.Equal(2, st.countEntries(dbFeeds[0].FeedURL))
	a.Equal(0, st.countEntries(dbFeeds[1].FeedURL))

	a.True(existf("Feed A"))
	a.False(existf("Feed X"))
}

func TestDeleteFeedsOkMultiple(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: wrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{Title: "Entry A1"},
				{Title: "Entry A2"},
			},
		},
		{
			Title:   "Feed P",
			FeedURL: "http://p.com/feed.xml",
			Updated: wrapNullString("2022-04-02T10:16:00.471+02:00"),
			Entries: []*Entry{
				{Title: "Entry P5"},
				{Title: "Entry P6"},
				{Title: "Entry P7"},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: wrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{Title: "Entry X1"},
			},
		},
	}
	keys := st.addFeeds(dbFeeds)
	r.Equal(3, st.countFeeds())
	a.Equal(2, st.countEntries(dbFeeds[0].FeedURL))
	a.Equal(3, st.countEntries(dbFeeds[1].FeedURL))
	a.Equal(1, st.countEntries(dbFeeds[2].FeedURL))

	existf := func(title string) bool {
		return st.rowExists(`SELECT * FROM feeds WHERE title = ?`, title)
	}

	a.True(existf("Feed A"))
	a.True(existf("Feed P"))
	a.True(existf("Feed X"))

	err := st.DeleteFeeds(context.Background(), []ID{keys["Feed A"].ID, keys["Feed P"].ID})
	r.NoError(err)
	a.Equal(1, st.countFeeds())
	a.Equal(0, st.countEntries(dbFeeds[0].FeedURL))
	a.Equal(0, st.countEntries(dbFeeds[1].FeedURL))
	a.Equal(1, st.countEntries(dbFeeds[2].FeedURL))

	a.False(existf("Feed A"))
	a.False(existf("Feed P"))
	a.True(existf("Feed X"))
}

func TestDeleteFeedsErrHasMissing(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: wrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{Title: "Entry A1"},
				{Title: "Entry A2"},
			},
		},
		{
			Title:   "Feed P",
			FeedURL: "http://p.com/feed.xml",
			Updated: wrapNullString("2022-04-02T10:16:00.471+02:00"),
			Entries: []*Entry{
				{Title: "Entry P5"},
				{Title: "Entry P6"},
				{Title: "Entry P7"},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: wrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{Title: "Entry X1"},
			},
		},
	}
	keys := st.addFeeds(dbFeeds)
	r.Equal(3, st.countFeeds())
	a.Equal(2, st.countEntries(dbFeeds[0].FeedURL))
	a.Equal(3, st.countEntries(dbFeeds[1].FeedURL))
	a.Equal(1, st.countEntries(dbFeeds[2].FeedURL))

	existf := func(title string) bool {
		return st.rowExists(`SELECT * FROM feeds WHERE title = ?`, title)
	}

	a.True(existf("Feed A"))
	a.True(existf("Feed P"))
	a.True(existf("Feed X"))

	err := st.DeleteFeeds(context.Background(), []ID{keys["Feed A"].ID, 99})
	a.EqualError(err, "SQLite.DeleteFeeds: feed with ID=99 not found")

	r.Equal(3, st.countFeeds())
	a.Equal(2, st.countEntries(dbFeeds[0].FeedURL))
	a.Equal(3, st.countEntries(dbFeeds[1].FeedURL))
	a.Equal(1, st.countEntries(dbFeeds[2].FeedURL))

	a.True(existf("Feed A"))
	a.True(existf("Feed P"))
	a.True(existf("Feed X"))
}
