package store

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// feedKey is a helper struct for testing.
type feedKey struct {
	DBID    DBID
	Title   string
	Entries map[string]DBID
}

type testStore struct {
	*SQLite
	t      *testing.T
	parser *MockFeedParser
}

func newTestStore(t *testing.T) testStore {
	t.Helper()

	// TODO: Avoid global states like this.
	zerolog.SetGlobalLevel(zerolog.Disabled)

	dbPath := filepath.Join(t.TempDir(), t.Name()+".db")
	prs := NewMockFeedParser(gomock.NewController(t))
	s, err := NewSQLite(dbPath, prs)
	require.NoError(t, err)

	return testStore{s, t, prs}
}

func (s *testStore) tx() *sql.Tx {
	s.t.Helper()

	tx, err := s.db.Begin()
	require.NoError(s.t, err)

	return tx
}

func (s *testStore) countTableRows(tableName string) int {
	s.t.Helper()

	tx := s.tx()
	stmt, err := tx.Prepare(fmt.Sprintf(`SELECT count(id) FROM %s`, tableName))
	require.NoError(s.t, err)

	var count int
	row := stmt.QueryRow()
	require.NoError(s.t, row.Scan(&count))
	require.NoError(s.t, tx.Rollback())

	return count
}

func (s *testStore) rowExists(
	query string,
	args ...any,
) bool {
	s.t.Helper()

	tx := s.tx()
	stmt, err := tx.Prepare(fmt.Sprintf("SELECT EXISTS (%s)", query))
	require.NoError(s.t, err)

	var exists bool
	row := stmt.QueryRow(args...)
	require.NoError(s.t, row.Scan(&exists))
	require.NoError(s.t, tx.Rollback())

	return exists
}

func (s *testStore) countFeeds() int {
	return s.countTableRows("feeds")
}

func (s *testStore) countEntries(xmlURL string) int {
	s.t.Helper()

	tx := s.tx()
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
	require.NoError(s.t, err)

	var count int
	row := stmt.QueryRow(xmlURL)
	require.NoError(s.t, row.Scan(&count))
	require.NoError(s.t, tx.Rollback())

	return count
}

func (s *testStore) countFeedTags() int {
	return s.countTableRows("feed_tags")
}

func (s *testStore) getFeedUpdateTime(feedURL string) sql.NullString {
	s.t.Helper()

	tx := s.tx()
	stmt1, err := tx.Prepare(`SELECT update_time FROM feeds WHERE feed_url = ?`)
	require.NoError(s.t, err)

	var updateTime string
	err = stmt1.QueryRow(feedURL, feedURL).Scan(&updateTime)
	require.NoError(s.t, err)

	return sql.NullString{String: updateTime, Valid: true}
}

func (s *testStore) getFeedSubscriptionTime(feedURL string) string {
	s.t.Helper()

	tx := s.tx()
	stmt1, err := tx.Prepare(`SELECT subscription_time FROM feeds WHERE feed_url = ?`)
	require.NoError(s.t, err)

	var subscriptionTime string
	err = stmt1.QueryRow(feedURL, feedURL).Scan(&subscriptionTime)
	require.NoError(s.t, err)

	return subscriptionTime
}

func (s *testStore) getEntryDBID(feedURL string, entryExtID string) DBID {
	s.t.Helper()

	tx := s.tx()
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
	require.NoError(s.t, err)

	var entryDBID DBID
	err = stmt1.QueryRow(feedURL, entryExtID).Scan(&entryDBID)
	require.NoError(s.t, err)

	return entryDBID
}

func (s *testStore) getEntryUpdateTime(feedURL string, entryExtID string) sql.NullString {
	s.t.Helper()

	tx := s.tx()
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
	require.NoError(s.t, err)

	var entryUpdateTime sql.NullString
	err = stmt1.QueryRow(feedURL, entryExtID).Scan(&entryUpdateTime)
	require.NoError(s.t, err)

	return entryUpdateTime
}

func (s *testStore) getEntryPublicationTime(feedURL string, entryExtID string) sql.NullString {
	s.t.Helper()

	tx := s.tx()
	stmt1, err := tx.Prepare(`
		SELECT
			e.publication_time
		FROM
			entries e
			INNER JOIN feeds f ON e.feed_id = f.id
		WHERE
			f.feed_url = ?
			AND e.external_id = ?
	`)
	require.NoError(s.t, err)

	var entryPublicationTime sql.NullString
	err = stmt1.QueryRow(feedURL, entryExtID).Scan(&entryPublicationTime)
	require.NoError(s.t, err)

	return entryPublicationTime
}

func (s *testStore) addFeeds(feeds []*Feed) map[string]feedKey {
	s.t.Helper()

	tx := s.tx()
	stmt1, err := tx.Prepare(`
		INSERT INTO
			feeds(title, feed_url, site_url, description, is_starred, update_time)
			VALUES (?, ?, ?, ?, ?, ?)
		RETURNING
			id
	`)
	require.NoError(s.t, err)
	stmt2, err := tx.Prepare(`
		INSERT INTO entries(
			feed_id,
			external_id,
			title,
			url,
			is_read,
			update_time
		)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`)
	require.NoError(s.t, err)

	keys := make(map[string]feedKey)
	for _, feed := range feeds {
		var feedDBID DBID
		err = stmt1.QueryRow(
			feed.Title,
			feed.FeedURL,
			feed.SiteURL,
			feed.Description,
			feed.IsStarred,
			feed.Updated).
			Scan(&feedDBID)
		require.NoError(s.t, err)

		entries := make(map[string]DBID)

		for i, entry := range feed.Entries {
			var (
				entryDBID  DBID
				extID      = entry.ExtID
				updateTime = entry.Updated
			)
			if extID == "" {
				extID = fmt.Sprintf("%s-entry-%d", s.t.Name(), i)
			}
			if updateTime.String == "" && !updateTime.Valid {
				updateTime.String = time.Now().UTC().Format(time.RFC822)
				updateTime.Valid = true
			}
			err = stmt2.QueryRow(
				feedDBID,
				extID,
				entry.Title,
				entry.URL,
				entry.IsRead,
				updateTime,
			).Scan(&entryDBID)
			require.NoError(s.t, err)
			entries[entry.Title] = entryDBID
		}

		keys[feed.Title] = feedKey{DBID: feedDBID, Title: feed.Title, Entries: entries}

		if len(feed.Tags) > 0 {
			require.NoError(s.t, addFeedTags(context.Background(), tx, feedDBID, feed.Tags))
		}
	}
	require.NoError(s.t, tx.Commit())

	return keys
}

func (s *testStore) addFeedWithURL(url string) {
	s.t.Helper()

	tx := s.tx()
	stmt, err := tx.Prepare(`INSERT INTO feeds(title, feed_url) VALUES (?, ?)`)
	require.NoError(s.t, err)

	_, err = stmt.Exec(s.t.Name(), url)
	require.NoError(s.t, err)
	require.NoError(s.t, tx.Commit())
}

func ts(t *testing.T, value string) *time.Time {
	t.Helper()
	tv, err := DeserializeTime(&value)
	require.NoError(t, err)
	return tv
}
