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
	a.Nil(stats.LastPullTime)
	a.Nil(stats.MostRecentUpdateTime)
}

func TestGetGlobalStatsExtendedOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	st := newTestStore(t)

	dbFeeds := []*feedRecord{
		{
			title:      "Feed A",
			feedURL:    "http://a.com/feed.xml",
			subscribed: mustTime(t, "2022-07-18T22:04:37Z"),
			lastPulled: mustTime(t, "2022-07-18T22:04:37Z"),
			updated:    toNullString("2022-04-20T16:32:30.760+02:00"),
			entries: []*entryRecord{
				{
					title:   "Entry A1",
					extID:   "A1",
					isRead:  true,
					updated: toNullString("2022-07-16T23:39:07.383+02:00"),
					url:     toNullString("http://a.com/a1.html"),
				},
				{
					title:   "Entry A2",
					extID:   "A2",
					isRead:  false,
					updated: toNullString("2022-07-16T23:42:24.988+02:00"),
					url:     toNullString("http://a.com/a2.html"),
				},
				{
					title:   "Entry A3",
					extID:   "A3",
					isRead:  true,
					updated: toNullString("2022-03-18T22:51:49.404+02:00"),
					url:     toNullString("http://a.com/a3.html"),
				},
			},
		},
		{
			title:      "Feed X",
			feedURL:    "http://x.com/feed.xml",
			subscribed: mustTime(t, "2022-07-18T22:04:45Z"),
			lastPulled: mustTime(t, "2022-07-18T22:04:45Z"),
			updated:    toNullString("2022-03-19T16:23:18.600+02:00"),
			entries: []*entryRecord{
				{
					title:   "Entry X1",
					extID:   "X1",
					isRead:  true,
					updated: toNullString("2022-07-16T23:43:12.759+02:00"),
					url:     toNullString("http://x.com/x1.html"),
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
	a.Equal(
		"2022-07-18T22:04:45Z",
		stats.LastPullTime.UTC().Format(time.RFC3339),
	)
	a.Equal(
		"2022-04-20T14:32:30Z",
		stats.MostRecentUpdateTime.UTC().Format(time.RFC3339),
	)
}
