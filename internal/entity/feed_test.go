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

func TestFeedEntries(t *testing.T) {
	a := assert.New(t)
	f := Feed{
		Entries: map[ID]*Entry{
			ID(1): {
				ID:        1,
				Title:     "A",
				IsRead:    true,
				Published: &twoWeeksAgo,
				Updated:   nil,
			},
			ID(2): {
				ID:        2,
				Title:     "B",
				IsRead:    false,
				Published: &lastWeek,
				Updated:   &yesterday,
			},
			ID(3): {
				ID:        3,
				Title:     "C",
				IsRead:    false,
				Published: &lastWeek,
				Updated:   &yesterday,
			},
			ID(4): {
				ID:        4,
				Title:     "D",
				IsRead:    false,
				Published: nil,
				Updated:   &threeDaysAgo,
			},
			ID(5): {
				ID:        5,
				Title:     "E",
				IsRead:    false,
				Published: &oneHourAgo,
				Updated:   nil,
			},
			ID(6): {
				ID:        6,
				Title:     "F",
				IsRead:    true,
				Published: nil,
				Updated:   nil,
			},
		},
	}
	slice := f.EntriesSlice()

	want := []string{"E", "B", "C", "D", "A", "F"}
	got := make([]string, len(slice))

	for i, entry := range slice {
		title := entry.Title
		got[i] = title
	}

	a.Len(got, 6)
	a.Equal(want, got)
}
