// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type feedsPane struct {
	*tview.TreeView

	theme   *Theme
	navRoot *tview.TreeNode
}

func newFeedsPane(theme *Theme) *feedsPane {

	fp := feedsPane{
		TreeView: tview.NewTreeView(),
		theme:    theme,
	}
	fp.setupNavTree()

	var titleUF, titleF string
	if theme.FeedsPaneTitle != "" {
		titleUF = fmt.Sprintf(" %s ", theme.FeedsPaneTitle)
		titleF = fmt.Sprintf("[::b]» %s[::-] ", theme.FeedsPaneTitle)
	} else {
		titleF = "[::b]»[::-] "
	}

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
			lineStyle := fp.theme.lineStyle()
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, lineStyle)
			}

			// Write the title text.
			tview.Print(
				screen,
				title,
				x+leftPad,
				y,
				width-2,
				tview.AlignLeft,
				fp.theme.TitleForeground,
			)

			return x + 2, y + 1, width - 2, height - 1
		}
	}

	focusf := drawf(true)
	ufocusf := drawf(false)

	fp.SetDrawFunc(ufocusf)
	fp.SetFocusFunc(func() { fp.SetDrawFunc(focusf) })
	fp.SetBlurFunc(func() { fp.SetDrawFunc(ufocusf) })

	return &fp
}

func (fp *feedsPane) setupNavTree() {

	navRoot := tview.NewTreeNode("")
	fp.navRoot = navRoot

	fp.SetRoot(navRoot).
		SetCurrentNode(navRoot).
		SetTopLevel(1)
	updateGroups := []string{"Today", "This Week", "This Month", "This Year"}
	for _, ug := range updateGroups {
		node := tview.NewTreeNode(ug).
			SetSelectable(true).
			SetColor(fp.theme.FeedsGroup)
		navRoot.AddChild(node)
	}
}
