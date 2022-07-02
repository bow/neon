package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"

	"github.com/bow/courier/internal/migration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"
)

type DBID = int

type SQLite struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewSQLite(filename string) (*SQLite, error) {

	log.Debug().Msgf("preparing '%s' as data store", filename)
	fail := failF("NewSQLiteStore")

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
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fail(err)
	}

	store := SQLite{db: db}

	return &store, nil
}

func (s *SQLite) withTx(
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

// nullIf returns nil if the given string is empty or only contains whitespaces, otherwise
// it returns a pointer to the string value.
func nullIf[T any](value T, pred func(T) bool) *T {
	if pred(value) {
		return nil
	}
	return &value
}

func nullIfTextEmpty(v string) *string {
	return nullIf(v, func(s string) bool { return s == "" || strings.TrimSpace(s) == "" })
}

// resolve returns the dereferenced pointer value if the pointer is non-nil,
// otherwise it returns the given default.
func resolve[T any](v *T, def T) T {
	if v != nil {
		return *v
	}
	return def
}

func setTableField[T any](
	tableName, columnName string,
) func(context.Context, *sql.Tx, DBID, *T) error {

	if tableName != "feeds" && tableName != "entries" {
		panic("unexpected tableName: " + tableName)
	}

	return func(ctx context.Context, tx *sql.Tx, id DBID, fieldValue *T) error {

		// nil pointers mean no value is given and so no updates are needed.
		if fieldValue == nil {
			return nil
		}

		sql1 := `UPDATE ` + tableName + ` SET ` + columnName + ` = $2 WHERE id = $1 RETURNING id`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return err
		}
		defer stmt1.Close()

		var updatedID DBID
		err = stmt1.QueryRowContext(ctx, id, fieldValue).Scan(&updatedID)
		if err != nil {
			return err
		}
		if updatedID == 0 {
			switch tableName {
			case "feeds":
				return FeedNotFoundError{id}
			case "entries":
				return EntryNotFoundError{id}
			}
			panic("unexpected tableName: " + tableName)
		}
		return nil
	}
}
