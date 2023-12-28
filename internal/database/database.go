// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/database/migration"
)

type ID = uint32

type SQLite struct {
	mu     sync.RWMutex
	handle *sql.DB
	parser internal.Parser
}

func NewSQLite(filename string) (*SQLite, error) {
	return NewSQLiteWithParser(filename, gofeed.NewParser())
}

func NewSQLiteWithParser(filename string, parser internal.Parser) (*SQLite, error) {

	fail := failF("NewSQLite")

	pkgLogger.Debug().Msgf("migrating database")
	m, err := migration.New(filename)
	if err != nil {
		return nil, fail(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fail(err)
	}
	dsv, dsd, dserr := m.Version()
	if dserr != nil {
		return nil, fail(err)
	}
	sv := fmt.Sprintf("%d", dsv)
	if dsd {
		sv = fmt.Sprintf("%s*", sv)
	}

	pkgLogger.Debug().
		Str("database_schema_version", sv).
		Msg("migrated database")

	handle, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, fail(err)
	}
	_, err = handle.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fail(err)
	}

	db := SQLite{handle: handle, parser: parser}

	return &db, nil
}

func (db *SQLite) withTx(
	ctx context.Context,
	dbFunc func(context.Context, *sql.Tx) error,
) (err error) {
	tx, err := db.handle.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	rb := func(tx *sql.Tx) {
		rerr := tx.Rollback()
		if rerr != nil {
			pkgLogger.Error().Err(rerr).Msg("failed to roll back transaction")
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

type editableTable interface {
	name() string
	errNotFound(id ID) error
}

type feedsTableType struct{}

func (t *feedsTableType) name() string            { return "feeds" }
func (t *feedsTableType) errNotFound(id ID) error { return internal.FeedNotFoundError{ID: id} }

type entriesTableType struct{}

func (t *entriesTableType) name() string            { return "entries" }
func (t *entriesTableType) errNotFound(id ID) error { return internal.EntryNotFoundError{ID: id} }

var (
	feedsTable   = &feedsTableType{}
	entriesTable = &entriesTableType{}
)

func tableFieldSetter[T any](
	table editableTable,
	columnName string,
) func(context.Context, *sql.Tx, ID, *T) error {

	return func(ctx context.Context, tx *sql.Tx, id ID, fieldValue *T) error {

		// nil pointers mean no value is given and so no updates are needed.
		if fieldValue == nil {
			return nil
		}

		// https://github.com/golang/go/issues/18478
		// nolint: gosec
		sql1 := `UPDATE ` + table.name() + ` SET ` + columnName + ` = $2 WHERE id = $1 RETURNING id`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return err
		}
		defer stmt1.Close()

		var updatedID ID
		err = stmt1.QueryRowContext(ctx, id, fieldValue).Scan(&updatedID)
		if err != nil {
			return err
		}
		if updatedID == 0 {
			return table.errNotFound(id)
		}
		return nil
	}
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

// failF creates a function for wrapping other error functions.
func failF(funcName string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("%s: %w", funcName, err)
	}
}

func SetLogger(logger zerolog.Logger) {
	pkgLogger = logger
}

// pkgLogger is the server package pkgLogger.
var pkgLogger = zerolog.Nop()
