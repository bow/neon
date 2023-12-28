// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEntryOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestSQLiteDB(t)

	dbFeeds := []*feedRecord{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-03-19T16:23:18.600+02:00")),
			entries: []*entryRecord{
				{title: "Entry A1", isRead: true},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
			entries: []*entryRecord{
				{title: "Entry X1", isRead: false},
				{title: "Entry X2", isRead: true},
			},
		},
		{
			title:   "Feed B",
			feedURL: "http://b.com/feed.xml",
			updated: toNullTime(mustTime(t, "2023-04-09T09:49:22.685+02:00")),
		},
	}
	keys := db.addFeeds(dbFeeds)

	r.Equal(3, db.countFeeds())
	r.Equal(2, db.countEntries(dbFeeds[1].feedURL))

	dbEntry, err := db.GetEntry(
		context.Background(),
		keys[dbFeeds[1].title].Entries["Entry X2"],
	)
	r.NoError(err)
	r.NotNil(dbEntry)

	a.Equal("Entry X2", dbEntry.Title)
	a.True(dbEntry.IsRead)
}

func TestGetEntryErr(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestSQLiteDB(t)

	r.Equal(0, db.countFeeds())

	dbEntry, err := db.GetEntry(context.Background(), 86)
	r.Nil(dbEntry)
	r.Error(err)

	a.EqualError(err, "SQLite.GetEntry: entry with ID=86 not found")
}
