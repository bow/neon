//go:build !linux

// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"os"
	"path/filepath"

	"github.com/bow/lens/internal"
)

func stateDir() (string, error) {
	cd, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cd, internal.AppName())
}
