package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/mmcdole/gofeed"
)

// AddFeed adds the given feed into the database.
func (s *SQLite) AddFeed(
	ctx context.Context,
	feed *gofeed.Feed,
	title *string,
	desc *string,
	tags []string,
	isStarred bool,
) (*Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("SQLite.AddFeed")

	var created *Feed
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		now := time.Now()

		feedDBID, err := insertFeedRow(ctx, tx, feed, title, desc, isStarred, &now)
		if err != nil {
			return fail(err)
		}

		err = addFeedTags(ctx, tx, feedDBID, tags)
		if err != nil {
			return fail(err)
		}

		created, err = getFeed(ctx, tx, feedDBID)
		if err != nil {
			return err
		}

		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func insertFeedRow(
	ctx context.Context,
	tx *sql.Tx,
	feed *gofeed.Feed,
	title *string,
	desc *string,
	isStarred bool,
	subTime *time.Time,
) (DBID, error) {

	var feedDBID DBID
	sql1 := `
		INSERT INTO
			feeds(
				title,
				description,
				feed_url,
				site_url,
				is_starred,
				update_time,
				subscription_time
			)
			VALUES (?, ?, ?, ?, ?, ?, ?)
`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return feedDBID, err
	}
	defer stmt1.Close()

	res, err := stmt1.ExecContext(
		ctx,
		nullIfTextEmpty(resolve(title, feed.Title)),
		nullIfTextEmpty(resolve(desc, feed.Description)),
		feed.FeedLink,
		nullIfTextEmpty(feed.Link),
		isStarred,
		serializeTime(resolveFeedUpdateTime(feed)),
		serializeTime(subTime),
	)

	if err == nil {
		lid, ierr := res.LastInsertId()
		feedDBID = DBID(lid)
		if ierr != nil {
			return feedDBID, ierr
		}
	} else {
		if !isUniqueErr(err, "UNIQUE constraint failed: feeds.feed_url") {
			return feedDBID, err
		}
		if ierr := tx.QueryRowContext(
			ctx,
			`SELECT id FROM feeds WHERE feed_url = ?`,
			feed.FeedLink,
		).Scan(&feedDBID); ierr != nil {
			return feedDBID, ierr
		}
	}
	// TODO: Add and combine with proper update call.
	if ierr := upsertEntries(ctx, tx, feedDBID, feed.Items); ierr != nil {
		return feedDBID, ierr
	}

	return feedDBID, nil
}

func upsertEntries(
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
	sql2 := `UPDATE entries SET is_read = (update_time = ?)`

	upsert := func(entry *gofeed.Item, insertStmt, updateStmt *sql.Stmt) error {
		updateTime := serializeTime(resolveEntryUpdateTime(entry))
		_, err := insertStmt.ExecContext(
			ctx,
			feedDBID,
			entry.GUID,
			entry.Link,
			entry.Title,
			nullIfTextEmpty(entry.Description),
			nullIfTextEmpty(entry.Content),
			serializeTime(resolveEntryPublishedTime(entry)),
			updateTime,
		)
		if err != nil {
			if !isUniqueErr(err, "UNIQUE constraint failed: entries.feed_id, entries.external_id") {
				return err
			}
			if _, ierr := updateStmt.ExecContext(ctx, updateTime); ierr != nil {
				return ierr
			}
		}
		return nil
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()

	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()

	for _, entry := range entries {
		if err := upsert(entry, stmt1, stmt2); err != nil {
			return err
		}
	}

	return nil
}

func addFeedTags(
	ctx context.Context,
	tx *sql.Tx,
	feedDBID DBID,
	tags []string,
) error {

	sql1 := `INSERT OR IGNORE INTO feed_tags(name) VALUES (?)`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()
	for _, tag := range tags {
		_, err = stmt1.ExecContext(ctx, tag)
		if err != nil {
			return err
		}
	}

	sql2 := `SELECT id FROM feed_tags WHERE name = ?`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()
	ids := make(map[string]DBID)
	for _, tag := range tags {
		if _, exists := ids[tag]; exists {
			continue
		}
		var id DBID
		row := stmt2.QueryRowContext(ctx, tag)
		if err = row.Scan(&id); err != nil {
			return err
		}
		ids[tag] = id
	}

	sql3 := `INSERT OR IGNORE INTO feeds_x_feed_tags(feed_id, feed_tag_id) VALUES (?, ?)`
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
