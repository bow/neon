package internal

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testDB struct {
	t  *testing.T
	db *sql.DB
}

func newTestDB(t *testing.T, filename string) testDB {
	t.Helper()

	db, err := sql.Open("sqlite", filename)
	require.NoError(t, err)

	return testDB{t, db}
}

func (tdb *testDB) tx() *sql.Tx {
	tdb.t.Helper()

	tx, err := tdb.db.Begin()
	require.NoError(tdb.t, err)

	return tx
}

func (tdb *testDB) countTableRows(tableName string) int {
	tdb.t.Helper()

	tx := tdb.tx()
	stmt, err := tx.Prepare(fmt.Sprintf(`SELECT count(id) FROM %s`, tableName))
	require.NoError(tdb.t, err)

	var count int
	row := stmt.QueryRow()
	require.NoError(tdb.t, row.Scan(&count))
	require.NoError(tdb.t, tx.Rollback())

	return count
}

func (tdb *testDB) rowExists(
	query string,
	args ...any,
) bool {
	tdb.t.Helper()

	tx := tdb.tx()
	stmt, err := tx.Prepare(fmt.Sprintf("SELECT EXISTS (%s)", query))
	require.NoError(tdb.t, err)

	var exists bool
	row := stmt.QueryRow(args...)
	require.NoError(tdb.t, row.Scan(&exists))
	require.NoError(tdb.t, tx.Rollback())

	return exists
}

func (tdb *testDB) countFeeds() int {
	return tdb.countTableRows("feeds")
}

func (tdb *testDB) countEntries(xmlURL string) int {
	tdb.t.Helper()

	tx := tdb.tx()
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
	require.NoError(tdb.t, err)

	var count int
	row := stmt.QueryRow(xmlURL)
	require.NoError(tdb.t, row.Scan(&count))
	require.NoError(tdb.t, tx.Rollback())

	return count
}

func (tdb *testDB) countFeedCategories() int {
	return tdb.countTableRows("feed_categories")
}

func (tdb *testDB) addFeedWithURL(url string) {
	tdb.t.Helper()

	tx := tdb.tx()
	stmt, err := tx.Prepare(`INSERT INTO feeds(title, feed_url) VALUES (?, ?)`)
	require.NoError(tdb.t, err)

	_, err = stmt.Exec(tdb.t.Name(), url)
	require.NoError(tdb.t, err)
	require.NoError(tdb.t, tx.Commit())
}
