// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// feedKey is a helper struct for testing.
type feedKey struct {
	ID      ID
	Title   string
	Entries map[string]ID
}

// toNullString wraps the given string into an sql.NullString value. An empty string input is
// considered a database NULL value.
func toNullString(v string) sql.NullString {
	return sql.NullString{String: v, Valid: v != ""}
}

func toNullTime(v time.Time) sql.NullTime {
	return sql.NullTime{Time: v, Valid: !v.IsZero()}
}

func mustTime(t *testing.T, value string) time.Time {
	t.Helper()
	tv := mustTimeP(t, value)
	return *tv
}

func mustTimeP(t *testing.T, value string) *time.Time {
	t.Helper()
	tv, err := deserializeTime(value)
	require.NoError(t, err)
	return tv
}

func deserializeTime(v string) (*time.Time, error) {
	if v == "" {
		return nil, nil
	}
	pv, err := time.Parse(time.RFC3339Nano, v)
	if err != nil {
		return nil, err
	}
	upv := pv.UTC()
	return &upv, nil
}
