//go:build !linux

// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bow/neon/internal"
)

// FIXME: Define this for non-linux.
var defaultDBPath = ""

// FIXME: Define this for non-linux.
func resolveDBPath(path string) (string, error) {
	return "", fmt.Errorf("not yet supported")
}

func stateDir() (string, error) {
	cd, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cd, internal.AppName())
}
