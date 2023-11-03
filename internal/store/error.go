// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"errors"
	"fmt"
	"strings"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var ErrEmptyPayload = errors.New("payload is empty")

type FeedNotFoundError struct{ FeedID ID }

func (e FeedNotFoundError) Error() string {
	return fmt.Sprintf("feed with ID=%d not found", e.FeedID)
}

type EntryNotFoundError struct{ EntryID ID }

func (e EntryNotFoundError) Error() string {
	return fmt.Sprintf("entry with ID=%d not found", e.EntryID)
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
