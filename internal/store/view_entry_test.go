// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewEntryOk(t *testing.T) {
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
				{Title: "Entry A1", IsRead: true},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: WrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{Title: "Entry X1", IsRead: false},
				{Title: "Entry X2", IsRead: true},
			},
		},
		{
			Title:   "Feed B",
			FeedURL: "http://b.com/feed.xml",
			Updated: WrapNullString("2023-04-09T09:49:22.685+02:00"),
		},
	}
	keys := st.addFeeds(dbFeeds)

	r.Equal(3, st.countFeeds())
	r.Equal(2, st.countEntries(dbFeeds[1].FeedURL))

	dbEntry, err := st.ViewEntry(
		context.Background(),
		keys[dbFeeds[1].Title].Entries["Entry X2"],
	)
	r.NoError(err)
	r.NotNil(dbEntry)

	a.Equal("Entry X2", dbEntry.Title)
	a.True(dbEntry.IsRead)
}

func TestViewEntryErr(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	dbEntry, err := st.ViewEntry(context.Background(), 86)
	r.Nil(dbEntry)
	r.Error(err)

	a.EqualError(err, "SQLite.ViewFeed: entry with ID=86 not found")
}
