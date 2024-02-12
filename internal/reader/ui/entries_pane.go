// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type entriesPane struct {
	tview.Box

	theme *Theme
	lang  *Lang
}

func newEntriesPane(theme *Theme, lang *Lang) *entriesPane {
	ep := entriesPane{theme: theme, lang: lang}
	box := tview.NewBox()
	ep.Box = *box

	focusf, unfocusf := ep.makeDrawFuncs()
	ep.SetDrawFunc(unfocusf)
	ep.SetFocusFunc(func() { ep.SetDrawFunc(focusf) })
	ep.SetBlurFunc(func() { ep.SetDrawFunc(unfocusf) })

	return &ep
}

func (ep *entriesPane) makeDrawFuncs() (focusf, unfocusf drawFunc) {

	titleUF, titleF := fmtPaneTitle(ep.lang.entriesPaneTitle)

	drawf := func(
		focused bool,
	) func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		var (
			title   string
			leftPad int
		)
		if focused {
			title = titleF
			leftPad = 0
		} else {
			title = titleUF
			leftPad = 1
		}

		return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			style := ep.theme.lineStyle()
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
			}

			// Write the title text.
			tview.Print(
				screen,
				title,
				x+leftPad,
				y,
				width-2,
				tview.AlignLeft,
				ep.theme.titleFG,
			)

			return x + 2, y + 1, width - 2, height - 1
		}
	}

	focusf = drawf(true)
	unfocusf = drawf(false)

	return focusf, unfocusf
}
