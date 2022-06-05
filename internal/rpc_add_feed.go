package internal

import (
	"context"

	"github.com/bow/courier/api"
)

// AddFeed satisfies the service API.
func (svc *service) AddFeed(
	ctx context.Context,
	req *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {

	feed, err := svc.parser.ParseURL(req.GetUrl())
	if err != nil {
		return nil, err
	}

	err = svc.store.AddFeed(ctx, feed, req.Title, req.Description, req.GetCategories())
	if err != nil {
		return nil, err
	}

	rsp := api.AddFeedResponse{}

	return &rsp, nil
}
