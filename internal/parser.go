package internal

import (
	"context"
	"io"

	"github.com/mmcdole/gofeed"
)

// FeedParser captures the gofeed parser as a pluggable interface.
type FeedParser interface {
	Parse(io.Reader) (*gofeed.Feed, error)
	ParseURL(string) (*gofeed.Feed, error)
	ParseURLWithContext(feedURL string, ctx context.Context) (feed *gofeed.Feed, err error)
}
