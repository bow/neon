// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import "time"

type Entry struct {
	ID           ID
	FeedID       ID
	Title        string
	IsRead       bool
	IsBookmarked bool
	ExtID        string
	Updated      *time.Time
	Published    *time.Time
	Description  *string
	Content      *string
	URL          *string
}

type EntryEditOp struct {
	ID           ID
	IsRead       *bool
	IsBookmarked *bool
}
