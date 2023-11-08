// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"

	"github.com/bow/iris/internal"
	"github.com/bow/iris/internal/store/migration"
)

type ID = uint32

// FeedStore describes the persistence layer interface.
type FeedStore interface {
	AddFeed(
		ctx context.Context,
		feedURL string,
		title *string,
		desc *string,
		tags []string,
		isStarred *bool,
	) (
		addedFeed *internal.Feed,
		err error,
	)

	EditFeeds(
		ctx context.Context,
		ops []*internal.FeedEditOp,
	) (
		feeds []*internal.Feed,
		err error,
	)

	ListFeeds(
		ctx context.Context,
	) (
		feeds []*internal.Feed,
		err error,
	)

	PullFeeds(
		ctx context.Context,
		ids []internal.ID,
	) (
		results <-chan internal.PullResult,
	)

	DeleteFeeds(
		ctx context.Context,
		ids []internal.ID,
	) (
		err error,
	)

	ListEntries(
		ctx context.Context,
		feedID internal.ID,
	) (
		entries []*internal.Entry,
		err error,
	)

	EditEntries(
		ctx context.Context,
		ops []*internal.EntryEditOp,
	) (
		entries []*internal.Entry,
		err error,
	)

	GetEntry(
		ctx context.Context,
		id internal.ID,
	) (
		entry *internal.Entry,
		err error,
	)

	// TODO: Export OPML structs instead.
	ExportOPML(
		ctx context.Context,
		exportTitle *string,
	) (
		payload []byte,
		err error,
	)

	// TODO: Import OPML structs instead.
	ImportOPML(
		ctx context.Context,
		payload []byte,
	) (
		processed int,
		imported int,
		err error,
	)

	GetGlobalStats(
		ctx context.Context,
	) (
		stats *internal.Stats,
		err error,
	)
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

type editableTable interface {
	name() string
	errNotFound(id ID) error
}

type feedsTableType struct{}

func (t *feedsTableType) name() string            { return "feeds" }
func (t *feedsTableType) errNotFound(id ID) error { return FeedNotFoundError{id} }

type entriesTableType struct{}

func (t *entriesTableType) name() string            { return "entries" }
func (t *entriesTableType) errNotFound(id ID) error { return EntryNotFoundError{id} }

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
