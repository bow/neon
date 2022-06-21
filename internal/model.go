package internal

import (
	"database/sql"
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
	Entries     []*Entry
}

func (f *Feed) Proto() (*api.Feed, error) {
	proto := api.Feed{
		Id:          int32(f.DBID),
		Title:       f.Title,
		FeedUrl:     f.FeedURL,
		SiteUrl:     UnwrapNullString(f.SiteURL),
		Categories:  []string(f.Categories),
		Description: UnwrapNullString(f.Description),
	}

	stv, err := deserializeTime(&f.Subscribed)
	if err != nil {
		return nil, err
	}
	proto.SubscriptionTime = timestamppb.New(*stv)

	utv, err := deserializeTime(UnwrapNullString(f.Updated))
	if err != nil {
		return nil, err
	}
	if utv != nil {
		proto.UpdateTime = timestamppb.New(*utv)
	}

	return &proto, nil
}

type Entry struct {
	Title string
}

func resolveFeedUpdateTime(feed *gofeed.Feed) *time.Time {
	// Use feed value if defined.
	var latest *time.Time = feed.UpdatedParsed
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
	pv, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil, err
	}
	upv := pv.UTC()
	return &upv, nil
}
