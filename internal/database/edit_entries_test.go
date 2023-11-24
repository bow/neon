// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bow/iris/internal"
)

func TestEditEntriesOkEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestDB(t)

	entries, err := db.EditEntries(context.Background(), nil)
	r.NoError(err)

	a.Empty(entries)
}

func TestEditEntriesOkMinimal(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestDB(t)

	dbFeeds := []*feedRecord{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-03-19T16:23:18.600+02:00")),
			entries: []*entryRecord{
				{title: "Entry A1", isRead: true},
			},
		},
	}
	keys := db.addFeeds(dbFeeds)

	r.Equal(1, db.countFeeds())

	existe := func(title string, isRead bool) bool {
		return db.rowExists(
			`SELECT * FROM entries e WHERE e.title = ? AND e.is_read = ?`,
			title,
			isRead,
		)
	}

	a.True(existe("Entry A1", true))
	a.False(existe("Entry A1", false))

	ops := []*internal.EntryEditOp{
		{ID: keys["Feed A"].Entries["Entry A1"], IsRead: pointer(false)},
	}
	entries, err := db.EditEntries(context.Background(), ops)
	r.NoError(err)

	a.Len(entries, 1)

	a.True(existe("Entry A1", false))
	a.False(existe("Entry A1", true))
}

func TestEditEntriesOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestDB(t)

	dbFeeds := []*feedRecord{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-03-19T16:23:18.600+02:00")),
			entries: []*entryRecord{
				{title: "Entry A1", isRead: false, isBookmarked: false},
				{title: "Entry A2", isRead: false, isBookmarked: true},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
			entries: []*entryRecord{
				{title: "Entry X1", isRead: false, isBookmarked: false},
			},
		},
	}
	keys := db.addFeeds(dbFeeds)

	r.Equal(2, db.countFeeds())

	existe := func(title string, isRead, isBookmarked bool) bool {
		return db.rowExists(
			`SELECT * FROM entries e WHERE e.title = ? AND e.is_read = ? AND e.is_bookmarked = ?`,
			title,
			isRead,
			isBookmarked,
		)
	}

	a.True(existe("Entry A1", false, false))
	a.False(existe("Entry A1", true, true))

	a.True(existe("Entry A2", false, true))
	a.False(existe("Entry A2", true, true))

	a.True(existe("Entry X1", false, false))
	a.False(existe("Entry X1", true, true))

	setOps := []*internal.EntryEditOp{
		{ID: keys["Feed X"].Entries["Entry X1"], IsRead: pointer(true), IsBookmarked: pointer(true)},
		{ID: keys["Feed A"].Entries["Entry A2"], IsRead: pointer(true)},
	}
	entries, err := db.EditEntries(context.Background(), setOps)
	r.NoError(err)

	a.Len(entries, 2)

	a.True(existe("Entry A1", false, false))
	a.False(existe("Entry A1", true, true))

	a.False(existe("Entry A2", false, true))
	a.True(existe("Entry A2", true, true))

	a.False(existe("Entry X1", false, false))
	a.True(existe("Entry X1", true, true))
}
