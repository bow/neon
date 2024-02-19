// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type KeyHandler = func(*tcell.EventKey) *tcell.EventKey

type drawFunc func(screen tcell.Screen, x int, y int, w int, h int) (ix int, iy int, iw int, ih int)

func fmtPaneTitle(title string) (unfocused, focused string) {
	if title != "" {
		unfocused = fmt.Sprintf(" %s ", title)
		focused = fmt.Sprintf("[::b] %s[::-] ● ", title)
	} else {
		focused = "[::b] ● [::-]"
	}
	return unfocused, focused
}
