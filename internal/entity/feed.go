// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

import (
	"time"

	"github.com/bow/neon/internal/opml"
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
	Entries     []*Entry
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
