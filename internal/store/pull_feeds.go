package store

import (
	"context"
	"database/sql"
	"sync"
)

func (s *SQLite) PullFeeds(ctx context.Context) (<-chan PullResult, error) {

	fail := failF("SQLite.PullFeeds")

	var c <-chan PullResult
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		pks, err := getAllPullKeys(ctx, tx)
		if err != nil {
			return err
		}
		if len(pks) == 0 {
			return nil
		}

		pullCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		chs := make([]<-chan PullResult, len(pks))
		for i, pk := range pks {
			chs[i] = pullNewFeedEntries(pullCtx, tx, pk, s.parser)
		}

		c = merge(chs)

		return nil
	}

	err := s.withTx(ctx, dbFunc, nil)
	if err != nil {
		return nil, fail(err)
	}

	return c, nil
}

// PullResult is a container for a pull operation.
type PullResult struct {
	pk     pullKey
	status pullStatus
	ok     *Feed
	err    error
}

func (msg PullResult) Result() *Feed {
	if msg.status == pullSuccess {
		return msg.ok
	}
	return nil
}

func (msg PullResult) Error() error {
	if msg.status == pullFail {
		return msg.err
	}
	return nil
}

type pullStatus int

const (
	pullSuccess pullStatus = iota
	pullFail
)

type pullKey struct {
	feedDBID DBID
	feedURL  string
}

func (pk pullKey) ok(feed *Feed) PullResult {
	return PullResult{pk: pk, status: pullSuccess, ok: feed, err: nil}
}

func (pk pullKey) err(e error) PullResult {
	return PullResult{pk: pk, status: pullFail, ok: nil, err: e}
}

func getAllPullKeys(ctx context.Context, tx *sql.Tx) ([]pullKey, error) {

	sql1 := `SELECT id, feed_url FROM feeds`

	scanRow := func(rows *sql.Rows) (pullKey, error) {
		var pk pullKey
		err := rows.Scan(&pk.feedDBID, &pk.feedURL)
		return pk, err
	}

	stmt1, err := tx.PrepareContext(ctx, sql1)
	if err != nil {
		return nil, err
	}

	rows, err := stmt1.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	pks := make([]pullKey, 0)
	for rows.Next() {
		pk, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		pks = append(pks, pk)
	}

	return pks, nil
}

func pullNewFeedEntries(
	ctx context.Context,
	tx *sql.Tx,
	pk pullKey,
	parser FeedParser,
) <-chan PullResult {

	pullf := func() PullResult {

		gfeed, err := parser.ParseURLWithContext(pk.feedURL, ctx)
		if err != nil {
			return pk.err(err)
		}

		if len(gfeed.Items) == 0 {
			return pk.ok(nil)
		}

		if err = upsertEntries(ctx, tx, pk.feedDBID, gfeed.Items); err != nil {
			return pk.err(err)
		}

		entries, err := getAllFeedEntries(ctx, tx, pk.feedDBID, pointer(false))
		if err != nil {
			return pk.err(err)
		}

		feed, err := getFeed(ctx, tx, pk.feedDBID)
		if err != nil {
			return pk.err(err)
		}

		feed.Entries = entries

		return pk.ok(feed)
	}

	ic := make(chan PullResult)
	go func() {
		defer close(ic)
		ic <- pullf()
	}()

	oc := make(chan PullResult)
	go func() {
		defer close(oc)
		select {
		case <-ctx.Done():
			oc <- pk.err(ctx.Err())
		case msg := <-ic:
			oc <- msg
		}
	}()

	return oc
}

func merge[T any](chs []<-chan T) <-chan T {
	var (
		wg     sync.WaitGroup
		merged = make(chan T, len(chs))
	)

	forward := func(ch <-chan T) {
		for msg := range ch {
			merged <- msg
		}
		wg.Done()
	}

	wg.Add(len(chs))
	for _, ch := range chs {
		go forward(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
