package internal

import (
	"context"

	"github.com/mmcdole/gofeed"

	"github.com/bow/courier/internal/store"
)

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
		categories []string,
	) (err error)

	ListFeeds(ctx context.Context) (feeds []*store.Feed, err error)

	SetEntryFields(
		ctx context.Context,
		setOps []*store.EntrySetOp,
	) (entries []*store.Entry, err error)
}
