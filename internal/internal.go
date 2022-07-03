package internal

import (
	"context"

	"github.com/mmcdole/gofeed"

	"github.com/bow/courier/internal/store"
)

// AppName returns the application name.
func AppName() string {
	return "courier"
}

// FeedParser captures the gofeed parser as a pluggable interface.
type FeedParser interface {
	ParseURLWithContext(feedURL string, ctx context.Context) (feed *gofeed.Feed, err error)
}

// FeedStore describes the persistence layer interface.
type FeedStore interface {
	AddFeed(
		ctx context.Context,
		feed *gofeed.Feed,
		title *string,
		desc *string,
		tags []string,
		isStarred bool,
	) (addedFeed *store.Feed, err error)

	EditFeeds(ctx context.Context, ops []*store.FeedEditOp) (feeds []*store.Feed, err error)

	ListFeeds(ctx context.Context) (feeds []*store.Feed, err error)

	DeleteFeeds(ctx context.Context, ids []store.DBID) (err error)

	EditEntries(ctx context.Context, ops []*store.EntryEditOp) (entries []*store.Entry, err error)

	ExportOPML(ctx context.Context) (payload []byte, err error)
}
