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
				Id:      int32(feed.DBID),
				FeedUrl: feed.inner.FeedLink,
			},
		)
	}

	return &rsp, nil
}

func (s *sqliteStore) ListFeeds(ctx context.Context) ([]*Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("sqliteStore.ListFeeds")

	feeds := make([]*Feed, 0)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		sql1 := `
			SELECT
				id, feed_url
			FROM
				feeds
			ORDER BY
				COALESCE(update_time, subscription_time) DESC
`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return fail(err)
		}
		defer stmt1.Close()

		rows, err := stmt1.QueryContext(ctx)
		if err != nil {
			return fail(err)
		}

		for rows.Next() {
			var feed Feed
			if err := rows.Scan(&feed.DBID, &feed.inner.FeedLink); err != nil {
				return fail(err)
			}
			feeds = append(feeds, &feed)
		}

		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)

	return feeds, err
}
