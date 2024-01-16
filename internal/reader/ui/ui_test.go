// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import "time"

var (
	now          = time.Now()
	day          = 24 * time.Hour
	yesterday    = now.Add(-1 * day)
	threeDaysAgo = now.Add(-3 * day)
	twoWeeksAgo  = now.Add(-14 * day)
)
