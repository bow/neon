package store

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bow/courier/api"
	"github.com/mmcdole/gofeed"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Feed struct {
	DBID        DBID
	Title       string
	Description sql.NullString
	FeedURL     string
	SiteURL     sql.NullString
	Subscribed  string
	Updated     sql.NullString
	Categories  jsonArrayString
	IsStarred   bool
	Entries     []*Entry
}

func (f *Feed) Proto() (*api.Feed, error) {
	proto := api.Feed{
		Id:          int32(f.DBID),
		Title:       f.Title,
		FeedUrl:     f.FeedURL,
		SiteUrl:     unwrapNullString(f.SiteURL),
		Categories:  []string(f.Categories),
		Description: unwrapNullString(f.Description),
		IsStarred:   f.IsStarred,
	}

	var err error

	proto.SubscriptionTime, err = toProtoTime(&f.Subscribed)
	if err != nil {
		return nil, err
	}

	proto.UpdateTime, err = toProtoTime(unwrapNullString(f.Updated))
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

type FeedEditOp struct {
	DBID        DBID
	Title       *string
	Description *string
	Categories  *[]string
}

func NewFeedEditOp(proto *api.EditFeedsRequest_Op) *FeedEditOp {
	return &FeedEditOp{
		DBID:        DBID(proto.Id),
		Title:       proto.Fields.Title,
		Description: proto.Fields.Description,
		Categories:  &proto.Fields.Categories,
	}
}

type Entry struct {
	DBID        DBID
	FeedDBID    DBID
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
		Id:          int32(e.DBID),
		FeedId:      int32(e.FeedDBID),
		Title:       e.Title,
		IsRead:      e.IsRead,
		ExtId:       e.ExtID,
		Description: unwrapNullString(e.Description),
		Content:     unwrapNullString(e.Content),
		Url:         unwrapNullString(e.URL),
	}

	var err error

	proto.PublicationTime, err = toProtoTime(unwrapNullString(e.Published))
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
	DBID   DBID
	IsRead *bool
}

func NewEntryEditOp(proto *api.EditEntriesRequest_Op) *EntryEditOp {
	return &EntryEditOp{DBID: DBID(proto.Id), IsRead: proto.Fields.IsRead}
}

// WrapNullString wraps the given string into an sql.NullString value. An empty string input is
// considered a database NULL value.
func WrapNullString(v string) sql.NullString {
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

func DeserializeTime(v *string) (*time.Time, error) {
	if v == nil {
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
	tv, err := DeserializeTime(v)
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
