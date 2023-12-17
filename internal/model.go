// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bow/lens/api"
	"github.com/bow/lens/internal/opml"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ID = uint32

func ToFeedIDs(raw []string) ([]ID, error) {
	ids := make([]ID, 0)
	for _, item := range raw {
		id, err := ToFeedID(item)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func ToFeedID(raw string) (ID, error) {
	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, FeedNotFoundError{ID: raw}
	}
	return ID(id), nil
}

type Subscription struct {
	Title *string
	Feeds []*Feed
}

func NewSubscriptionFromRawOPML(payload []byte) (*Subscription, error) {

	doc, err := opml.Parse(payload)
	if err != nil {
		return nil, err
	}

	return NewSubscriptionFromOPML(doc)
}

func NewSubscriptionFromOPML(doc *opml.Doc) (*Subscription, error) {

	feeds := make([]*Feed, len(doc.Body.Outlines))
	for i, outl := range doc.Body.Outlines {
		if outl.Text == "" {
			return nil, fmt.Errorf(
				"missing title for feed with URL=%s in OPML document", outl.XMLURL,
			)
		}
		feed := Feed{
			Title:       outl.Text,
			Description: outl.Description,
			FeedURL:     outl.XMLURL,
			SiteURL:     outl.HTMLURL,
			Tags:        outl.Categories,

			ID:         0,
			Subscribed: time.Time{},
			LastPulled: time.Time{},
		}
		if star := outl.IsStarred; star != nil {
			feed.IsStarred = *star
		}
		feeds[i] = &feed
	}

	sub := Subscription{Feeds: feeds}

	return &sub, nil
}

func (sub *Subscription) Export() ([]byte, error) {
	var et = defaultExportTitle
	if sub.Title != nil {
		et = *sub.Title
	}

	doc := opml.New(et, time.Now())
	for _, feed := range sub.Feeds {
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

type FeedEditOp struct {
	ID          ID
	Title       *string
	Description *string
	Tags        *[]string
	IsStarred   *bool
}

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

type Stats struct {
	NumFeeds             uint32
	NumEntries           uint32
	NumEntriesUnread     uint32
	LastPullTime         *time.Time
	MostRecentUpdateTime *time.Time
}

func FromStatsPb(pb *api.GetStatsResponse_Stats) *Stats {
	return &Stats{
		NumFeeds:             pb.GetNumFeeds(),
		NumEntries:           pb.GetNumEntries(),
		NumEntriesUnread:     pb.GetNumEntriesUnread(),
		LastPullTime:         FromTimestampPb(pb.GetLastPullTime()),
		MostRecentUpdateTime: FromTimestampPb(pb.GetMostRecentUpdateTime()),
	}
}

func FromTimestampPb(pb *timestamppb.Timestamp) *time.Time {
	if pb == nil {
		return nil
	}
	v := pb.AsTime()
	return &v
}

// PullResult is a container for a pull operation.
type PullResult struct {
	status PullStatus
	url    *string
	feed   *Feed
	err    error
}

func NewPullResultFromFeed(url *string, feed *Feed) PullResult {
	return PullResult{status: PullSuccess, url: url, feed: feed}
}

func NewPullResultFromError(url *string, err error) PullResult {
	return PullResult{status: PullFail, url: url, err: err}
}

func (msg PullResult) Feed() *Feed {
	if msg.status == PullSuccess {
		return msg.feed
	}
	return nil
}

func (msg PullResult) Error() error {
	if msg.status == PullFail {
		return msg.err
	}
	return nil
}

func (msg PullResult) URL() string {
	if msg.url != nil {
		return *msg.url
	}
	return ""
}

func (msg *PullResult) SetError(err error) {
	msg.err = err
}

func (msg *PullResult) SetStatus(status PullStatus) {
	msg.status = status
}

type PullStatus int

const (
	PullSuccess PullStatus = iota
	PullFail
)

const defaultExportTitle = "lens export"
