package store

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type testStore struct {
	*SQLite
	t *testing.T
}

func newTestStore(t *testing.T) testStore {
	t.Helper()

	// TODO: Avoid global states like this.
	zerolog.SetGlobalLevel(zerolog.Disabled)

	dbPath := filepath.Join(t.TempDir(), t.Name()+".db")
	s, err := NewSQLite(dbPath)
	require.NoError(t, err)

	return testStore{s, t}
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

func (ts *testStore) countFeedCategories() int {
	return ts.countTableRows("feed_categories")
}

func (ts *testStore) addFeeds(feeds []*Feed) {
	ts.t.Helper()

	tx := ts.tx()
	stmt, err := tx.Prepare(`INSERT INTO feeds(title, feed_url, update_time) VALUES (?, ?, ?)`)
	require.NoError(ts.t, err)

	for _, feed := range feeds {
		_, err = stmt.Exec(ts.t.Name(), feed.FeedURL, feed.Updated)
		require.NoError(ts.t, err)
	}
	require.NoError(ts.t, tx.Commit())
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
