// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEntriesOkMinimal(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: toNullString("2022-03-19T16:23:18.600+02:00"),
		},
	}
	keys := st.addFeeds(dbFeeds)

	r.Equal(1, st.countFeeds())
	r.Equal(0, st.countEntries(dbFeeds[0].FeedURL))

	entries, err := st.ListEntries(context.Background(), keys[dbFeeds[0].Title].ID)
	r.NoError(err)

	a.Len(entries, 0)
}

func TestListEntriesOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: toNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{Title: "Entry A1", IsRead: true},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: toNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{Title: "Entry X1", IsRead: false},
				{Title: "Entry X2", IsRead: true},
			},
		},
		{
			Title:   "Feed B",
			FeedURL: "http://b.com/feed.xml",
			Updated: toNullString("2023-04-09T09:49:22.685+02:00"),
		},
	}
	keys := st.addFeeds(dbFeeds)

	r.Equal(3, st.countFeeds())
	r.Equal(2, st.countEntries(dbFeeds[1].FeedURL))

	entries, err := st.ListEntries(context.Background(), keys[dbFeeds[1].Title].ID)
	r.NoError(err)

	a.Len(entries, 2)
}

func TestListEntriesErrFeedIDNotFound(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: toNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{Title: "Entry A1", IsRead: true},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: toNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{Title: "Entry X1", IsRead: false},
				{Title: "Entry X2", IsRead: true},
			},
		},
		{
			Title:   "Feed B",
			FeedURL: "http://b.com/feed.xml",
			Updated: toNullString("2023-04-09T09:49:22.685+02:00"),
		},
	}
	st.addFeeds(dbFeeds)

	r.Equal(3, st.countFeeds())
	r.Equal(2, st.countEntries(dbFeeds[1].FeedURL))

	entries, err := st.ListEntries(context.Background(), 404)
	r.Len(entries, 0)

	a.EqualError(err, "SQLite.ListEntries: feed with ID=404 not found")
}
