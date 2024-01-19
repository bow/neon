// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/bow/neon/internal/entity"
)

type ID = uint32

// Datastore describes the persistence layer interface.
type Datastore interface {
	AddFeed(
		ctx context.Context,
		feedURL string,
		title *string,
		desc *string,
		tags []string,
		isStarred *bool,
	) (
		feed *entity.Feed,
		added bool,
		err error,
	)

	EditFeeds(
		ctx context.Context,
		ops []*entity.FeedEditOp,
	) (
		feeds []*entity.Feed,
		err error,
	)

	ListFeeds(
		ctx context.Context,
	) (
		feeds []*entity.Feed,
		err error,
	)

	PullFeeds(
		ctx context.Context,
		ids []entity.ID,
		isRead *bool,
	) (
		results <-chan entity.PullResult,
	)

	DeleteFeeds(
		ctx context.Context,
		ids []entity.ID,
	) (
		err error,
	)

	ListEntries(
		ctx context.Context,
		feedIDs []entity.ID,
		isBookmarked *bool,
	) (
		entries []*entity.Entry,
		err error,
	)

	EditEntries(
		ctx context.Context,
		ops []*entity.EntryEditOp,
	) (
		entries []*entity.Entry,
		err error,
	)

	GetEntry(
		ctx context.Context,
		id entity.ID,
	) (
		entry *entity.Entry,
		err error,
	)

	ExportSubscription(
		ctx context.Context,
		title *string,
	) (
		subscription *entity.Subscription,
		err error,
	)

	ImportSubscription(
		ctx context.Context,
		sub *entity.Subscription,
	) (
		processed int,
		imported int,
		err error,
	)

	GetGlobalStats(
		ctx context.Context,
	) (
		stats *entity.Stats,
		err error,
	)
}

func SetLogger(logger zerolog.Logger) {
	pkgLogger = logger
}

type editableTable interface {
	name() string
	errNotFound(id ID) error
}

type feedsTableType struct{}

func (t *feedsTableType) name() string            { return "feeds" }
func (t *feedsTableType) errNotFound(id ID) error { return entity.FeedNotFoundError{ID: id} }

type entriesTableType struct{}

func (t *entriesTableType) name() string            { return "entries" }
func (t *entriesTableType) errNotFound(id ID) error { return entity.EntryNotFoundError{ID: id} }

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

// pkgLogger is the server package pkgLogger.
var pkgLogger = zerolog.Nop()
