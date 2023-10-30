// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package tui

import (
	"fmt"

	"github.com/bow/iris/internal/store"
)

// Show displays a reader for the given datastore.
func Show(_ store.FeedStore) error {
	return fmt.Errorf("not implemented")
}
