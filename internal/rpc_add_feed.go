package internal

import (
	"context"
	"database/sql"

	"github.com/bow/courier/api"
	"github.com/mmcdole/gofeed"
)

// AddFeed satisfies the service API.
func (r *rpc) AddFeed(
	ctx context.Context,
	req *api.AddFeedRequest,
) (*api.AddFeedResponse, error) {

	feed, err := r.parser.ParseURLWithContext(req.GetUrl(), ctx)
	if err != nil {
		return nil, err
	}

	err = r.store.AddFeed(ctx, feed, req.Title, req.Description, req.GetCategories())
	if err != nil {
		return nil, err
	}

	rsp := api.AddFeedResponse{}

	return &rsp, nil
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

	fail := failF("sqliteStore.AddFeed")

	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		sql1 := `
			INSERT INTO
				feeds(title, description, feed_url, site_url, update_time)
				VALUES (?, ?, ?, ?, ?)
`
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
			serializeTime(resolveFeedUpdateTime(feed)),
		)
		var feedDBID int64
		if err == nil {
			feedDBID, err = res.LastInsertId()
			if err != nil {
				return fail(err)
			}
		} else {
			if !isUniqueErr(err, "UNIQUE constraint failed: feeds.feed_url") {
				return fail(err)
			}
			if ierr := tx.QueryRowContext(
				ctx,
				`SELECT id FROM feeds WHERE feed_url = ?`,
				feed.FeedLink,
			).Scan(&feedDBID); ierr != nil {
				return fail(ierr)
			}
			// TODO: Add and combine with proper update call.
			if ierr := s.upsertEntries(ctx, tx, DBID(feedDBID), feed.Items); ierr != nil {
				return fail(ierr)
			}
		}

		err = s.addFeedCategories(ctx, tx, DBID(feedDBID), categories)
		if err != nil {
			return fail(err)
		}

		err = s.upsertEntries(ctx, tx, DBID(feedDBID), feed.Items)
		if err != nil {
			return fail(err)
		}

		return nil
	}

	return s.withTx(ctx, dbFunc, nil)
}

func (s *sqliteStore) upsertEntries(
	ctx context.Context,
	tx *sql.Tx,
	feedDBID DBID,
	entries []*gofeed.Item,
) error {

	sql1 := `
	INSERT INTO
		entries(
			feed_id,
			external_id,
			url,
			title,
			description,
			content,
			publication_time,
			update_time
		)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()

	sql2 := `UPDATE entries SET is_read = (update_time = ?)`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()

	for _, entry := range entries {
		updateTime := serializeTime(resolveEntryUpdateTime(entry))
		_, err = stmt1.ExecContext(
			ctx,
			feedDBID,
			entry.GUID,
			entry.Link,
			entry.Title,
			nullIf(entry.Description, textEmpty),
			nullIf(entry.Content, textEmpty),
			nullIf(entry.Published, textEmpty),
			serializeTime(resolveEntryPublishedTime(entry)),
			updateTime,
		)
		if err != nil {
			if isUniqueErr(err, "UNIQUE constraint failed: entries.feed_id, entries.external_id") {
				if _, ierr := stmt2.ExecContext(ctx, updateTime); ierr != nil {
					return ierr
				}
			} else {
				return err
			}
		}
	}

	return nil
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

	sql3 := `INSERT OR IGNORE INTO feeds_x_feed_categories(feed_id, feed_category_id) VALUES (?, ?)`
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
