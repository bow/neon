// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import "time"

type Stats struct {
	NumFeeds             uint32
	NumEntries           uint32
	NumEntriesUnread     uint32
	LastPullTime         *time.Time
	MostRecentUpdateTime *time.Time
}
