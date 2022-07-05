package internal

import (
	"context"

	"github.com/mmcdole/gofeed"
)

// AppName returns the application name.
func AppName() string {
	return "courier"
}

// FeedParser captures the gofeed parser as a pluggable interface.
type FeedParser interface {
	ParseURLWithContext(feedURL string, ctx context.Context) (feed *gofeed.Feed, err error)
}
