package internal

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/bow/courier/internal/migration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type FeedStore interface {
	AddFeed(context.Context, *gofeed.Feed, *string, *string, []string) error
	ListFeeds(context.Context) ([]*Feed, error)
}

type DBID = int

type sqliteStore struct {
	db *sql.DB
	mu sync.RWMutex
}

func newSQLiteStore(filename string) (*sqliteStore, error) {

	log.Debug().Msgf("preparing '%s' as data store", filename)
	fail := failF("newFeedDB")

	m, err := migration.New(filename)
	if err != nil {
		return nil, fail(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fail(err)
	}
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, fail(err)
	}

	store := sqliteStore{db: db}

	return &store, nil
}

func (s *sqliteStore) withTx(
	ctx context.Context,
	dbFunc func(context.Context, *sql.Tx) error,
	txOpts *sql.TxOptions,
) (err error) {
	tx, err := s.db.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}

	rb := func(tx *sql.Tx) {
		rerr := tx.Rollback()
		if rerr != nil {
			log.Error().Err(rerr).Msg("failed to roll back transaction")
		}
	}

	defer func() {
		if p := recover(); p != nil {
			rb(tx)
			panic(p)
		}
		if err != nil {
			rb(tx)
		} else {
			err = tx.Commit()
		}
	}()

	// Store txFunc results in err first so defer call above sees return value.
	err = dbFunc(ctx, tx)

	return err
}

// isUniqueErr returns true if the given error represents or wraps an SQLite unique constraint
// violation.
func isUniqueErr(err error, txtMatch string) bool {
	serr, matches := err.(*sqlite.Error)
	if matches {
		return serr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE &&
			(txtMatch == "" || strings.Contains(serr.Error(), txtMatch))
	}
	if ierr := errors.Unwrap(err); ierr != nil {
		return isUniqueErr(ierr, txtMatch)
	}
	return false
}

// WrapNullString wraps the given string into an sql.NullString value. An empty string input is
// considered a database NULL value.
func WrapNullString(v string) sql.NullString {
	return sql.NullString{String: v, Valid: v != ""}
}

// UnwrapNullString unwraps the given sql.NullString value into a string pointer. If the input value
// is NULL (i.e. its `Valid` field is `false`), `nil` is returned.
func UnwrapNullString(v sql.NullString) *string {
	if v.Valid {
		s := v.String
		return &s
	}
	return nil
}

// jsonArrayString is a wrapper type that implements Scan() for database-compatible
// (de)serialization.
type jsonArrayString []string

// Value implements the database valuer interface for serializing into the database.
func (arr *jsonArrayString) Value() (driver.Value, error) {
	if arr == nil {
		return nil, nil
	}
	return json.Marshal([]string(*arr))
}

// Scan implements the database scanner interface for deserialization out of the database.
func (arr *jsonArrayString) Scan(value any) error {
	var bv []byte

	switch v := value.(type) {
	case []byte:
		bv = v
	case string:
		bv = []byte(v)
	default:
		return fmt.Errorf("value of type %T can not be scanned into a string slice", v)
	}

	return json.Unmarshal(bv, arr)
}
