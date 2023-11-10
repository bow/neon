// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

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
	db := newTestDB(t)

	dbFeeds := []*feedRecord{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-03-19T16:23:18.600+02:00")),
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
		},
	}
	db.addFeeds(dbFeeds)
	r.Equal(2, db.countFeeds())

	err := db.DeleteFeeds(context.Background(), []ID{})
	r.NoError(err)

	a.Equal(2, db.countFeeds())
}

func TestDeleteFeedsOkSingle(t *testing.T) {
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
				{title: "Entry A1"},
				{title: "Entry A2"},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
			entries: []*entryRecord{
				{title: "Entry X1"},
			},
		},
	}
	keys := db.addFeeds(dbFeeds)
	r.Equal(2, db.countFeeds())
	a.Equal(2, db.countEntries(dbFeeds[0].feedURL))
	a.Equal(1, db.countEntries(dbFeeds[1].feedURL))

	existf := func(title string) bool {
		return db.rowExists(`SELECT * FROM feeds WHERE title = ?`, title)
	}

	a.True(existf("Feed A"))
	a.True(existf("Feed X"))

	err := db.DeleteFeeds(context.Background(), []ID{keys["Feed X"].ID})
	r.NoError(err)
	a.Equal(1, db.countFeeds())
	a.Equal(2, db.countEntries(dbFeeds[0].feedURL))
	a.Equal(0, db.countEntries(dbFeeds[1].feedURL))

	a.True(existf("Feed A"))
	a.False(existf("Feed X"))
}

func TestDeleteFeedsOkMultiple(t *testing.T) {
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
				{title: "Entry A1"},
				{title: "Entry A2"},
			},
		},
		{
			title:   "Feed P",
			feedURL: "http://p.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-02T10:16:00.471+02:00")),
			entries: []*entryRecord{
				{title: "Entry P5"},
				{title: "Entry P6"},
				{title: "Entry P7"},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
			entries: []*entryRecord{
				{title: "Entry X1"},
			},
		},
	}
	keys := db.addFeeds(dbFeeds)
	r.Equal(3, db.countFeeds())
	a.Equal(2, db.countEntries(dbFeeds[0].feedURL))
	a.Equal(3, db.countEntries(dbFeeds[1].feedURL))
	a.Equal(1, db.countEntries(dbFeeds[2].feedURL))

	existf := func(title string) bool {
		return db.rowExists(`SELECT * FROM feeds WHERE title = ?`, title)
	}

	a.True(existf("Feed A"))
	a.True(existf("Feed P"))
	a.True(existf("Feed X"))

	err := db.DeleteFeeds(context.Background(), []ID{keys["Feed A"].ID, keys["Feed P"].ID})
	r.NoError(err)
	a.Equal(1, db.countFeeds())
	a.Equal(0, db.countEntries(dbFeeds[0].feedURL))
	a.Equal(0, db.countEntries(dbFeeds[1].feedURL))
	a.Equal(1, db.countEntries(dbFeeds[2].feedURL))

	a.False(existf("Feed A"))
	a.False(existf("Feed P"))
	a.True(existf("Feed X"))
}

func TestDeleteFeedsErrHasMissing(t *testing.T) {
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
				{title: "Entry A1"},
				{title: "Entry A2"},
			},
		},
		{
			title:   "Feed P",
			feedURL: "http://p.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-02T10:16:00.471+02:00")),
			entries: []*entryRecord{
				{title: "Entry P5"},
				{title: "Entry P6"},
				{title: "Entry P7"},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
			entries: []*entryRecord{
				{title: "Entry X1"},
			},
		},
	}
	keys := db.addFeeds(dbFeeds)
	r.Equal(3, db.countFeeds())
	a.Equal(2, db.countEntries(dbFeeds[0].feedURL))
	a.Equal(3, db.countEntries(dbFeeds[1].feedURL))
	a.Equal(1, db.countEntries(dbFeeds[2].feedURL))

	existf := func(title string) bool {
		return db.rowExists(`SELECT * FROM feeds WHERE title = ?`, title)
	}

	a.True(existf("Feed A"))
	a.True(existf("Feed P"))
	a.True(existf("Feed X"))

	err := db.DeleteFeeds(context.Background(), []ID{keys["Feed A"].ID, 99})
	a.EqualError(err, "SQLite.DeleteFeeds: feed with ID=99 not found")

	r.Equal(3, db.countFeeds())
	a.Equal(2, db.countEntries(dbFeeds[0].feedURL))
	a.Equal(3, db.countEntries(dbFeeds[1].feedURL))
	a.Equal(1, db.countEntries(dbFeeds[2].feedURL))

	a.True(existf("Feed A"))
	a.True(existf("Feed P"))
	a.True(existf("Feed X"))
}
