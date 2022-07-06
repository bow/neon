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

func (ts *testStore) tx() *sql.Tx {
	ts.t.Helper()

	tx, err := ts.db.Begin()
	require.NoError(ts.t, err)

	return tx
}

func (ts *testStore) countTableRows(tableName string) int {
	ts.t.Helper()

	tx := ts.tx()
	stmt, err := tx.Prepare(fmt.Sprintf(`SELECT count(id) FROM %s`, tableName))
	require.NoError(ts.t, err)

	var count int
	row := stmt.QueryRow()
	require.NoError(ts.t, row.Scan(&count))
	require.NoError(ts.t, tx.Rollback())

	return count
}

func (ts *testStore) rowExists(
	query string,
	args ...any,
) bool {
	ts.t.Helper()

	tx := ts.tx()
	stmt, err := tx.Prepare(fmt.Sprintf("SELECT EXISTS (%s)", query))
	require.NoError(ts.t, err)

	var exists bool
	row := stmt.QueryRow(args...)
	require.NoError(ts.t, row.Scan(&exists))
	require.NoError(ts.t, tx.Rollback())

	return exists
}

func (ts *testStore) countFeeds() int {
	return ts.countTableRows("feeds")
}

func (ts *testStore) countEntries(xmlURL string) int {
	ts.t.Helper()

	tx := ts.tx()
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
	require.NoError(ts.t, err)

	var count int
	row := stmt.QueryRow(xmlURL)
	require.NoError(ts.t, row.Scan(&count))
	require.NoError(ts.t, tx.Rollback())

	return count
}

func (ts *testStore) countFeedTags() int {
	return ts.countTableRows("feed_tags")
}

func (ts *testStore) addFeeds(feeds []*Feed) map[string]feedKey {
	ts.t.Helper()

	tx := ts.tx()
	stmt1, err := tx.Prepare(`
		INSERT INTO feeds(title, feed_url, is_starred, update_time) VALUES (?, ?, ?, ?)
		RETURNING id
	`)
	require.NoError(ts.t, err)
	stmt2, err := tx.Prepare(`
		INSERT INTO entries(feed_id, external_id, title, is_read) VALUES (?, ?, ?, ?)
		RETURNING id
	`)
	require.NoError(ts.t, err)

	keys := make(map[string]feedKey)
	for _, feed := range feeds {
		var feedDBID DBID
		err = stmt1.QueryRow(feed.Title, feed.FeedURL, feed.IsStarred, feed.Updated).
			Scan(&feedDBID)
		require.NoError(ts.t, err)

		entries := make(map[string]DBID)

		for i, entry := range feed.Entries {
			var entryDBID DBID
			extID := fmt.Sprintf("%s-entry-%d", ts.t.Name(), i)
			err = stmt2.QueryRow(feedDBID, extID, entry.Title, entry.IsRead).Scan(&entryDBID)
			require.NoError(ts.t, err)
			entries[entry.Title] = entryDBID
		}

		keys[feed.Title] = feedKey{DBID: feedDBID, Title: feed.Title, Entries: entries}

		if len(feed.Tags) > 0 {
			require.NoError(ts.t, addFeedTags(context.Background(), tx, feedDBID, feed.Tags))
		}
	}
	require.NoError(ts.t, tx.Commit())

	return keys
}

func (ts *testStore) addFeedWithURL(url string) {
	ts.t.Helper()

	tx := ts.tx()
	stmt, err := tx.Prepare(`INSERT INTO feeds(title, feed_url) VALUES (?, ?)`)
	require.NoError(ts.t, err)

	_, err = stmt.Exec(ts.t.Name(), url)
	require.NoError(ts.t, err)
	require.NoError(ts.t, tx.Commit())
}

func ts(t *testing.T, value string) *time.Time {
	t.Helper()
	tv, err := DeserializeTime(&value)
	require.NoError(t, err)
	return tv
}
