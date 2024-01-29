// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	now          = time.Now()
	oneHourAgo   = now.Add(-1 * time.Hour)
	yesterday    = now.Add(-1 * 24 * time.Hour)
	threeDaysAgo = now.Add(-3 * 24 * time.Hour)
	lastWeek     = now.Add(-7 * 24 * time.Hour)
	twoWeeksAgo  = now.Add(-14 * 24 * time.Hour)
)

func TestFeedSortEntries(t *testing.T) {
	a := assert.New(t)
	f := Feed{
		Entries: []*Entry{
			{
				Title:     "A",
				IsRead:    true,
				Published: &twoWeeksAgo,
				Updated:   nil,
			},
			{
				Title:     "B",
				IsRead:    false,
				Published: &lastWeek,
				Updated:   &yesterday,
			},
			{
				Title:     "C",
				IsRead:    false,
				Published: &lastWeek,
				Updated:   &yesterday,
			},
			{
				Title:     "D",
				IsRead:    false,
				Published: nil,
				Updated:   &threeDaysAgo,
			},
			{
				Title:     "E",
				IsRead:    false,
				Published: &oneHourAgo,
				Updated:   nil,
			},
			{
				Title:     "F",
				IsRead:    true,
				Published: nil,
				Updated:   nil,
			},
		},
	}
	f.SortEntries()

	want := []string{"E", "B", "C", "D", "A", "F"}
	got := make([]string, len(f.Entries))

	for i, entry := range f.Entries {
		title := entry.Title
		got[i] = title
	}

	a.Len(got, 6)
	a.Equal(want, got)
}
