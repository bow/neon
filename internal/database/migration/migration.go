// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

// package migration contains SQL migration schemas.
package migration

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	// blank import required by migrate.
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	// blank import required by migrate.
	_ "modernc.org/sqlite"
)

//go:embed *.sql
var fs embed.FS

// New creates a new migration instance for the given SQLite3 database.
func New(filename string) (*migrate.Migrate, error) {
	d, err := iofs.New(fs, ".")
	if err != nil {
		return nil, fmt.Errorf("iofs: %w", err)
	}

	p, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("abs: %w", err)
	}

	source := fmt.Sprintf("sqlite://%s", p)
	m, err := migrate.NewWithSourceInstance("iofs", d, source)
	if err != nil {
		return nil, fmt.Errorf("migrate source: %w", err)
	}
	return m, nil
}
