// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package state

type NullState struct{}

func newNullState() *NullState { return &NullState{} }

func (s *NullState) MarkIntroSeen() {}

func (s *NullState) IntroSeen() bool { return true }

var _ State = new(NullState)
