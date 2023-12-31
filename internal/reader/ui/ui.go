// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import "github.com/gdamore/tcell/v2"

type KeyHandler = func(*tcell.EventKey) *tcell.EventKey

//nolint:unused
type drawFunc func(screen tcell.Screen, x int, y int, w int, h int) (ix int, iy int, iw int, ih int)

const (
	leftPopupMargin      = 2
	rightPopupMargin     = 2
	verticalPopupPadding = 4
)
