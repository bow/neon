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

func TestPullFeedsAllOkEmptyDB(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	r.Equal(0, st.countFeeds())

	st.parser.EXPECT().
		ParseURLWithContext(gomock.Any(), gomock.Any()).
		MaxTimes(0)

	c := st.PullFeeds(context.Background(), nil)
	a.Empty(c)
}

func TestPullFeedsAllOkEmptyEntries(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullString("2022-03-19T16:23:18.600+02:00"),
			entries: []*Entry{},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullString("2022-04-20T16:32:30.760+02:00"),
			entries: []*Entry{},
		},
	}

	st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[0].feedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, dbFeeds[0]), nil)

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[1].feedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, dbFeeds[1]), nil)

	c := st.PullFeeds(context.Background(), nil)

	got := make([]PullResult, 0)
	for res := range c {
		got = append(got, res)
	}

	want := []PullResult{
		{
			url:    &dbFeeds[0].feedURL,
			status: pullSuccess,
			feed:   nil,
			err:    nil,
		},
		{
			url:    &dbFeeds[1].feedURL,
			status: pullSuccess,
			feed:   nil,
			err:    nil,
		},
	}

	a.ElementsMatch(want, got)
}

func TestPullFeedsAllOkNoNewEntries(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			title:   "Feed A",
			feedURL: "http://a.com/feed.xml",
			updated: toNullString("2022-03-19T16:23:18.600+02:00"),
			entries: []*Entry{
				{
					Title:   "Entry A1",
					ExtID:   "A1",
					IsRead:  true,
					Updated: toNullString("2022-07-16T23:39:07.383+02:00"),
					URL:     toNullString("http://a.com/a1.html"),
				},
				{
					Title:   "Entry A2",
					ExtID:   "A2",
					IsRead:  true,
					Updated: toNullString("2022-07-16T23:42:24.988+02:00"),
					URL:     toNullString("http://a.com/a2.html"),
				},
			},
		},
		{
			title:   "Feed X",
			feedURL: "http://x.com/feed.xml",
			updated: toNullString("2022-04-20T16:32:30.760+02:00"),
			entries: []*Entry{
				{
					Title:   "Entry X1",
					ExtID:   "X1",
					IsRead:  true,
					Updated: toNullString("2022-07-16T23:43:12.759+02:00"),
					URL:     toNullString("http://x.com/x1.html"),
				},
			},
		},
	}

	st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	pulledFeeds := []*Feed{
		{
			title:   dbFeeds[0].title,
			feedURL: dbFeeds[0].feedURL,
			updated: dbFeeds[0].updated,
			entries: []*Entry{
				{
					Title:   dbFeeds[0].entries[0].Title,
					ExtID:   dbFeeds[0].entries[0].ExtID,
					Updated: dbFeeds[0].entries[0].Updated,
					URL:     dbFeeds[0].entries[0].URL,
				},
				{
					Title:   dbFeeds[0].entries[1].Title,
					ExtID:   dbFeeds[0].entries[1].ExtID,
					Updated: dbFeeds[0].entries[1].Updated,
					URL:     dbFeeds[0].entries[1].URL,
				},
			},
		},
		{
			title:   dbFeeds[1].title,
			feedURL: dbFeeds[1].feedURL,
			updated: dbFeeds[1].updated,
			entries: []*Entry{
				{
					Title:   dbFeeds[1].entries[0].Title,
					ExtID:   dbFeeds[1].entries[0].ExtID,
					Updated: dbFeeds[1].entries[0].Updated,
					URL:     dbFeeds[1].entries[0].URL,
				},
			},
		},
	}

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[0].feedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeeds[0]), nil)

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[1].feedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeeds[1]), nil)

	c := st.PullFeeds(context.Background(), nil)

	got := make([]PullResult, 0)
	for res := range c {
		got = append(got, res)
	}

	want := []PullResult{
		{
			url:    &pulledFeeds[0].feedURL,
			status: pullSuccess,
			err:    nil,
			feed:   nil,
		},
		{
			url:    &pulledFeeds[1].feedURL,
			status: pullSuccess,
			err:    nil,
			feed:   nil,
		},
	}

	a.ElementsMatch(want, got)
}

