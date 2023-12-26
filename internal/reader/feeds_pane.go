// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type feedsPane struct {
	*tview.TreeView

	theme *Theme
}

func newFeedsPane(theme *Theme) *feedsPane {

	fp := feedsPane{theme: theme}
	fp.setupNavTree()

	focusf, unfocusf := fp.makeDrawFuncs()
	fp.SetDrawFunc(unfocusf)
	fp.SetFocusFunc(func() { fp.SetDrawFunc(focusf) })
	fp.SetBlurFunc(func() { fp.SetDrawFunc(unfocusf) })

	return &fp
}

func (fp *feedsPane) makeDrawFuncs() (focusf, unfocusf drawFunc) {

	var titleUF, titleF string
	if fp.theme.FeedsPaneTitle != "" {
		titleUF = fmt.Sprintf(" %s ", fp.theme.FeedsPaneTitle)
		titleF = fmt.Sprintf("[::b]▶ %s[::-] ", fp.theme.FeedsPaneTitle)
	} else {
		titleF = "[::b]▶[::-] "
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

	focusf = drawf(true)
	unfocusf = drawf(false)

	return focusf, unfocusf
}

func (fp *feedsPane) setupNavTree() {

	root := tview.NewTreeNode("")

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root).
		SetTopLevel(1)

	updateGroups := []string{"Today", "This Week", "This Month", "This Year"}

	for _, ug := range updateGroups {
		node := tview.NewTreeNode(ug).
			SetSelectable(true).
			SetColor(fp.theme.FeedsGroup)
		root.AddChild(node)
	}

	fp.TreeView = tree
}

func (fp *feedsPane) refreshColors() {
	for _, node := range fp.TreeView.GetRoot().GetChildren() {
		node.SetColor(fp.theme.FeedsGroup)
	}
}
