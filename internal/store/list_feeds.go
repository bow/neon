package store

import (
	"context"
	"database/sql"
)

func (s *SQLite) ListFeeds(ctx context.Context) ([]*Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	feeds := make([]*Feed, 0)
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		ifeeds, err := getAllFeeds(ctx, tx)
		if err != nil {
			return err
		}
		for _, ifeed := range ifeeds {
			ifeed := ifeed
			entries, err := getAllFeedEntries(ctx, tx, ifeed.DBID)
			if err != nil {
				return err
			}
			ifeed.Entries = entries
		}
		feeds = ifeeds

		return nil
	}

	fail := failF("SQLite.ListFeeds")

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, fail(err)
	}
	return feeds, nil
}

func getAllFeeds(ctx context.Context, tx *sql.Tx) ([]*Feed, error) {

	sql1 := `
		SELECT
			f.id AS id,
			f.title AS title,
			f.description AS description,
			f.feed_url AS feed_url,
			f.site_url AS site_url,
			f.is_starred AS is_starred,
			f.subscription_time AS subscription_time,
			f.update_time AS update_time,
			json_group_array(fc.name) FILTER (WHERE fc.name IS NOT NULL) AS tags
		FROM
			feeds f
			LEFT JOIN feeds_x_feed_tags fxfc ON fxfc.feed_id = f.id
			LEFT JOIN feed_tags fc ON fxfc.feed_tag_id = fc.id
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
			&feed.IsStarred,
			&feed.Subscribed,
			&feed.Updated,
			&feed.Tags,
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

func getAllFeedEntries(ctx context.Context, tx *sql.Tx, feedDBID DBID) ([]*Entry, error) {

	sql1 := `
		SELECT
			e.id AS id,
			e.feed_id AS feed_id,
			e.title AS title,
			e.is_read AS is_read,
			e.external_id AS ext_id,
			e.description AS description,
			e.content AS content,
			e.url AS url,
			e.update_time AS update_time,
			e.publication_time AS publication_time
		FROM
			entries e
		WHERE
			e.feed_id = ?
		ORDER BY
			COALESCE(e.update_time, e.publication_time) DESC
`
	scanRow := func(rows *sql.Rows) (*Entry, error) {
		var entry Entry
		if err := rows.Scan(
			&entry.DBID,
			&entry.FeedDBID,
			&entry.Title,
			&entry.IsRead,
			&entry.ExtID,
			&entry.Description,
			&entry.Content,
			&entry.URL,
			&entry.Updated,
			&entry.Published,
		); err != nil {
			return nil, err
		}
		return &entry, nil
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	rows, err := stmt1.QueryContext(ctx, feedDBID)
	if err != nil {
		return nil, err
	}

	entries := make([]*Entry, 0)
	for rows.Next() {
		entry, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
