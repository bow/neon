// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportSubscriptionOkEmpty(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	sub, err := st.ExportSubscription(context.Background(), pointer("export"))
	r.NoError(err)
	r.NotNil(sub)

	a.Equal(pointer("export"), sub.Title)
	a.Empty(sub.Feeds)
}

func TestExportSubscriptionOkExtended(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*feedRecord{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-03-19T16:23:18.600+02:00")),
			entries: []*entryRecord{
				{title: "Entry A1", isRead: false},
				{title: "Entry A2", isRead: false},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullTime(mustTime(t, "2022-04-20T16:32:30.760+02:00")),
			entries: []*entryRecord{
				{title: "Entry X1", isRead: false},
			},
			tags: []string{"foo", "baz"},
		},
		{
			title:     "Feed Q",
			feedURL:   "http://q.com/feed.xml",
			updated:   toNullTime(mustTime(t, "2022-05-02T11:47:33.683+02:00")),
			isStarred: true,
			entries: []*entryRecord{
				{title: "Entry Q1", isRead: false},
			},
		},
	}
	st.addFeeds(dbFeeds)
	r.Equal(3, st.countFeeds())

	sub, err := st.ExportSubscription(context.Background(), pointer("Test Export"))
	r.NoError(err)

	a.NotNil(sub.Title)
	a.Equal(pointer("Test Export"), sub.Title)
	a.Len(sub.Feeds, 3)

	a.Equal(sub.Feeds[0].Title, dbFeeds[2].title)
	a.Equal(sub.Feeds[0].FeedURL, dbFeeds[2].feedURL)
	a.Equal(sub.Feeds[0].IsStarred, dbFeeds[2].isStarred)
	a.ElementsMatch(sub.Feeds[0].Tags, dbFeeds[2].tags)

	a.Equal(sub.Feeds[1].Title, dbFeeds[1].title)
	a.Equal(sub.Feeds[1].FeedURL, dbFeeds[1].feedURL)
	a.Equal(sub.Feeds[1].IsStarred, dbFeeds[1].isStarred)
	a.ElementsMatch(sub.Feeds[1].Tags, dbFeeds[1].tags)

	a.Equal(sub.Feeds[2].Title, dbFeeds[0].title)
	a.Equal(sub.Feeds[2].FeedURL, dbFeeds[0].feedURL)
	a.Equal(sub.Feeds[2].IsStarred, dbFeeds[0].isStarred)
	a.ElementsMatch(sub.Feeds[2].Tags, dbFeeds[0].tags)
}
