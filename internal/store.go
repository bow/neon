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

func (f *feedDB) AddFeed(
	ctx context.Context,
	feed *gofeed.Feed,
	title *string,
	desc *string,
	categories []string,
) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	fail := failF("FeedStore.AddFeed")

	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		sql1 := `INSERT INTO feeds(title, description, xml_url, html_url) VALUES (?, ?, ?, ?)`
		stmt1, err := tx.PrepareContext(ctx, sql1)
		if err != nil {
			return fail(err)
		}
		defer stmt1.Close()

		res, err := stmt1.ExecContext(
			ctx,
			nullIf(resolve(title, feed.Title), textEmpty),
			nullIf(resolve(desc, feed.Description), textEmpty),
			feed.FeedLink,
			nullIf(feed.Link, textEmpty),
		)
		if err != nil {
			return fail(err)
		}

		feedDBID, err := res.LastInsertId()
		if err != nil {
			return fail(err)
		}

		err = f.addFeedCategories(ctx, tx, DBID(feedDBID), categories)
		if err != nil {
			return fail(err)
		}

		return nil
	}

	return f.withTx(ctx, dbFunc, nil)
}

func (f *feedDB) addFeedCategories(
	ctx context.Context,
	tx *sql.Tx,
	feedDBID DBID,
	cats []string,
) error {

	sql1 := `INSERT OR IGNORE INTO feed_categories(name) VALUES (?)`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return err
	}
	defer stmt1.Close()
	for _, cat := range cats {
		_, err = stmt1.ExecContext(ctx, cat)
		if err != nil {
			return err
		}
	}

	sql2 := `SELECT id FROM feed_categories WHERE name = ?`
	stmt2, err := tx.PrepareContext(ctx, sql2)
	if err != nil {
		return err
	}
	defer stmt2.Close()
	ids := make(map[string]DBID)
	for _, cat := range cats {
		if _, exists := ids[cat]; exists {
			continue
		}
		var id DBID
		row := stmt2.QueryRowContext(ctx, cat)
		if err = row.Scan(&id); err != nil {
			return err
		}
		ids[cat] = id
	}

	sql3 := `INSERT INTO feeds_x_feed_categories(feed_id, feed_category_id) VALUES (?, ?)`
	stmt3, err := tx.PrepareContext(ctx, sql3)
	if err != nil {
		return err
	}
	defer stmt3.Close()

	for _, catDBID := range ids {
		if _, err := stmt3.ExecContext(ctx, feedDBID, catDBID); err != nil {
			return err
		}
	}

	return nil
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
