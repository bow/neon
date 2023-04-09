// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"

	"github.com/bow/iris/internal/store/migration"
)

type DBID = int

// FeedStore describes the persistence layer interface.
type FeedStore interface {
	AddFeed(
		ctx context.Context,
		feedURL string,
		title *string,
		desc *string,
		tags []string,
		isStarred *bool,
	) (addedFeed *Feed, err error)
	EditFeeds(ctx context.Context, ops []*FeedEditOp) (feeds []*Feed, err error)
	ListFeeds(ctx context.Context) (feeds []*Feed, err error)
	PullFeeds(ctx context.Context) (results <-chan PullResult)
	DeleteFeeds(ctx context.Context, ids []DBID) (err error)
	ListEntries(ctx context.Context, feedID DBID) (entries []*Entry, err error)
	EditEntries(ctx context.Context, ops []*EntryEditOp) (entries []*Entry, err error)
	ExportOPML(ctx context.Context, title *string) (payload []byte, err error)
	ImportOPML(ctx context.Context, payload []byte) (processed int, imported int, err error)
}

type SQLite struct {
	db     *sql.DB
	mu     sync.RWMutex
	parser FeedParser
}

func NewSQLite(filename string) (*SQLite, error) {
	return NewSQLiteWithParser(filename, gofeed.NewParser())
}

func NewSQLiteWithParser(filename string, parser FeedParser) (*SQLite, error) {

	fail := failF("NewSQLiteStore")

	log.Debug().Msgf("migrating data store")
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
	dsvt := fmt.Sprintf("%d", dsv)
	if dsd {
		dsvt = fmt.Sprintf("%s*", dsvt)
	}

	log.Debug().
		Str("data_store_version", dsvt).
		Msg("migrated data store")

	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, fail(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fail(err)
	}

	store := SQLite{db: db, parser: parser}

	return &store, nil
}

func (s *SQLite) withTx(
	ctx context.Context,
	dbFunc func(context.Context, *sql.Tx) error,
) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
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

func pointerOrNil(v string) *string {
	if v == "" || strings.TrimSpace(v) == "" {
		return nil
	}
	return &v
}

// deref returns the dereferenced pointer value if the pointer is non-nil,
// otherwise it returns the given default.
func deref[T any](v *T, def T) T {
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
