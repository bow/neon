// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type testSQLiteDB struct {
	*SQLite
	t      *testing.T
	parser *MockParser
}

func newTestSQLiteDB(t *testing.T) testSQLiteDB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), t.Name()+".db")
	prs := NewMockParser(gomock.NewController(t))
	s, err := newSQLiteWithParser(dbPath, prs)
	require.NoError(t, err)

	return testSQLiteDB{s, t, prs}
}

func (db *testSQLiteDB) tx() *sql.Tx {
	db.t.Helper()

	tx, err := db.handle.Begin()
	require.NoError(db.t, err)

	return tx
}

func (db *testSQLiteDB) countTableRows(tableName string) int {
	db.t.Helper()

	tx := db.tx()
	stmt, err := tx.Prepare(fmt.Sprintf(`SELECT count(id) FROM %s`, tableName))
	require.NoError(db.t, err)

	var count int
	row := stmt.QueryRow()
	require.NoError(db.t, row.Scan(&count))
	require.NoError(db.t, tx.Rollback())

	return count
}

func (db *testSQLiteDB) rowExists(
	query string,
	args ...any,
) bool {
	db.t.Helper()

	tx := db.tx()
	stmt, err := tx.Prepare(fmt.Sprintf("SELECT EXISTS (%s)", query))
	require.NoError(db.t, err)

	var exists bool
	row := stmt.QueryRow(args...)
	require.NoError(db.t, row.Scan(&exists))
	require.NoError(db.t, tx.Rollback())

	return exists
}

func (db *testSQLiteDB) countFeeds() int {
	return db.countTableRows("feeds")
}

func (db *testSQLiteDB) countEntries(xmlURL string) int {
	db.t.Helper()

	tx := db.tx()
	stmt, err := tx.Prepare(`
	SELECT
		count(e.id)
	FROM
		entries e
		INNER JOIN feeds f ON e.feed_id = f.id
	WHERE
		f.feed_url = ?
`,
	)
	require.NoError(db.t, err)

	var count int
	row := stmt.QueryRow(xmlURL)
	require.NoError(db.t, row.Scan(&count))
	require.NoError(db.t, tx.Rollback())

	return count
}

func (db *testSQLiteDB) countFeedTags() int {
	return db.countTableRows("feed_tags")
}

func (db *testSQLiteDB) getFeedUpdateTime(feedURL string) *time.Time {
	db.t.Helper()

	tx := db.tx()
	stmt1, err := tx.Prepare(`SELECT update_time FROM feeds WHERE feed_url = ?`)
	require.NoError(db.t, err)

	var updateTime string
	err = stmt1.QueryRow(feedURL, feedURL).Scan(&updateTime)
	require.NoError(db.t, err)

	return mustTimeP(db.t, updateTime)
}

func (db *testSQLiteDB) getFeedSubTime(feedURL string) time.Time {
	db.t.Helper()

	tx := db.tx()
	stmt1, err := tx.Prepare(`SELECT sub_time FROM feeds WHERE feed_url = ?`)
	require.NoError(db.t, err)

	var subTime time.Time
	err = stmt1.QueryRow(feedURL, feedURL).Scan(&subTime)
	require.NoError(db.t, err)

	return subTime
}

func (db *testSQLiteDB) getEntryID(feedURL string, entryExtID string) ID {
	db.t.Helper()

	tx := db.tx()
	stmt1, err := tx.Prepare(`
		SELECT
			e.id
		FROM
			entries e
			INNER JOIN feeds f ON e.feed_id = f.id
		WHERE
			f.feed_url = ?
			AND e.external_id = ?
	`)
	require.NoError(db.t, err)

	var entryID ID
	err = stmt1.QueryRow(feedURL, entryExtID).Scan(&entryID)
	require.NoError(db.t, err)

	return entryID
}

