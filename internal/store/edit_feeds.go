package store

import (
	"context"
	"database/sql"
	"fmt"
)

// EditFeed updates fields of an feed.
func (s *SQLite) EditFeeds(
	ctx context.Context,
	ops []*FeedEditOp,
) ([]*Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("SQLite.EditFeed")

	updateFunc := func(ctx context.Context, tx *sql.Tx, op *FeedEditOp) (*Feed, error) {
		if err := updateFeedTitle(ctx, tx, op.DBID, op.Title); err != nil {
			return nil, err
		}
		if err := updateFeedDescription(ctx, tx, op.DBID, op.Description); err != nil {
			return nil, err
		}
		if err := s.updateFeedCategories(ctx, tx, op.DBID, op.Categories); err != nil {
			return nil, err
		}
		return getFeed(ctx, tx, op.DBID)
	}

	var entries = make([]*Feed, len(ops))
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		for i, op := range ops {
			feed, err := updateFunc(ctx, tx, op)
			if err != nil {
				return fail(err)
			}
			entries[i] = feed
		}
		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func getFeed(ctx context.Context, tx *sql.Tx, feedDBID DBID) (*Feed, error) {

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
		WHERE
			f.id = ?
		GROUP BY
			f.id
		ORDER BY
			COALESCE(f.update_time, f.subscription_time) DESC
`
	scanRow := func(row *sql.Row) (*Feed, error) {
		var feed Feed
		if err := row.Scan(
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

	return scanRow(stmt1.QueryRowContext(ctx, feedDBID))
}

func updateFeedField[T any](columnName string) func(context.Context, *sql.Tx, DBID, *T) error {

	return func(ctx context.Context, tx *sql.Tx, feedDBID DBID, fieldValue *T) error {

		if fieldValue == nil {
			return nil
		}

		sql1 := `UPDATE feeds SET ` + columnName + ` = $2 WHERE id = $1 RETURNING id`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return err
		}
		defer stmt1.Close()

		var updatedID DBID
		err = stmt1.QueryRowContext(ctx, feedDBID, fieldValue).Scan(&updatedID)
		if err != nil {
			return err
		}
		if updatedID == 0 {
			// TODO: Wrap in proper gRPC errors.
			return fmt.Errorf("feed id %d does not exist", updatedID)
		}
		return nil
	}
}

var (
	updateFeedTitle       = updateFeedField[string]("title")
	updateFeedDescription = updateFeedField[string]("description")
)

func (s *SQLite) updateFeedCategories(
	ctx context.Context,
	tx *sql.Tx,
	feedDBID DBID,
	categories *[]string,
) error {

	if categories == nil {
		return nil
	}

	sql1 := `DELETE FROM feeds_x_feed_categories WHERE feed_id = ?`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()

	if _, err = stmt1.ExecContext(ctx); err != nil {
		return err
	}

	if err = s.addFeedCategories(ctx, tx, feedDBID, *categories); err != nil {
		return err
	}

	sql2 := `
		DELETE
			feed_categories
		WHERE
			id IN (
				SELECT
					fc.id
				FROM
					feed_categories fc
					LEFT JOIN feeds_x_feed_categories fxfc ON fxfc.feed_category_id = fc.id
				WHERE
					fxfc.feed_id IS NULL
			)
	`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()

	return nil
}
