package internal

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/bow/courier/internal/migration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

type FeedStore interface {
	AddFeed(context.Context, *gofeed.Feed, *string, *string, []string) error
}

type DBID = int

type feedDB struct {
	db *sql.DB
	mu sync.RWMutex
}

func newFeedDB(filename string) (*feedDB, error) {

	log.Debug().Msgf("preparing '%s' as data store", filename)
	fail := failF("newFeedDB")

	m, err := migration.New(filename)
	if err != nil {
		return nil, fail(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fail(err)
	}
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, fail(err)
	}

	store := feedDB{db: db}

	return &store, nil
}

func (f *feedDB) withTx(
	ctx context.Context,
	dbFunc func(context.Context, *sql.Tx) error,
	txOpts *sql.TxOptions,
) (err error) {
	tx, err := f.db.BeginTx(ctx, txOpts)
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
