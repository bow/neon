package internal

import (
	"context"
	"database/sql"

	"github.com/bow/courier/api"
	"github.com/mmcdole/gofeed"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddFeed satisfies the service API.
func (r *rpc) AddFeed(
	ctx context.Context,
	req *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {

	url := req.GetUrl()
	errExistsF := func(url string) error {
		return status.Errorf(codes.AlreadyExists, "feed with URL '%s' already added", url)
	}

	hasFeed, err := r.store.HasFeedURL(ctx, url)
	if err != nil {
		return nil, err
	}
	if hasFeed {
		return nil, errExistsF(url)
	}

	feed, err := r.parser.ParseURL(url)
	if err != nil {
		return nil, err
	}

	err = r.store.AddFeed(ctx, feed, req.Title, req.Description, req.GetCategories())
	if err != nil {
		if isUniqueErr(err, "UNIQUE constraint failed: feeds.xml_url") {
			return nil, errExistsF(feed.FeedLink)
		}
		return nil, err
	}

	rsp := api.AddFeedResponse{}

	return &rsp, nil
}

// HasFeedURL checks if a feed with the given URL already exists in the database.
func (s *sqliteStore) HasFeedURL(ctx context.Context, url string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("FeedStore.HasFeedURL")

	var exists bool
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		sql1 := `SELECT EXISTS (SELECT id FROM feeds WHERE xml_url = ?)`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return fail(err)
		}
		defer stmt1.Close()

		if err := stmt1.QueryRowContext(ctx, url).Scan(&exists); err != nil {
			return fail(err)
		}

		return nil
	}

	if err := s.withTx(ctx, dbFunc, nil); err != nil {
		return exists, err
	}

	return exists, nil
}

// AddFeed adds the given feed into the database.
func (s *sqliteStore) AddFeed(
	ctx context.Context,
	feed *gofeed.Feed,
	title *string,
	desc *string,
	categories []string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("FeedStore.AddFeed")

	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		sql1 := `INSERT INTO feeds(title, description, xml_url, html_url) VALUES (?, ?, ?, ?)`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return fail(err)
		}
		defer stmt1.Close()

		res, err := stmt1.ExecContext(
			ctx,
			nullIf(resolve(title, feed.Title), textEmpty),
			nullIf(resolve(desc, feed.Description), textEmpty),
			feed.FeedLink,
			nullIf(feed.Link, textEmpty),
		)
		if err != nil {
			return fail(err)
		}

		feedDBID, err := res.LastInsertId()
		if err != nil {
			return fail(err)
		}

		err = s.addFeedCategories(ctx, tx, DBID(feedDBID), categories)
		if err != nil {
			return fail(err)
		}

		return nil
	}

	return s.withTx(ctx, dbFunc, nil)
}

func (s *sqliteStore) addFeedCategories(
	ctx context.Context,
	tx *sql.Tx,
	feedDBID DBID,
	cats []string,
) error {

	sql1 := `INSERT OR IGNORE INTO feed_categories(name) VALUES (?)`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()
	for _, cat := range cats {
		_, err = stmt1.ExecContext(ctx, cat)
		if err != nil {
			return err
		}
	}

	sql2 := `SELECT id FROM feed_categories WHERE name = ?`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()
	ids := make(map[string]DBID)
	for _, cat := range cats {
		if _, exists := ids[cat]; exists {
			continue
		}
		var id DBID
		row := stmt2.QueryRowContext(ctx, cat)
		if err = row.Scan(&id); err != nil {
			return err
		}
		ids[cat] = id
	}

	sql3 := `INSERT INTO feeds_x_feed_categories(feed_id, feed_category_id) VALUES (?, ?)`
	stmt3, err := tx.PrepareContext(ctx, sql3)
	if err != nil {
		return err
	}
	defer stmt3.Close()

	for _, catDBID := range ids {
		if _, err := stmt3.ExecContext(ctx, feedDBID, catDBID); err != nil {
			return err
		}
	}

	return nil
}
