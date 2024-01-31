// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import (
	"time"

	"github.com/bow/neon/internal/opml"
	"github.com/bow/neon/internal/sliceutil"
)

type Feed struct {
	ID          ID
	Title       string
	Description *string
	FeedURL     string
	SiteURL     *string
	Subscribed  time.Time
	LastPulled  time.Time
	Updated     *time.Time
	IsStarred   bool
	Tags        []string
	Entries     map[ID]*Entry
}

func (f *Feed) NumEntriesTotal() int {
	return len(f.Entries)
}

func (f *Feed) NumEntriesRead() int {
	var n int
	for _, entry := range f.Entries {
		if entry.IsRead {
			n++
		}
	}
	return n
}

func (f *Feed) NumEntriesUnread() int {
	return f.NumEntriesTotal() - f.NumEntriesRead()
}

// EntriesSlice returns a slice of entries sorted by read status (unread first),
// update date (oldest first), and published date (oldest first).
func (f *Feed) EntriesSlice() []*Entry { // nolint:revive
	sortDate := func(e *Entry) *time.Time {
		if e.Updated != nil {
			return e.Updated
		}
		if e.Published != nil {
			return e.Published
		}
		return nil
	}
	isRead := func(e1, e2 *Entry) int {
		if e1.IsRead && !e2.IsRead {
			return 1
		}
		if !e1.IsRead && e2.IsRead {
			return -1
		}
		return 0
	}
	date := func(e1, e2 *Entry) int {
		d1 := sortDate(e1)
		d2 := sortDate(e2)
		if d1 != nil && d2 != nil {
			if d1.Before(*d2) {
				return 1
			}
			if d2.Before(*d1) {
				return -1
			}
			return 0
		}
		if d1 != nil {
			return -1
		}
		if d2 != nil {
			return 1
		}
		return 0
	}

	entries := make([]*Entry, 0)
	for _, entry := range f.Entries {
		entries = append(entries, entry)
	}

	sliceutil.Ordered[*Entry]().
		By(isRead, date).
		Sort(entries)

	return entries
}

func (f *Feed) Outline() (*opml.Outline, error) {
	outl := opml.Outline{
		Text:        f.Title,
		Type:        "rss",
		XMLURL:      f.FeedURL,
		HTMLURL:     f.SiteURL,
		Description: f.Description,
	}
	if v := f.IsStarred; v {
		outl.IsStarred = &v
	}
	if len(f.Tags) > 0 {
		outl.Categories = opml.Categories(f.Tags)
	}

	return &outl, nil
}

type FeedEditOp struct {
	ID          ID
	Title       *string
	Description *string
	Tags        *[]string
	IsStarred   *bool
}
