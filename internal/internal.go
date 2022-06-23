package internal

import (
	"context"

	"github.com/mmcdole/gofeed"

	st "github.com/bow/courier/internal/store"
)

// FeedParser captures the gofeed parser as a pluggable interface.
type FeedParser interface {
	ParseURLWithContext(feedURL string, ctx context.Context) (feed *gofeed.Feed, err error)
}

// FeedStore describes the persistence layer interface.
type FeedStore interface {
	AddFeed(context.Context, *gofeed.Feed, *string, *string, []string) error
	ListFeeds(context.Context) ([]*st.Feed, error)
}
