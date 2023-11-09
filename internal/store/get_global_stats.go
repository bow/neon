// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bow/iris/internal"
)

func (s *SQLite) GetGlobalStats(ctx context.Context) (*internal.Stats, error) {

	aggr := &statsAggregateRecord{}
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		iaggr, err := getGlobalStats(ctx, tx)
		if err != nil {
			return err
		}
		aggr = iaggr
		return nil
	}

	fail := failF("SQLite.GetGlobalStats")

	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}

	return aggr.stats(), nil
}

func getGlobalStats(ctx context.Context, tx *sql.Tx) (*statsAggregateRecord, error) {

	var stats statsAggregateRecord

	stmt1, err := tx.PrepareContext(
		ctx,
		`
			SELECT
				COUNT(DISTINCT f.id) AS num_feeds,
				COUNT(DISTINCT e.id) AS num_entries
			FROM
				feeds f
				INNER JOIN entries e ON f.id = e.feed_id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer stmt1.Close()

	stmt2, err := tx.PrepareContext(
		ctx, `SELECT COUNT(DISTINCT e.id) FROM entries e WHERE NOT e.is_read`,
	)
	if err != nil {
		return nil, err
	}
	defer stmt2.Close()

	stmt3, err := tx.PrepareContext(
		ctx,
		`SELECT f.last_pull_time FROM feeds f ORDER BY f.last_pull_time DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer stmt3.Close()

	stmt4, err := tx.PrepareContext(
		ctx,
		`
			SELECT
				f.update_time
			FROM
				feeds f
			WHERE
				f.update_time IS NOT NULL
			ORDER BY
				f.update_time DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer stmt4.Close()

	if err = stmt1.QueryRowContext(ctx).Scan(&stats.numFeeds, &stats.numEntries); err != nil {
		return nil, err
	}
	if err = stmt2.QueryRowContext(ctx).Scan(&stats.numEntriesUnread); err != nil {
		return nil, err
	}
	if stats.numFeeds == 0 {
		return &stats, err
	}
	if err = stmt3.QueryRowContext(ctx).Scan(&stats.lastPullTime); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	if err = stmt4.QueryRowContext(ctx).Scan(&stats.mostRecentUpdateTime); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return &stats, nil
}
