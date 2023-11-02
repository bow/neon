// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"sort"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullFeedsOkEmptyDB(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	st.parser.EXPECT().
		ParseURLWithContext(gomock.Any(), gomock.Any()).
		MaxTimes(0)

	c := st.PullFeeds(context.Background())
	a.Empty(c)
}

func TestPullFeedsOkEmptyEntries(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: WrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: WrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{},
		},
	}

	st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[0].FeedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, dbFeeds[0]), nil)

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[1].FeedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, dbFeeds[1]), nil)

	c := st.PullFeeds(context.Background())

	got := make([]PullResult, 0)
	for res := range c {
		got = append(got, res)
	}

	want := []PullResult{
		{
			url:    &dbFeeds[0].FeedURL,
			status: pullSuccess,
			feed:   nil,
			err:    nil,
		},
		{
			url:    &dbFeeds[1].FeedURL,
			status: pullSuccess,
			feed:   nil,
			err:    nil,
		},
	}

	a.ElementsMatch(want, got)
}

func TestPullFeedsOkNoNewEntries(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:   "Feed A",
			FeedURL: "http://a.com/feed.xml",
			Updated: WrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{
					Title:   "Entry A1",
					ExtID:   "A1",
					IsRead:  true,
					Updated: WrapNullString("2022-07-16T23:39:07.383+02:00"),
					URL:     WrapNullString("http://a.com/a1.html"),
				},
				{
					Title:   "Entry A2",
					ExtID:   "A2",
					IsRead:  true,
					Updated: WrapNullString("2022-07-16T23:42:24.988+02:00"),
					URL:     WrapNullString("http://a.com/a2.html"),
				},
			},
		},
		{
			Title:   "Feed X",
			FeedURL: "http://x.com/feed.xml",
			Updated: WrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{
					Title:   "Entry X1",
					ExtID:   "X1",
					IsRead:  true,
					Updated: WrapNullString("2022-07-16T23:43:12.759+02:00"),
					URL:     WrapNullString("http://x.com/x1.html"),
				},
			},
		},
	}

	st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	pulledFeeds := []*Feed{
		{
			Title:   dbFeeds[0].Title,
			FeedURL: dbFeeds[0].FeedURL,
			Updated: dbFeeds[0].Updated,
			Entries: []*Entry{
				{
					Title:   dbFeeds[0].Entries[0].Title,
					ExtID:   dbFeeds[0].Entries[0].ExtID,
					Updated: dbFeeds[0].Entries[0].Updated,
					URL:     dbFeeds[0].Entries[0].URL,
				},
				{
					Title:   dbFeeds[0].Entries[1].Title,
					ExtID:   dbFeeds[0].Entries[1].ExtID,
					Updated: dbFeeds[0].Entries[1].Updated,
					URL:     dbFeeds[0].Entries[1].URL,
				},
			},
		},
		{
			Title:   dbFeeds[1].Title,
			FeedURL: dbFeeds[1].FeedURL,
			Updated: dbFeeds[1].Updated,
			Entries: []*Entry{
				{
					Title:   dbFeeds[1].Entries[0].Title,
					ExtID:   dbFeeds[1].Entries[0].ExtID,
					Updated: dbFeeds[1].Entries[0].Updated,
					URL:     dbFeeds[1].Entries[0].URL,
				},
			},
		},
	}

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[0].FeedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeeds[0]), nil)

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[1].FeedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeeds[1]), nil)

	c := st.PullFeeds(context.Background())

	got := make([]PullResult, 0)
	for res := range c {
		got = append(got, res)
	}

	want := []PullResult{
		{
			url:    &pulledFeeds[0].FeedURL,
			status: pullSuccess,
			err:    nil,
			feed:   nil,
		},
		{
			url:    &pulledFeeds[1].FeedURL,
			status: pullSuccess,
			err:    nil,
			feed:   nil,
		},
	}

	a.ElementsMatch(want, got)
}

