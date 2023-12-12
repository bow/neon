//go:build linux

// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/bow/lens/internal"
)

func stateDir() (string, error) {
	return filepath.Join(xdg.StateHome, internal.AppName()), nil
}
