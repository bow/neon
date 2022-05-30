package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/bow/courier/internal/migration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

type FeedStore interface {
	AddFeed(*gofeed.Feed) error
}

type feedDB struct {
	db *sql.DB
	mu sync.RWMutex
}

func newFeedDB(filename string) (*feedDB, error) {
	log.Debug().Msgf("preparing '%s' as data store", filename)

	m, err := migration.New(filename)
	if err != nil {
		return nil, fmt.Errorf("migration setup: %w", err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("migration up: %w", err)
	}
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	store := feedDB{db: db}

	return &store, nil
}

func (f *feedDB) AddFeed(_ *gofeed.Feed) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return nil
}
