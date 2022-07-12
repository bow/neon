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
	feedURL string,
	title *string,
	desc *string,
	tags []string,
	isStarred *bool,
) (*Feed, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fail := failF("SQLite.AddFeed")

	feed, err := s.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fail(err)
	}
	// Handle possible specs deviations.
	if feed.FeedLink == "" {
		feed.FeedLink = feedURL
	}

	var created *Feed
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {

		now := time.Now()

		feedDBID, ierr := upsertFeed(
			ctx,
			tx,
			feed.FeedLink,
			pointerOrNil(deref(title, feed.Title)),
			pointerOrNil(deref(desc, feed.Description)),
			pointerOrNil(feed.Link),
			isStarred,
			resolveFeedUpdateTime(feed),
			&now,
		)
		if ierr != nil {
			return ierr
		}

		if _, ierr = upsertEntries(ctx, tx, feedDBID, feed.Items); ierr != nil {
			return ierr
		}

		if ierr = addFeedTags(ctx, tx, feedDBID, tags); ierr != nil {
			return ierr
		}

		if created, ierr = getFeed(ctx, tx, feedDBID); ierr != nil {
			return ierr
		}

		return nil
	}

	err = s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, fail(err)
	}
	return created, nil
}

func upsertFeed(
	ctx context.Context,
	tx *sql.Tx,
	feedURL string,
	title *string,
	desc *string,
	siteURL *string,
	isStarred *bool,
	updateTime *time.Time,
	subTime *time.Time,
) (DBID, error) {

	var feedDBID DBID
	sql1 := `
		INSERT INTO
			feeds(
				feed_url,
				title,
				description,
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
		feedURL,
		title,
		desc,
		siteURL,
		deref(isStarred, false),
		serializeTime(updateTime),
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
		var ierr error
		if feedDBID, ierr = editFeedWithFeedURL(
			ctx,
			tx,
			feedURL,
			title,
			desc,
			siteURL,
			isStarred,
		); ierr != nil {
			return feedDBID, ierr
		}
	}

	return feedDBID, nil
}

func editFeedWithFeedURL(
	ctx context.Context,
	tx *sql.Tx,
	feedURL string,
	title *string,
	desc *string,
	siteURL *string,
	isStarred *bool,
) (DBID, error) {

	var feedDBID DBID

	sql1 := `SELECT id FROM feeds WHERE feed_url = ?`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return 0, err
	}

	if err := stmt1.QueryRowContext(ctx, feedURL).Scan(&feedDBID); err != nil {
		return 0, err
	}
	if err := setFeedTitle(ctx, tx, feedDBID, title); err != nil {
		return 0, err
	}
	if err := setFeedDescription(ctx, tx, feedDBID, desc); err != nil {
		return 0, err
	}
	if err := setFeedIsStarred(ctx, tx, feedDBID, isStarred); err != nil {
		return 0, err
	}
	if err := setFeedSiteURL(ctx, tx, feedDBID, siteURL); err != nil {
		return 0, err
	}
	return feedDBID, nil
}

func upsertEntries(
	ctx context.Context,
	tx *sql.Tx,
	feedDBID DBID,
	entries []*gofeed.Item,
) ([]DBID, error) {

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
		return nil, err
	}
	defer stmt1.Close()

	sql2 := `UPDATE entries SET is_read = (update_time = ?)`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return nil, err
	}
	defer stmt2.Close()

	sql3 := `SELECT id, is_read FROM entries WHERE feed_id = ? AND external_id = ?`
	stmt3, err := tx.PrepareContext(ctx, sql3)
	if err != nil {
		return nil, err
	}
	defer stmt3.Close()

	upsert := func(entry *gofeed.Item, insertStmt, updateStmt, getStmt *sql.Stmt) (DBID, error) {
		updateTime := serializeTime(resolveEntryUpdateTime(entry))
		_, err := insertStmt.Exec(
			feedDBID,
			entry.GUID,
			entry.Link,
			entry.Title,
			pointerOrNil(entry.Description),
			pointerOrNil(entry.Content),
			serializeTime(resolveEntryPublishedTime(entry)),
			updateTime,
		)
		if err != nil {
			if !isUniqueErr(err, "UNIQUE constraint failed: entries.feed_id, entries.external_id") {
				return 0, err
			}
			if _, ierr := updateStmt.Exec(ctx, updateTime); ierr != nil {
				return 0, ierr
			}
		}
		var (
			isRead    bool
			entryDBID DBID
		)
		if err = getStmt.QueryRow(feedDBID, entry.GUID).Scan(&entryDBID, &isRead); err != nil {
			return 0, err
		}
		if isRead {
			return 0, nil
		}
		return entryDBID, nil
	}

	ids := make([]DBID, 0)
	for _, entry := range entries {
		entryDBID, err := upsert(entry, stmt1, stmt2, stmt3)
		if err != nil {
			return nil, err
		}
		ids = append(ids, entryDBID)
	}

	return ids, nil
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
