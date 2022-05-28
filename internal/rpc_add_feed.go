package internal

import (
	"context"

	"github.com/bow/courier/api"
)

// AddFeed satisfies the service API.
func (svc *service) AddFeed(
	_ context.Context,
	_ *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {
	return nil, svc.store.AddFeed(nil)
}
