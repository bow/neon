package internal

import (
	"context"

	"github.com/bow/courier/api"
	"github.com/mmcdole/gofeed"
)

// AddFeed satisfies the service API.
func (svc *service) AddFeed(
	_ context.Context,
	req *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {

	var (
		err  error
		feed *gofeed.Feed
	)

	if feed, err = svc.parser.ParseURL(req.GetUrl()); err != nil {
		return nil, err
	}

	if err = svc.store.AddFeed(feed); err != nil {
		return nil, err
	}

	rsp := api.AddFeedResponse{}

	return &rsp, nil
}
