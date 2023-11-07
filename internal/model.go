// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"time"

	"github.com/bow/iris/internal/opml"
)

type ID = uint32

type Subscription []*Feed

func (sub Subscription) Export(title *string) ([]byte, error) {
	et := defaultExportTitle
	if title != nil {
		et = *title
	}
	doc := opml.New(et, time.Now())
	for _, feed := range sub {
		if err := doc.AddOutline(feed); err != nil {
			return nil, err
		}
	}
	return doc.XML()
}

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

type Entry struct {
	ID          ID
	FeedID      ID
	Title       string
	IsRead      bool
	ExtID       string
	Updated     *time.Time
	Published   *time.Time
	Description *string
	Content     *string
	URL         *string
}

type Stats struct {
	NumFeeds             uint32
	NumEntries           uint32
	NumEntriesUnread     uint32
	LastPullTime         *time.Time
	MostRecentUpdateTime *time.Time
}

const defaultExportTitle = "iris export"
