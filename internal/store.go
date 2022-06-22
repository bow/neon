package internal

import (
	"context"

	"github.com/mmcdole/gofeed"

	st "github.com/bow/courier/internal/store"
)

type FeedStore interface {
	AddFeed(context.Context, *gofeed.Feed, *string, *string, []string) error
	ListFeeds(context.Context) ([]*st.Feed, error)
}
