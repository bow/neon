// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetGlobalStatsOk(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	st := newTestStore(t)

	stats, err := st.GetGlobalStats(context.Background())
	r.Nil(stats)
	r.EqualError(err, "SQLite.GetGlobalStats: unimplemented")
}
