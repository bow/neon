// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/bow/iris/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportSubscriptionOkNoFeeds(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	sub := internal.Subscription{}

	nproc, nimp, err := st.ImportSubscription(context.Background(), &sub)
	r.NoError(err)

	a.Equal(0, nproc)
	a.Equal(0, nimp)
	a.Equal(0, st.countFeeds())
}

func TestImportSubscriptionOkMinimal(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	existf := func() bool {
		return st.rowExists(
			feedExistSQL,
			"Feed A",
			nil,
			"http://a.com/feed.xml",
			nil,
			false,
		)
	}

	r.Equal(0, st.countFeeds())
	a.False(existf())

	sub := internal.Subscription{
		Feeds: []*internal.Feed{
			{Title: "Feed A", FeedURL: "http://a.com/feed.xml"},
		},
	}

	nproc, nimp, err := st.ImportSubscription(context.Background(), &sub)
	r.NoError(err)

	a.Equal(1, nproc)
	a.Equal(1, nimp)
	a.Equal(1, st.countFeeds())
	a.True(existf())
}

func TestImportSubscriptionOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*feedRecord{
		{
			title:     "Feed BC",
			feedURL:   "http://bc.com/feed.xml",
			updated:   toNullTime(mustTime(t, "2022-03-19T16:23:18.600+02:00")),
			isStarred: false,
			entries: []*entryRecord{
				{title: "Entry BC1", isRead: false},
				{title: "Entry BC2", isRead: true},
			},
		},
		{
			title:     "Feed D",
			feedURL:   "http://d.com/feed.xml",
			updated:   toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
			isStarred: true,
			entries: []*entryRecord{
				{title: "Entry D1", isRead: false},
			},
			tags: []string{"foo", "baz"},
		},
	}
	st.addFeeds(dbFeeds)

	existfA := func() bool {
		return st.rowExists(
			feedExistSQL,
			"Feed A",
			"New feed",
			"http://a.com/feed.xml",
			"http://a.com",
			false,
		)
	}
	existfBC := func() bool {
		return st.rowExists(
			feedExistSQL,
			"Feed BC",
			"Updated feed",
			"http://bc.com/feed.xml",
			nil,
			true,
		)
	}

	sub := internal.Subscription{
		Feeds: []*internal.Feed{
			{
				Title:       "Feed A",
				FeedURL:     "http://a.com/feed.xml",
				SiteURL:     pointer("http://a.com"),
				Description: pointer("New feed"),
			},
			{
				Title:       "Feed BC",
				FeedURL:     "http://bc.com/feed.xml",
				Description: pointer("Updated feed"),
				IsStarred:   true,
			},
		},
	}

	r.Equal(2, st.countFeeds())
	a.False(existfA())
	a.False(existfBC())

	nproc, nimp, err := st.ImportSubscription(context.Background(), &sub)
	r.NoError(err)

	a.Equal(2, nproc)
	a.Equal(1, nimp)
	a.Equal(3, st.countFeeds())
	a.True(existfA())
	a.True(existfBC())
}
