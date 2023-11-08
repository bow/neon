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

	"github.com/bow/iris/internal"
)

type feedRecord struct {
	id          ID
	title       string
	description sql.NullString
	feedURL     string
	siteURL     sql.NullString
	subscribed  time.Time
	lastPulled  time.Time
	updated     sql.NullTime
	isStarred   bool
	tags        jsonArrayString
	entries     []*entryRecord
}

func (rec *feedRecord) feed() *internal.Feed {
	return &internal.Feed{
		ID:          rec.id,
		Title:       rec.title,
		Description: fromNullString(rec.description),
		FeedURL:     rec.feedURL,
		SiteURL:     fromNullString(rec.siteURL),
		Subscribed:  rec.subscribed,
		LastPulled:  rec.lastPulled,
		Updated:     fromNullTime(rec.updated),
		IsStarred:   rec.isStarred,
		Tags:        []string(rec.tags),
		Entries:     entryRecords(rec.entries).entries(),
	}
}

type feedRecords []*feedRecord

func (recs feedRecords) feeds() []*internal.Feed {

	feeds := make([]*internal.Feed, len(recs))
	for i, rec := range recs {
		feeds[i] = rec.feed()
	}

	return feeds
}

type entryRecord struct {
	id          ID
	feedID      ID
	title       string
	isRead      bool
	extID       string
	updated     sql.NullTime
	published   sql.NullTime
	description sql.NullString
	content     sql.NullString
	url         sql.NullString
}

func (rec *entryRecord) entry() *internal.Entry {
	return &internal.Entry{
		ID:          rec.id,
		FeedID:      rec.feedID,
		Title:       rec.title,
		IsRead:      rec.isRead,
		ExtID:       rec.extID,
		Updated:     fromNullTime(rec.updated),
		Published:   fromNullTime(rec.published),
		Description: fromNullString(rec.description),
		Content:     fromNullString(rec.content),
		URL:         fromNullString(rec.url),
	}
}

type entryRecords []*entryRecord

func (recs entryRecords) entries() []*internal.Entry {

	entries := make([]*internal.Entry, len(recs))
	for i, rec := range recs {
		entries[i] = rec.entry()
	}

	return entries
}

type statsAggregateRecord struct {
	numFeeds             uint32
	numEntries           uint32
	numEntriesUnread     uint32
	lastPullTime         time.Time
	mostRecentUpdateTime sql.NullTime
}

func (aggr *statsAggregateRecord) stats() *internal.Stats {

	var mrut *time.Time
	if aggr.mostRecentUpdateTime.Valid {
		mrut = &aggr.mostRecentUpdateTime.Time
	}

	stats := internal.Stats{
		NumFeeds:             aggr.numFeeds,
		NumEntries:           aggr.numEntries,
		NumEntriesUnread:     aggr.numEntriesUnread,
		LastPullTime:         &aggr.lastPullTime,
		MostRecentUpdateTime: mrut,
	}

	return &stats
}

func toFeedID(raw string) (ID, error) {
	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, FeedNotFoundError{ID: raw}
	}
	return ID(id), nil
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

// fromNullString unwraps the given sql.NullString value into a string pointer. If the input value
// is NULL (i.e. its `Valid` field is `false`), `nil` is returned.
func fromNullString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}

func fromNullTime(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	return &v.Time
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
