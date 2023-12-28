// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFeedsOkMinimal(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestSQLiteDB(t)

	feeds, err := db.ListFeeds(context.Background())
	r.NoError(err)

	a.Empty(feeds)
}

func TestListFeedsOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestSQLiteDB(t)

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

	feeds, err := db.ListFeeds(context.Background())
	r.NoError(err)
	r.NotEmpty(feeds)

	a.Len(feeds, 2)

	feed0 := feeds[0]
	a.Equal(feed0.FeedURL, dbFeeds[1].feedURL)

	feed1 := feeds[1]
	a.Equal(feed1.FeedURL, dbFeeds[0].feedURL)
}
