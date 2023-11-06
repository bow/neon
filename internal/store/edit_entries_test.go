// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditEntriesOkEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	entries, err := st.EditEntries(context.Background(), nil)
	r.NoError(err)

	a.Empty(entries)
}

func TestEditEntriesOkMinimal(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*FeedRecord{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullString("2022-03-19T16:23:18.600+02:00"),
			entries: []*Entry{
				{Title: "Entry A1", IsRead: true},
			},
		},
	}
	keys := st.addFeeds(dbFeeds)

	r.Equal(1, st.countFeeds())

	existe := func(title string, isRead bool) bool {
		return st.rowExists(
			`SELECT * FROM entries e WHERE e.title = ? AND e.is_read = ?`,
			title,
			isRead,
		)
	}

	a.True(existe("Entry A1", true))
	a.False(existe("Entry A1", false))

	ops := []*EntryEditOp{
		{ID: keys["Feed A"].Entries["Entry A1"], IsRead: pointer(false)},
	}
	entries, err := st.EditEntries(context.Background(), ops)
	r.NoError(err)

	a.Len(entries, 1)

	a.True(existe("Entry A1", false))
	a.False(existe("Entry A1", true))
}

func TestEditEntriesOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*FeedRecord{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullString("2022-03-19T16:23:18.600+02:00"),
			entries: []*Entry{
				{Title: "Entry A1", IsRead: false},
				{Title: "Entry A2", IsRead: false},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullString("2022-04-20T16:32:30.760+02:00"),
			entries: []*Entry{
				{Title: "Entry X1", IsRead: false},
			},
		},
	}
	keys := st.addFeeds(dbFeeds)

	r.Equal(2, st.countFeeds())

	existe := func(title string, isRead bool) bool {
		return st.rowExists(
			`SELECT * FROM entries e WHERE e.title = ? AND e.is_read = ?`,
			title,
			isRead,
		)
	}

	a.True(existe("Entry A1", false))
	a.False(existe("Entry A1", true))

	a.True(existe("Entry A2", false))
	a.False(existe("Entry A2", true))

	a.True(existe("Entry X1", false))
	a.False(existe("Entry X1", true))

	setOps := []*EntryEditOp{
		{ID: keys["Feed X"].Entries["Entry X1"], IsRead: pointer(true)},
		{ID: keys["Feed A"].Entries["Entry A2"], IsRead: pointer(true)},
	}
	entries, err := st.EditEntries(context.Background(), setOps)
	r.NoError(err)

	a.Len(entries, 2)

	a.True(existe("Entry A1", false))
	a.False(existe("Entry A1", true))

	a.False(existe("Entry A2", false))
	a.True(existe("Entry A2", true))

	a.False(existe("Entry X1", false))
	a.True(existe("Entry X1", true))
}
