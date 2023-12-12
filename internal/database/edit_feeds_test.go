// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bow/lens/internal"
)

func TestEditFeedsOkEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestDB(t)

	feeds, err := db.EditFeeds(context.Background(), nil)
	r.NoError(err)

	a.Empty(feeds)
}

func TestEditFeedsOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	db := newTestDB(t)

	dbFeeds := []*feedRecord{
		{
			title:     "Feed A",
			feedURL:   "http://a.com/feed.xml",
			updated:   toNullTime(mustTime(t, "2022-03-19T16:23:18.600+02:00")),
			isStarred: false,
		},
	}
	keys := db.addFeeds(dbFeeds)

	r.Equal(1, db.countFeeds())

	existf := func(title string, isStarred bool) bool {
		return db.rowExists(
			`SELECT * FROM feeds WHERE title = ? AND is_starred = ?`,
			title,
			isStarred,
		)
	}

	a.True(existf("Feed A", false))
	a.False(existf("Feed X", true))

	ops := []*internal.FeedEditOp{
		{ID: keys["Feed A"].ID, Title: pointer("Feed X"), IsStarred: pointer(true)},
	}
	feeds, err := db.EditFeeds(context.Background(), ops)
	r.NoError(err)

	a.Len(feeds, 1)

	a.False(existf("Feed A", false))
	a.True(existf("Feed X", true))
}