func TestPullFeedsOkSomeNewEntries(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:      "Feed A",
			FeedURL:    "http://a.com/feed.xml",
			Subscribed: "2022-07-18T22:04:37Z",
			LastPulled: "2022-07-18T22:04:37Z",
			Updated:    WrapNullString("2022-03-19T16:23:18.600+02:00"),
			Entries: []*Entry{
				{
					// This entry should not be returned later; 'updated' remains the same.
					Title:   "Entry A1",
					ExtID:   "A1",
					IsRead:  true,
					Updated: WrapNullString("2022-07-16T23:39:07.383+02:00"),
					URL:     WrapNullString("http://a.com/a1.html"),
				},
				{
					// This entry should not be returned later; 'updated' remains the same.
					Title:   "Entry A2",
					ExtID:   "A2",
					IsRead:  false,
					Updated: WrapNullString("2022-07-16T23:42:24.988+02:00"),
					URL:     WrapNullString("http://a.com/a2.html"),
				},
				{
					// This entry should be returned later; 'updated' will be changed.
					Title:   "Entry A3",
					ExtID:   "A3",
					IsRead:  true,
					Updated: WrapNullString("2022-03-18T22:51:49.404+02:00"),
					URL:     WrapNullString("http://a.com/a3.html"),
				},
			},
		},
		{
			Title:      "Feed X",
			FeedURL:    "http://x.com/feed.xml",
			Subscribed: "2022-07-18T22:04:45Z",
			LastPulled: "2022-07-18T22:04:45Z",
			Updated:    WrapNullString("2022-04-20T16:32:30.760+02:00"),
			Entries: []*Entry{
				{
					// This entry should not be returned later; 'updated' remains the same.
					Title:   "Entry X1",
					ExtID:   "X1",
					IsRead:  true,
					Updated: WrapNullString("2022-07-16T23:43:12.759+02:00"),
					URL:     WrapNullString("http://x.com/x1.html"),
				},
			},
		},
	}

	keys := st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	pulledFeeds := []*Feed{
		{
			Title:   dbFeeds[0].Title,
			FeedURL: dbFeeds[0].FeedURL,
			Updated: WrapNullString("2022-07-18T22:51:49.404+02:00"),
			Entries: []*Entry{
				{
					Title:   dbFeeds[0].Entries[0].Title,
					ExtID:   dbFeeds[0].Entries[0].ExtID,
					Updated: dbFeeds[0].Entries[0].Updated,
					URL:     dbFeeds[0].Entries[0].URL,
				},
				{
					Title:   dbFeeds[0].Entries[1].Title,
					ExtID:   dbFeeds[0].Entries[1].ExtID,
					Updated: dbFeeds[0].Entries[1].Updated,
					URL:     dbFeeds[0].Entries[1].URL,
				},
				{
					Title:   dbFeeds[0].Entries[2].Title,
					ExtID:   dbFeeds[0].Entries[2].ExtID,
					Updated: WrapNullString("2022-07-19T16:23:18.600+02:00"),
					URL:     dbFeeds[0].Entries[2].URL,
				},
			},
		},
		{
			Title:   dbFeeds[1].Title,
			FeedURL: dbFeeds[1].FeedURL,
			Updated: WrapNullString("2022-07-18T22:21:41.647+02:00"),
			Entries: []*Entry{
				{
					Title:   dbFeeds[1].Entries[0].Title,
					ExtID:   dbFeeds[1].Entries[0].ExtID,
					Updated: dbFeeds[1].Entries[0].Updated,
					URL:     dbFeeds[1].Entries[0].URL,
				},
				{
					Title:   "Entry X2",
					ExtID:   "X2",
					Updated: WrapNullString("2022-07-18T22:21:41.647+02:00"),
					URL:     WrapNullString("http://x.com/x2.html"),
				},
			},
		},
	}

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[0].FeedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeeds[0]), nil)

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[1].FeedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeeds[1]), nil)

	c := st.PullFeeds(context.Background())

	got := make([]PullResult, 0)
	for res := range c {
		got = append(got, res)
	}

	feedURL0 := pulledFeeds[0].FeedURL
	feedURL1 := pulledFeeds[1].FeedURL

	want := []PullResult{
		{
			url:    &dbFeeds[0].FeedURL,
			status: pullSuccess,
			err:    nil,
			feed: &Feed{
				DBID:       keys[pulledFeeds[0].Title].DBID,
				Title:      pulledFeeds[0].Title,
				FeedURL:    pulledFeeds[0].FeedURL,
				Updated:    st.getFeedUpdateTime(feedURL0),
				Subscribed: st.getFeedSubTime(feedURL0),
				LastPulled: "",
				Entries: []*Entry{
					{
						DBID:      st.getEntryDBID(feedURL0, pulledFeeds[0].Entries[1].ExtID),
						FeedDBID:  keys[pulledFeeds[0].Title].DBID,
						Title:     pulledFeeds[0].Entries[1].Title,
						ExtID:     pulledFeeds[0].Entries[1].ExtID,
						Updated:   st.getEntryUpdateTime(feedURL0, pulledFeeds[0].Entries[1].ExtID),
						Published: st.getEntryPubTime(feedURL0, pulledFeeds[0].Entries[1].ExtID),
						URL:       pulledFeeds[0].Entries[1].URL,
						IsRead:    false,
					},
					{
						DBID:      st.getEntryDBID(feedURL0, pulledFeeds[0].Entries[2].ExtID),
						FeedDBID:  keys[pulledFeeds[0].Title].DBID,
						Title:     pulledFeeds[0].Entries[2].Title,
						ExtID:     pulledFeeds[0].Entries[2].ExtID,
						Updated:   st.getEntryUpdateTime(feedURL0, pulledFeeds[0].Entries[2].ExtID),
						Published: st.getEntryPubTime(feedURL0, pulledFeeds[0].Entries[2].ExtID),
						URL:       pulledFeeds[0].Entries[2].URL,
						IsRead:    false,
					},
				},
			},
		},
		{
			url:    &dbFeeds[1].FeedURL,
			status: pullSuccess,
			err:    nil,
			feed: &Feed{
				DBID:       keys[pulledFeeds[1].Title].DBID,
				Title:      pulledFeeds[1].Title,
				FeedURL:    pulledFeeds[1].FeedURL,
				Updated:    st.getFeedUpdateTime(feedURL1),
				Subscribed: st.getFeedSubTime(feedURL1),
				LastPulled: "",
				Entries: []*Entry{
					{
						DBID:      st.getEntryDBID(feedURL1, pulledFeeds[1].Entries[1].ExtID),
						FeedDBID:  keys[pulledFeeds[1].Title].DBID,
						Title:     pulledFeeds[1].Entries[1].Title,
						ExtID:     pulledFeeds[1].Entries[1].ExtID,
						Updated:   st.getEntryUpdateTime(feedURL1, pulledFeeds[1].Entries[1].ExtID),
						Published: st.getEntryPubTime(feedURL1, pulledFeeds[1].Entries[1].ExtID),
						URL:       pulledFeeds[1].Entries[1].URL,
						IsRead:    false,
					},
				},
			},
		},
	}

	// Sort inner entries first, since ElementsMatch cares about inner array elements order.
	sortPullResultEntries(want)
	sortPullResultEntries(got)

	// Set LastPulled fields to empty strings as this value is always updated on every pull.
	for _, item := range got {
		item.feed.LastPulled = ""
	}

	a.ElementsMatch(want, got)
}

func sortPullResultEntries(arr []PullResult) {
	for _, item := range arr {
		sort.SliceStable(
			item.feed.Entries,
			func(i, j int) bool {
				return item.feed.Entries[i].ExtID < item.feed.Entries[j].ExtID
			},
		)
	}
}

func toGFeed(t *testing.T, feed *Feed) *gofeed.Feed {
	t.Helper()
	gfeed := gofeed.Feed{
		Title:    feed.Title,
		FeedLink: feed.FeedURL,
	}
	if feed.Updated.String != "" {
		gfeed.Updated = feed.Updated.String
		gfeed.UpdatedParsed = ts(t, feed.Updated.String)
	}
	for _, entry := range feed.Entries {
		item := gofeed.Item{
			GUID:    entry.ExtID,
			Link:    entry.URL.String,
			Title:   entry.Title,
			Content: entry.Content.String,
		}
		if entry.Published.String != "" {
			item.Published = entry.Published.String
			item.PublishedParsed = ts(t, entry.Published.String)
		}
		if entry.Updated.String != "" {
			item.Updated = entry.Updated.String
			item.UpdatedParsed = ts(t, entry.Updated.String)
		}
		gfeed.Items = append(gfeed.Items, &item)
	}
	return &gfeed
}
