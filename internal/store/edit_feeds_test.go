// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditFeedsOkEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	feeds, err := st.EditFeeds(context.Background(), nil)
	r.NoError(err)

	a.Empty(feeds)
}

func TestEditFeedsOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:     "Feed A",
			FeedURL:   "http://a.com/feed.xml",
			Updated:   wrapNullString("2022-03-19T16:23:18.600+02:00"),
			IsStarred: false,
		},
	}
	keys := st.addFeeds(dbFeeds)

	r.Equal(1, st.countFeeds())

	existf := func(title string, isStarred bool) bool {
		return st.rowExists(
			`SELECT * FROM feeds WHERE title = ? AND is_starred = ?`,
			title,
			isStarred,
		)
	}

	a.True(existf("Feed A", false))
	a.False(existf("Feed X", true))

	ops := []*FeedEditOp{
		{ID: keys["Feed A"].ID, Title: pointer("Feed X"), IsStarred: pointer(true)},
	}
	feeds, err := st.EditFeeds(context.Background(), ops)
	r.NoError(err)

	a.Len(feeds, 1)

	a.False(existf("Feed A", false))
	a.True(existf("Feed X", true))
}
