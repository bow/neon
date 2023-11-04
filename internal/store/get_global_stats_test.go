// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGlobalStatsEmptyOk(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	r := require.New(t)
	st := newTestStore(t)

	stats, err := st.GetGlobalStats(context.Background())
	r.NoError(err)
	r.NotNil(stats)

	a.Equal(uint32(0), stats.NumFeeds)
	a.Equal(uint32(0), stats.NumEntries)
	a.Equal(uint32(0), stats.NumEntriesUnread)

	lpt, err := stats.LastPullTime()
	r.NoError(err)
	a.Nil(lpt)

	mrpt, err := stats.MostRecentUpdateTime()
	r.NoError(err)
	a.Nil(mrpt)
}

func TestGetGlobalStatsExtendedOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	st := newTestStore(t)

	dbFeeds := []*Feed{
		{
			Title:      "Feed A",
			FeedURL:    "http://a.com/feed.xml",
			Subscribed: "2022-07-18T22:04:37Z",
			LastPulled: "2022-07-18T22:04:37Z",
			Updated:    WrapNullString("2022-04-20T16:32:30.760+02:00"),
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
					IsRead:  false,
					Updated: WrapNullString("2022-07-16T23:42:24.988+02:00"),
					URL:     WrapNullString("http://a.com/a2.html"),
				},
				{
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
			Updated:    WrapNullString("2022-03-19T16:23:18.600+02:00"),
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
	_ = st.addFeeds(dbFeeds)
	r.Equal(2, st.countFeeds())

	stats, err := st.GetGlobalStats(context.Background())
	r.NoError(err)
	r.NotNil(stats)

	a.Equal(uint32(2), stats.NumFeeds)
	a.Equal(uint32(4), stats.NumEntries)
	a.Equal(uint32(1), stats.NumEntriesUnread)

	lpt, err := stats.LastPullTime()
	r.NoError(err)
	a.Equal(
		"2022-07-18T22:04:45Z",
		lpt.UTC().Format(time.RFC3339),
	)

	mrpt, err := stats.MostRecentUpdateTime()
	r.NoError(err)
	a.Equal(
		"2022-04-20T14:32:30Z",
		mrpt.UTC().Format(time.RFC3339),
	)
}
