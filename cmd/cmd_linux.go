//go:build linux

// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"

	"github.com/bow/neon/internal"
)

var defaultDBPath = "$XDG_DATA_HOME/neon/neon.db"

func resolveDBPath(path string) (string, error) {
	var (
		err    error
		xdgDir = "$XDG_DATA_HOME/"
	)

	if strings.HasPrefix(path, xdgDir) {
		rel := strings.TrimPrefix(path, xdgDir)
		path, err = xdg.DataFile(rel)
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func stateDir() (string, error) {
	return filepath.Join(xdg.StateHome, internal.AppName()), nil
}
