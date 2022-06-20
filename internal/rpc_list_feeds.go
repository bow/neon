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
		// TODO: Use gRPC INTERNAL error for this.
		proto, err := feed.Proto()
		if err != nil {
			return nil, err
		}
		rsp.Feeds = append(rsp.Feeds, proto)
	}

	return &rsp, nil
}

func (s *sqliteStore) ListFeeds(ctx context.Context) ([]*Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("sqliteStore.ListFeeds")

	feeds := make([]*Feed, 0)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		ifeeds, err := s.getAllFeeds(ctx, tx)
		if err != nil {
			return fail(err)
		}
		for _, ifeed := range ifeeds {
			ifeed := ifeed
			if err := s.populateEntries(ctx, tx, ifeed); err != nil {
				return fail(err)
			}
		}
		feeds = ifeeds

		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)

	return feeds, err
}

func (s *sqliteStore) getAllFeeds(ctx context.Context, tx *sql.Tx) ([]*Feed, error) {

	sql1 := `
		SELECT
			f.id AS id,
			f.title AS title,
			f.description AS description,
			f.feed_url AS feed_url,
			f.site_url AS site_url,
			f.subscription_time AS subscription_time,
			f.update_time AS update_time,
			json_group_array(fc.name) FILTER (WHERE fc.name IS NOT NULL) AS categories
		FROM
			feeds f
			LEFT JOIN feeds_x_feed_categories fxfc ON fxfc.feed_id = f.id
			LEFT JOIN feed_categories fc ON fxfc.feed_category_id = fc.id
		GROUP BY
			f.id
		ORDER BY
			COALESCE(f.update_time, f.subscription_time) DESC
`
	scanRow := func(rows *sql.Rows) (*Feed, error) {
		var feed Feed
		if err := rows.Scan(
			&feed.DBID,
			&feed.Title,
			&feed.Description,
			&feed.FeedURL,
			&feed.SiteURL,
			&feed.Subscribed,
			&feed.Updated,
			&feed.Categories,
		); err != nil {
			return nil, err
		}
		return &feed, nil
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	rows, err := stmt1.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	feeds := make([]*Feed, 0)
	for rows.Next() {
		feed, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

func (s *sqliteStore) populateEntries(_ context.Context, _ *sql.Tx, _ *Feed) error {
	return nil
}
