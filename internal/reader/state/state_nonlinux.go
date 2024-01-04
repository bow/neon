//go:build !linux

// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package state

func stateDir() (string, error) {
	cd, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cd, internal.AppName())
}
