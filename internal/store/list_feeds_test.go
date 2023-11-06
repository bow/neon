// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

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
	st := newTestStore(t)

	feeds, err := st.ListFeeds(context.Background())
	r.NoError(err)

	a.Empty(feeds)
}

func TestListFeedsOkExtended(t *testing.T) {
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
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: toNullString("2022-04-20T16:32:30.760+02:00"),
		},
	}
	st.addFeeds(dbFeeds)

	r.Equal(2, st.countFeeds())

	feeds, err := st.ListFeeds(context.Background())
	r.NoError(err)
	r.NotEmpty(feeds)

	a.Len(feeds, 2)

	feed0 := feeds[0]
	a.Equal(feed0.FeedURL, dbFeeds[1].FeedURL)

	feed1 := feeds[1]
	a.Equal(feed1.FeedURL, dbFeeds[0].FeedURL)
}
