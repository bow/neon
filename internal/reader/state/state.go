// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package state

// State describes local state that persists between runs.
type State interface {
	MarkIntroSeen()
	IntroSeen() bool
}

func NewState() State {
	st, err := newFileSystemState()
	if err != nil {
		return newNullState()
	}
	return st
}