func (db *testSQLiteDB) getEntryUpdateTime(feedURL string, entryExtID string) *time.Time {
	db.t.Helper()

	tx := db.tx()
	stmt1, err := tx.Prepare(`
		SELECT
			e.update_time
		FROM
			entries e
			INNER JOIN feeds f ON e.feed_id = f.id
		WHERE
			f.feed_url = ?
			AND e.external_id = ?
	`)
	require.NoError(db.t, err)

	var entryUpdateTime sql.NullString
	err = stmt1.QueryRow(feedURL, entryExtID).Scan(&entryUpdateTime)
	require.NoError(db.t, err)

	if entryUpdateTime.Valid {
		return mustTimeP(db.t, entryUpdateTime.String)
	}
	return nil
}

func (db *testSQLiteDB) getEntryPubTime(feedURL string, entryExtID string) *time.Time {
	db.t.Helper()

	tx := db.tx()
	stmt1, err := tx.Prepare(`
		SELECT
			e.pub_time
		FROM
			entries e
			INNER JOIN feeds f ON e.feed_id = f.id
		WHERE
			f.feed_url = ?
			AND e.external_id = ?
	`)
	require.NoError(db.t, err)

	var entryPubTime sql.NullString
	err = stmt1.QueryRow(feedURL, entryExtID).Scan(&entryPubTime)
	require.NoError(db.t, err)

	if entryPubTime.Valid {
		return mustTimeP(db.t, entryPubTime.String)
	}
	return nil
}

func (db *testSQLiteDB) addFeeds(feeds []*feedRecord) map[string]feedKey {
	db.t.Helper()

	tx := db.tx()
	stmt1, err := tx.Prepare(`
		INSERT INTO
			feeds(
				title
				, feed_url
				, site_url
				, description
				, is_starred
				, sub_time
				, last_pull_time
				, update_time
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`)
	require.NoError(db.t, err)
	stmt2, err := tx.Prepare(`
		INSERT INTO entries(
			feed_id
			, external_id
			, title
			, url
			, is_read
			, is_bookmarked
			, update_time
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`)
	require.NoError(db.t, err)

	keys := make(map[string]feedKey)
	for _, feed := range feeds {
		var feedID ID
		subTime := feed.subscribed
		if subTime.IsZero() {
			subTime = time.Now().UTC()
		}
		err = stmt1.QueryRow(
			feed.title,
			feed.feedURL,
			feed.siteURL,
			feed.description,
			feed.isStarred,
			subTime,
			subTime, // last_pull_time defaults to sub_time
			feed.updated,
		).Scan(&feedID)
		require.NoError(db.t, err)

		entries := make(map[string]ID)

		for i, entry := range feed.entries {
			var (
				entryID    ID
				extID      = entry.extID
				updateTime = entry.updated
			)
			if extID == "" {
				extID = fmt.Sprintf("%s-entry-%d", db.t.Name(), i)
			}
			if updateTime.Time.IsZero() && !updateTime.Valid {
				updateTime.Time = time.Now().UTC()
				updateTime.Valid = true
			}
			err = stmt2.QueryRow(
				feedID,
				extID,
				entry.title,
				entry.url,
				entry.isRead,
				entry.isBookmarked,
				updateTime,
			).Scan(&entryID)
			require.NoError(db.t, err)
			entries[entry.title] = entryID
		}

		keys[feed.title] = feedKey{ID: feedID, Title: feed.title, Entries: entries}

		if len(feed.tags) > 0 {
			require.NoError(db.t, addFeedTags(context.Background(), tx, feedID, feed.tags))
		}
	}
	require.NoError(db.t, tx.Commit())

	return keys
}

func (db *testSQLiteDB) addFeedWithURL(url string) {
	db.t.Helper()

	tx := db.tx()
	stmt, err := tx.Prepare(`INSERT INTO feeds(title, feed_url, last_pull_time) VALUES (?, ?, ?)`)
	require.NoError(db.t, err)

	_, err = stmt.Exec(db.t.Name(), url, time.Now().UTC().Format(time.RFC3339))
	require.NoError(db.t, err)
	require.NoError(db.t, tx.Commit())
}
