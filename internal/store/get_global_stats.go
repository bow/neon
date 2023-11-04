// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"database/sql"
	"errors"
)

func (s *SQLite) GetGlobalStats(ctx context.Context) (*Stats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := &Stats{}
	dbFunc := func(ctx context.Context, tx *sql.Tx) error {
		istats, err := getGlobalStats(ctx, tx)
		if err != nil {
			return err
		}
		stats = istats
		return nil
	}

	fail := failF("SQLite.GetGlobalStats")

	err := s.withTx(ctx, dbFunc)
	if err != nil {
		return nil, fail(err)
	}
	return stats, nil
}

func getGlobalStats(ctx context.Context, tx *sql.Tx) (*Stats, error) {

	var stats Stats

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

	if err = stmt1.QueryRowContext(ctx).Scan(&stats.NumFeeds, &stats.NumEntries); err != nil {
		return nil, err
	}
	if err = stmt2.QueryRowContext(ctx).Scan(&stats.NumEntriesUnread); err != nil {
		return nil, err
	}
	if stats.NumFeeds == 0 {
		return &stats, err
	}
	if err = stmt3.QueryRowContext(ctx).Scan(&stats.rawLastPullTime); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	if err = stmt4.QueryRowContext(ctx).Scan(&stats.rawMostRecentUpdateTime); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return &stats, nil
}
