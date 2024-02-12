// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type readingPane struct {
	tview.Box

	theme *Theme
	lang  *Lang

	narrowBranchPoint int
}

func newReadingPane(theme *Theme, lang *Lang, narrowBranchPoint int) *readingPane {
	rp := readingPane{
		theme: theme,
		lang:  lang,

		narrowBranchPoint: narrowBranchPoint,
	}
	box := tview.NewBox()
	rp.Box = *box

	focusf, unfocusf := rp.makeDrawFuncs()
	rp.SetDrawFunc(unfocusf)
	rp.SetFocusFunc(func() { rp.SetDrawFunc(focusf) })
	rp.SetBlurFunc(func() { rp.SetDrawFunc(unfocusf) })

	return &rp
}

func (rp *readingPane) makeDrawFuncs() (focusf, unfocusf drawFunc) {

	titleUF, titleF := fmtPaneTitle(rp.lang.readingPaneTitle)

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
			style := rp.theme.lineStyle()
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				if cx == rp.narrowBranchPoint {
					screen.SetContent(cx, y, tview.BoxDrawingsLightUpAndHorizontal, nil, style)
				} else {
					screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
				}
			}
			screen.SetContent(x-1, y, tview.BoxDrawingsLightVerticalAndRight, nil, style)

			// Write the title text.
			tview.Print(
				screen,
				title,
				x+leftPad,
				y,
				width-2,
				tview.AlignLeft,
				rp.theme.titleFG,
			)

			return x + 2, y + 1, width - 2, height - 1
		}
	}

	focusf = drawf(true)
	unfocusf = drawf(false)

	return focusf, unfocusf
}
