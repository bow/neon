//go:build linux

// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package state

import (
	"path/filepath"

	"github.com/adrg/xdg"

	"github.com/bow/neon/internal"
)

func stateDir() (string, error) {
	sd := filepath.Join(xdg.StateHome, internal.AppName())
	return sd, nil
}
