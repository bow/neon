// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package state

import (
	"os"
	"path/filepath"
)

type FileSystemState struct {
	initPath string
}

func newFileSystemState() (*FileSystemState, error) {
	sd, err := stateDir()
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(sd)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := os.MkdirAll(sd, os.ModeDir|0o700); err != nil {
			return nil, err
		}
	}

	fst := FileSystemState{initPath: filepath.Join(sd, initFileName)}

	return &fst, nil
}

func (s *FileSystemState) MarkIntroSeen() {
	_, _ = os.Create(s.initPath)
}

func (s *FileSystemState) IntroSeen() bool {
	_, err := os.Stat(s.initPath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	// Default is to assume already seen.
	return true
}

var _ State = new(FileSystemState)

var initFileName = "reader.initialized"
