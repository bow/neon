package store

import (
	"context"

	"github.com/mmcdole/gofeed"
)

// FeedParser captures the gofeed parser as a pluggable interface.
type FeedParser interface {
	ParseURLWithContext(feedURL string, ctx context.Context) (feed *gofeed.Feed, err error)
}
