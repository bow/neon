package internal

import (
	"database/sql"
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

func (tdb *testDB) countFeeds() int {
	tdb.t.Helper()

	tx := tdb.tx()
	stmt, err := tx.Prepare(`SELECT count(id) FROM feeds`)
	require.NoError(tdb.t, err)

	var count int
	row := stmt.QueryRow()
	require.NoError(tdb.t, row.Scan(&count))
	require.NoError(tdb.t, tx.Rollback())

	return count
}