func TestPullFeedsAllOkSomeNewEntries(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			title:      "Feed A",
			feedURL:    "http://a.com/feed.xml",
			subscribed: "2022-07-18T22:04:37Z",
			lastPulled: "2022-07-18T22:04:37Z",
			updated:    toNullString("2022-03-19T16:23:18.600+02:00"),
			entries: []*Entry{
				{
					// This entry should not be returned later; 'updated' remains the same.
					Title:   "Entry A1",
					ExtID:   "A1",
					IsRead:  true,
					Updated: toNullString("2022-07-16T23:39:07.383+02:00"),
					URL:     toNullString("http://a.com/a1.html"),
				},
				{
					// This entry should not be returned later; 'updated' remains the same.
					Title:   "Entry A2",
					ExtID:   "A2",
					IsRead:  false,
					Updated: toNullString("2022-07-16T23:42:24.988+02:00"),
					URL:     toNullString("http://a.com/a2.html"),
				},
				{
					// This entry should be returned later; 'updated' will be changed.
					Title:   "Entry A3",
					ExtID:   "A3",
					IsRead:  true,
					Updated: toNullString("2022-03-18T22:51:49.404+02:00"),
					URL:     toNullString("http://a.com/a3.html"),
				},
			},
		},
		{
			title:      "Feed X",
			feedURL:    "http://x.com/feed.xml",
			subscribed: "2022-07-18T22:04:45Z",
			lastPulled: "2022-07-18T22:04:45Z",
			updated:    toNullString("2022-04-20T16:32:30.760+02:00"),
			entries: []*Entry{
				{
					// This entry should not be returned later; 'updated' remains the same.
					Title:   "Entry X1",
					ExtID:   "X1",
					IsRead:  true,
					Updated: toNullString("2022-07-16T23:43:12.759+02:00"),
					URL:     toNullString("http://x.com/x1.html"),
				},
			},
		},
	}

	keys := st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	pulledFeeds := []*Feed{
		{
			title:   dbFeeds[0].title,
			feedURL: dbFeeds[0].feedURL,
			updated: toNullString("2022-07-18T22:51:49.404+02:00"),
			entries: []*Entry{
				{
					Title:   dbFeeds[0].entries[0].Title,
					ExtID:   dbFeeds[0].entries[0].ExtID,
					Updated: dbFeeds[0].entries[0].Updated,
					URL:     dbFeeds[0].entries[0].URL,
				},
				{
					Title:   dbFeeds[0].entries[1].Title,
					ExtID:   dbFeeds[0].entries[1].ExtID,
					Updated: dbFeeds[0].entries[1].Updated,
					URL:     dbFeeds[0].entries[1].URL,
				},
				{
					Title:   dbFeeds[0].entries[2].Title,
					ExtID:   dbFeeds[0].entries[2].ExtID,
					Updated: toNullString("2022-07-19T16:23:18.600+02:00"),
					URL:     dbFeeds[0].entries[2].URL,
				},
			},
		},
		{
			title:   dbFeeds[1].title,
			feedURL: dbFeeds[1].feedURL,
			updated: toNullString("2022-07-18T22:21:41.647+02:00"),
			entries: []*Entry{
				{
					Title:   dbFeeds[1].entries[0].Title,
					ExtID:   dbFeeds[1].entries[0].ExtID,
					Updated: dbFeeds[1].entries[0].Updated,
					URL:     dbFeeds[1].entries[0].URL,
				},
				{
					Title:   "Entry X2",
					ExtID:   "X2",
					Updated: toNullString("2022-07-18T22:21:41.647+02:00"),
					URL:     toNullString("http://x.com/x2.html"),
				},
			},
		},
	}

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[0].feedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeeds[0]), nil)

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[1].feedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeeds[1]), nil)

	c := st.PullFeeds(context.Background(), nil)

	got := make([]PullResult, 0)
	for res := range c {
		got = append(got, res)
	}

	feedURL0 := pulledFeeds[0].feedURL
	feedURL1 := pulledFeeds[1].feedURL

	want := []PullResult{
		{
			url:    &dbFeeds[0].feedURL,
			status: pullSuccess,
			err:    nil,
			feed: &Feed{
				id:         keys[pulledFeeds[0].title].ID,
				title:      pulledFeeds[0].title,
				feedURL:    pulledFeeds[0].feedURL,
				updated:    st.getFeedUpdateTime(feedURL0),
				subscribed: st.getFeedSubTime(feedURL0),
				lastPulled: "",
				entries: []*Entry{
					{
						ID:        st.getEntryID(feedURL0, pulledFeeds[0].entries[1].ExtID),
						FeedID:    keys[pulledFeeds[0].title].ID,
						Title:     pulledFeeds[0].entries[1].Title,
						ExtID:     pulledFeeds[0].entries[1].ExtID,
						Updated:   st.getEntryUpdateTime(feedURL0, pulledFeeds[0].entries[1].ExtID),
						Published: st.getEntryPubTime(feedURL0, pulledFeeds[0].entries[1].ExtID),
						URL:       pulledFeeds[0].entries[1].URL,
						IsRead:    false,
					},
					{
						ID:        st.getEntryID(feedURL0, pulledFeeds[0].entries[2].ExtID),
						FeedID:    keys[pulledFeeds[0].title].ID,
						Title:     pulledFeeds[0].entries[2].Title,
						ExtID:     pulledFeeds[0].entries[2].ExtID,
						Updated:   st.getEntryUpdateTime(feedURL0, pulledFeeds[0].entries[2].ExtID),
						Published: st.getEntryPubTime(feedURL0, pulledFeeds[0].entries[2].ExtID),
						URL:       pulledFeeds[0].entries[2].URL,
						IsRead:    false,
					},
				},
			},
		},
		{
			url:    &dbFeeds[1].feedURL,
			status: pullSuccess,
			err:    nil,
			feed: &Feed{
				id:         keys[pulledFeeds[1].title].ID,
				title:      pulledFeeds[1].title,
				feedURL:    pulledFeeds[1].feedURL,
				updated:    st.getFeedUpdateTime(feedURL1),
				subscribed: st.getFeedSubTime(feedURL1),
				lastPulled: "",
				entries: []*Entry{
					{
						ID:        st.getEntryID(feedURL1, pulledFeeds[1].entries[1].ExtID),
						FeedID:    keys[pulledFeeds[1].title].ID,
						Title:     pulledFeeds[1].entries[1].Title,
						ExtID:     pulledFeeds[1].entries[1].ExtID,
						Updated:   st.getEntryUpdateTime(feedURL1, pulledFeeds[1].entries[1].ExtID),
						Published: st.getEntryPubTime(feedURL1, pulledFeeds[1].entries[1].ExtID),
						URL:       pulledFeeds[1].entries[1].URL,
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
		item.feed.lastPulled = ""
	}

	a.ElementsMatch(want, got)
}

func TestPullFeedsSelectedOkSomeNewEntries(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		// This feed should not be returned later, it is not selected.
		{
			title:      "Feed A",
			feedURL:    "http://a.com/feed.xml",
			subscribed: "2022-07-18T22:04:37Z",
			lastPulled: "2022-07-18T22:04:37Z",
			updated:    toNullString("2022-03-19T16:23:18.600+02:00"),
			entries: []*Entry{
				{
					Title:   "Entry A1",
					ExtID:   "A1",
					IsRead:  true,
					Updated: toNullString("2022-07-16T23:39:07.383+02:00"),
					URL:     toNullString("http://a.com/a1.html"),
				},
				{
					Title:   "Entry A2",
					ExtID:   "A2",
					IsRead:  false,
					Updated: toNullString("2022-07-16T23:42:24.988+02:00"),
					URL:     toNullString("http://a.com/a2.html"),
				},
				{
					Title:   "Entry A3",
					ExtID:   "A3",
					IsRead:  true,
					Updated: toNullString("2022-03-18T22:51:49.404+02:00"),
					URL:     toNullString("http://a.com/a3.html"),
				},
			},
		},
		// This feed should be returned later, it is selected.
		{
			title:      "Feed X",
			feedURL:    "http://x.com/feed.xml",
			subscribed: "2022-07-18T22:04:45Z",
			lastPulled: "2022-07-18T22:04:45Z",
			updated:    toNullString("2022-04-20T16:32:30.760+02:00"),
			entries: []*Entry{
				{
					// This entry should not be returned later; 'updated' remains the same.
					Title:   "Entry X1",
					ExtID:   "X1",
					IsRead:  true,
					Updated: toNullString("2022-07-16T23:43:12.759+02:00"),
					URL:     toNullString("http://x.com/x1.html"),
				},
			},
		},
	}

	keys := st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	pulledFeed := &Feed{
		title:   dbFeeds[1].title,
		feedURL: dbFeeds[1].feedURL,
		updated: toNullString("2022-07-18T22:21:41.647+02:00"),
		entries: []*Entry{
			{
				Title:   dbFeeds[1].entries[0].Title,
				ExtID:   dbFeeds[1].entries[0].ExtID,
				Updated: dbFeeds[1].entries[0].Updated,
				URL:     dbFeeds[1].entries[0].URL,
			},
			{
				Title:   "Entry X2",
				ExtID:   "X2",
				Updated: toNullString("2022-07-18T22:21:41.647+02:00"),
				URL:     toNullString("http://x.com/x2.html"),
			},
		},
	}

	st.parser.EXPECT().
		ParseURLWithContext(dbFeeds[1].feedURL, gomock.Any()).
		MaxTimes(1).
		Return(toGFeed(t, pulledFeed), nil)

	c := st.PullFeeds(context.Background(), []ID{keys[pulledFeed.title].ID})

	got := make([]PullResult, 0)
	for res := range c {
		got = append(got, res)
	}

	want := []PullResult{
		{
			url:    &dbFeeds[1].feedURL,
			status: pullSuccess,
			err:    nil,
			feed: &Feed{
				id:         keys[pulledFeed.title].ID,
				title:      pulledFeed.title,
				feedURL:    pulledFeed.feedURL,
				updated:    st.getFeedUpdateTime(pulledFeed.feedURL),
				subscribed: st.getFeedSubTime(pulledFeed.feedURL),
				lastPulled: "",
				entries: []*Entry{
					{
						ID:        st.getEntryID(pulledFeed.feedURL, pulledFeed.entries[1].ExtID),
						FeedID:    keys[pulledFeed.title].ID,
						Title:     pulledFeed.entries[1].Title,
						ExtID:     pulledFeed.entries[1].ExtID,
						Updated:   st.getEntryUpdateTime(pulledFeed.feedURL, pulledFeed.entries[1].ExtID),
						Published: st.getEntryPubTime(pulledFeed.feedURL, pulledFeed.entries[1].ExtID),
						URL:       pulledFeed.entries[1].URL,
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
		item.feed.lastPulled = ""
	}

	a.ElementsMatch(want, got)
}

func sortPullResultEntries(arr []PullResult) {
	for _, item := range arr {
		sort.SliceStable(
			item.feed.entries,
			func(i, j int) bool {
				return item.feed.entries[i].ExtID < item.feed.entries[j].ExtID
			},
		)
	}
}

func toGFeed(t *testing.T, feed *Feed) *gofeed.Feed {
	t.Helper()
	gfeed := gofeed.Feed{
		Title:    feed.title,
		FeedLink: feed.feedURL,
	}
	if feed.updated.String != "" {
		gfeed.Updated = feed.updated.String
		gfeed.UpdatedParsed = ts(t, feed.updated.String)
	}
	for _, entry := range feed.entries {
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
