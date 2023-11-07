// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/bow/iris/api"
	"github.com/bow/iris/internal"
)

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

type EntryEditOp struct {
	ID     ID
	IsRead *bool
}

func NewEntryEditOp(proto *api.EditEntriesRequest_Op) *EntryEditOp {
	return &EntryEditOp{ID: proto.Id, IsRead: proto.Fields.IsRead}
}

type feedRecord struct {
	id          ID
	title       string
	description sql.NullString
	feedURL     string
	siteURL     sql.NullString
	subscribed  string
	lastPulled  string
	updated     sql.NullString
	isStarred   bool
	tags        jsonArrayString
	entries     []*entryRecord
}

type entryRecord struct {
	id          ID
	feedID      ID
	title       string
	isRead      bool
	extID       string
	updated     sql.NullString
	published   sql.NullString
	description sql.NullString
	content     sql.NullString
	url         sql.NullString
}

type statsAggregateRecord struct {
	numFeeds             uint32
	numEntries           uint32
	numEntriesUnread     uint32
	lastPullTime         string
	mostRecentUpdateTime sql.NullString
}

func toFeedID(raw string) (ID, error) {
	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, FeedNotFoundError{ID: raw}
	}
	return ID(id), nil
}

func toFeed(rec *feedRecord) (*internal.Feed, error) {

	subt, err := deserializeTime(&rec.subscribed)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize Feed.Subscribed time: %w", err)
	}
	lpt, err := deserializeTime(&rec.lastPulled)
	if err != nil {
		return nil, err
	}
	var upt *time.Time
	if v := fromNullString(rec.updated); v != nil {
		upt, err = deserializeTime(v)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize Feed.Updated time: %w", err)
		}
	}
	entries, err := toEntries(rec.entries)
	if err != nil {
		return nil, err
	}

	feed := internal.Feed{
		ID:          rec.id,
		Title:       rec.title,
		Description: fromNullString(rec.description),
		FeedURL:     rec.feedURL,
		SiteURL:     fromNullString(rec.siteURL),
		Subscribed:  *subt,
		LastPulled:  *lpt,
		Updated:     upt,
		IsStarred:   rec.isStarred,
		Tags:        []string(rec.tags),
		Entries:     entries,
	}
	return &feed, nil
}

func toFeeds(recs []*feedRecord) ([]*internal.Feed, error) {

	feeds := make([]*internal.Feed, len(recs))
	for i, rec := range recs {
		feed, err := toFeed(rec)
		if err != nil {
			return nil, err
		}
		feeds[i] = feed
	}

	return feeds, nil
}

func toEntry(rec *entryRecord) (*internal.Entry, error) {

	ut, err := deserializeTime(fromNullString(rec.updated))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize Entry.Updated time: %w", err)
	}
	pt, err := deserializeTime(fromNullString(rec.published))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize Entry.Published time: %w", err)
	}

	entry := internal.Entry{
		ID:          rec.id,
		FeedID:      rec.feedID,
		Title:       rec.title,
		IsRead:      rec.isRead,
		ExtID:       rec.extID,
		Updated:     ut,
		Published:   pt,
		Description: fromNullString(rec.description),
		Content:     fromNullString(rec.content),
		URL:         fromNullString(rec.url),
	}

	return &entry, nil
}

func toEntries(recs []*entryRecord) ([]*internal.Entry, error) {

	entries := make([]*internal.Entry, len(recs))
	for i, rec := range recs {
		entry, err := toEntry(rec)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}

	return entries, nil
}

func toStats(aggr *statsAggregateRecord) (*internal.Stats, error) {

	lpt, err := deserializeTime(&aggr.lastPullTime)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize Stats.LastPullTime: %w", err)
	}

	mrut, err := deserializeTime(fromNullString(aggr.mostRecentUpdateTime))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize Stats.MostRecentUpdateTime: %w", err)
	}

	stats := internal.Stats{
		NumFeeds:             aggr.numFeeds,
		NumEntries:           aggr.numEntries,
		NumEntriesUnread:     aggr.numEntriesUnread,
		LastPullTime:         lpt,
		MostRecentUpdateTime: mrut,
	}

	return &stats, nil
}

// toNullString wraps the given string into an sql.NullString value. An empty string input is
// considered a database NULL value.
func toNullString(v string) sql.NullString {
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

// fromNullString unwraps the given sql.NullString value into a string pointer. If the input value
// is NULL (i.e. its `Valid` field is `false`), `nil` is returned.
func fromNullString(v sql.NullString) *string {
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
