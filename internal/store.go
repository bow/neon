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

	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback() // nolint: errcheck

	sql1 := `INSERT INTO feeds(title, description) VALUES (?, ?)`
	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return fail(err)
	}
	defer stmt1.Close()

	res, err := stmt1.ExecContext(ctx, resolve(title, feed.Title), resolve(desc, feed.Description))
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

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return nil
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
