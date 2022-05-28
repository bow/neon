package internal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bow/courier/migration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"
)

// TODO: Define actual methods.
type FeedsStore interface {
	Foo() string
}

type FeedsDB struct {
	db *sql.DB
}

func NewFeedsDB(filename string) (*FeedsDB, error) {
	log.Debug().Msgf("preparing '%s' as data store", filename)

	m, err := migration.New(filename)
	if err != nil {
		return nil, fmt.Errorf("migration setup: %w", err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("migration up: %w", err)
	}
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	store := FeedsDB{db}

	return &store, nil
}

func (f *FeedsDB) Foo() string {
	return "ok!"
}
