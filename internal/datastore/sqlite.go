// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/mmcdole/gofeed"

	"github.com/bow/neon/internal/datastore/migration"
)

type SQLite struct {
	mu     sync.RWMutex
	handle *sql.DB
	parser Parser
}

// Ensure SQLite implements Datastore.
var _ Datastore = new(SQLite)

func NewSQLite(filename string) (*SQLite, error) {
	return newSQLiteWithParser(filename, gofeed.NewParser())
}

func newSQLiteWithParser(filename string, parser Parser) (*SQLite, error) {

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
