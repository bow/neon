// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bow/iris/api"
	"github.com/bow/iris/internal/store/opml"
)

const defaultExportTitle = "iris export"

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
	Description sql.NullString
	FeedURL     string
	SiteURL     sql.NullString
	Subscribed  string
	Updated     sql.NullString
	LastPulled  string
	Tags        jsonArrayString
	IsStarred   bool
	Entries     []*Entry
}

func (f *Feed) Proto() (*api.Feed, error) {
	proto := api.Feed{
		Id:          f.ID,
		Title:       f.Title,
		FeedUrl:     f.FeedURL,
		SiteUrl:     unwrapNullString(f.SiteURL),
		Tags:        []string(f.Tags),
		Description: unwrapNullString(f.Description),
		IsStarred:   f.IsStarred,
	}

	var err error

	proto.SubTime, err = toProtoTime(&f.Subscribed)
	if err != nil {
		return nil, err
	}

	proto.UpdateTime, err = toProtoTime(unwrapNullString(f.Updated))
	if err != nil {
		return nil, err
	}

	proto.LastPullTime, err = toProtoTime(&f.LastPulled)
	if err != nil {
		return nil, err
	}

	for _, entry := range f.Entries {
		ep, err := entry.Proto()
		if err != nil {
			return nil, err
		}
		proto.Entries = append(proto.Entries, ep)
	}

	return &proto, nil
}

func (f *Feed) Outline() (*opml.Outline, error) {
	outl := opml.Outline{
		Text:   f.Title,
		Type:   "rss",
		XMLURL: f.FeedURL,
	}
	if f.SiteURL.Valid {
		outl.HTMLURL = pointer(f.SiteURL.String)
	}
	if f.Description.Valid {
		outl.Description = pointer(f.Description.String)
	}
	if len(f.Tags) > 0 {
		outl.Categories = opml.Categories(f.Tags)
	}
	if f.IsStarred {
		outl.IsStarred = &f.IsStarred
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

func NewFeedEditOp(proto *api.EditFeedsRequest_Op) *FeedEditOp {
	return &FeedEditOp{
		ID:          proto.Id,
		Title:       proto.Fields.Title,
		Description: proto.Fields.Description,
		Tags:        &proto.Fields.Tags,
		IsStarred:   proto.Fields.IsStarred,
	}
}

type Entry struct {
	ID          ID
	FeedID      ID
	Title       string
	IsRead      bool
	ExtID       string
	Updated     sql.NullString
	Published   sql.NullString
	Description sql.NullString
	Content     sql.NullString
	URL         sql.NullString
}

func (e *Entry) Proto() (*api.Entry, error) {
	proto := api.Entry{
		Id:          e.ID,
		FeedId:      e.FeedID,
		Title:       e.Title,
		IsRead:      e.IsRead,
		ExtId:       e.ExtID,
		Description: unwrapNullString(e.Description),
		Content:     unwrapNullString(e.Content),
		Url:         unwrapNullString(e.URL),
	}

	var err error

	proto.PubTime, err = toProtoTime(unwrapNullString(e.Published))
	if err != nil {
		return nil, err
	}

	proto.UpdateTime, err = toProtoTime(unwrapNullString(e.Updated))
	if err != nil {
		return nil, err
	}

	return &proto, nil
}

type EntryEditOp struct {
	ID     ID
	IsRead *bool
}

func NewEntryEditOp(proto *api.EditEntriesRequest_Op) *EntryEditOp {
	return &EntryEditOp{ID: proto.Id, IsRead: proto.Fields.IsRead}
}

type Stats struct {
	NumFeeds         uint32
	NumEntries       uint32
	NumEntriesUnread uint32

	RawLastPullTime         string
	RawMostRecentUpdateTime sql.NullString
}

func (s *Stats) LastPullTime() (*time.Time, error) {
	return deserializeTime(&s.RawLastPullTime)
}

func (s *Stats) MostRecentUpdateTime() (*time.Time, error) {
	return deserializeTime(unwrapNullString(s.RawMostRecentUpdateTime))
}

func (s *Stats) Proto() (*api.GetStatsResponse_Stats, error) {
	proto := api.GetStatsResponse_Stats{
		NumFeeds:         s.NumFeeds,
		NumEntries:       s.NumEntries,
		NumEntriesUnread: s.NumEntriesUnread,
	}

	var err error

	proto.LastPullTime, err = toProtoTime(&s.RawLastPullTime)
	if err != nil {
		return nil, err
	}

	proto.MostRecentUpdateTime, err = toProtoTime(unwrapNullString(s.RawMostRecentUpdateTime))
	if err != nil {
		return nil, err
	}

	return &proto, nil
}

// wrapNullString wraps the given string into an sql.NullString value. An empty string input is
// considered a database NULL value.
func wrapNullString(v string) sql.NullString {
	return sql.NullString{String: v, Valid: v != ""}
}

func resolveFeedUpdateTime(feed *gofeed.Feed) *time.Time {
	// Use feed value if defined.
	var latest = feed.UpdatedParsed
	if latest != nil {
		return latest
	}
	// Otherwise try to infer from entries.
	for _, entry := range feed.Items {
		etv := resolveEntryUpdateTime(entry)
		if latest == nil {
			latest = etv
		}
		if latest != nil && etv != nil {
			if etv.After(*latest) {
				latest = etv
			}
		}
	}
	return latest
}

func resolveEntryUpdateTime(entry *gofeed.Item) *time.Time {
	// Use value if defined.
	if tv := entry.UpdatedParsed; tv != nil {
		return tv
	}
	// Otherwise use published time.
	return entry.PublishedParsed
}

func resolveEntryPublishedTime(entry *gofeed.Item) *time.Time {
	// Use value if defined.
	if tv := entry.PublishedParsed; tv != nil {
		return tv
	}
	// Otherwise use update time.
	return entry.UpdatedParsed
}

func serializeTime(tv *time.Time) *string {
	if tv == nil {
		return nil
	}
	ts := tv.UTC().Format(time.RFC3339)
	return &ts
}

func deserializeTime(v *string) (*time.Time, error) {
	if v == nil {
		return nil, nil
	}
	if *v == "" {
		return nil, nil
	}
	pv, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil, err
	}
	upv := pv.UTC()
	return &upv, nil
}

func toProtoTime(v *string) (*timestamppb.Timestamp, error) {
	tv, err := deserializeTime(v)
	if err != nil {
		return nil, err
	}
	if tv == nil {
		return nil, nil
	}
	return timestamppb.New(*tv), nil
}

// unwrapNullString unwraps the given sql.NullString value into a string pointer. If the input value
// is NULL (i.e. its `Valid` field is `false`), `nil` is returned.
func unwrapNullString(v sql.NullString) *string {
	if v.Valid {
		s := v.String
		return &s
	}
	return nil
}

// jsonArrayString is a wrapper type that implements Scan() for database-compatible
// (de)serialization.
type jsonArrayString []string

// Value implements the database valuer interface for serializing into the database.
func (arr *jsonArrayString) Value() (driver.Value, error) {
	if arr == nil {
		return nil, nil
	}
	return json.Marshal([]string(*arr))
}

// Scan implements the database scanner interface for deserialization out of the database.
func (arr *jsonArrayString) Scan(value any) error {
	var bv []byte

	switch v := value.(type) {
	case []byte:
		bv = v
	case string:
		bv = []byte(v)
	default:
		return fmt.Errorf("value of type %T can not be scanned into a string slice", v)
	}

	return json.Unmarshal(bv, arr)
}

func pointer[T any](value T) *T { return &value }
