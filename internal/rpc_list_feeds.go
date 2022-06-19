package internal

import (
	"context"
	"database/sql"

	"github.com/bow/courier/api"
)

// ListFeeds satisfies the service API.
func (r *rpc) ListFeeds(
	ctx context.Context,
	_ *api.ListFeedsRequest,
) (*api.ListFeedsResponse, error) {

	feeds, err := r.store.ListFeeds(ctx)
	if err != nil {
		return nil, err
	}

	rsp := api.ListFeedsResponse{}
	for _, feed := range feeds {
		rsp.Feeds = append(
			rsp.Feeds,
			&api.Feed{
				Id:  int32(feed.DBID),
				Url: feed.FeedLink,
			},
		)
	}

	return &rsp, nil
}

func (s *sqliteStore) ListFeeds(ctx context.Context) ([]*Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	feeds := make([]*Feed, 0)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		// TODO: Full implementation.
		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)

	return feeds, err
}
