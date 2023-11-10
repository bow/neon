// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/iris/internal"
)

// Show displays a reader for the given datastore.
func Show(_ internal.FeedStore) error {

	lineForeground := tcell.ColorWhite
	titleForeground := tcell.ColorBlue

	feedsPane := newPane("Feeds", titleForeground, lineForeground, 1)
	entriesPane := newPane("Entries", titleForeground, lineForeground, 1)

	readingSection := tview.NewGrid().
		SetColumns(30, 1, 0).
		SetBorders(false).
		AddItem(feedsPane, 0, 0, 1, 1, 0, 0, false).
		AddItem(newVerticalDivider(lineForeground), 0, 1, 1, 1, 0, 0, false).
		AddItem(entriesPane, 0, 2, 1, 1, 0, 0, false)

	root := tview.NewGrid().
		SetRows(2, 0, 1).
		SetBorders(false).
		AddItem(newPlaceholderSection("<header>"), 0, 0, 1, 2, 0, 0, false).
		AddItem(readingSection, 1, 0, 1, 2, 0, 0, false).
		AddItem(newPlaceholderSection("<footer>"), 2, 0, 1, 2, 0, 0, false)

	app := tview.NewApplication()

	if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}

func newPlaceholderSection(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(fmt.Sprintf(" %s", text))
}

func newPane(
	text string,
	textForeground, lineForeground tcell.Color,
	titleLeftPad int,
) *tview.Box {

	lineStyle := tcell.StyleDefault.Foreground(lineForeground).Background(tcell.ColorBlack)

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		// Draw a top and bottom borders.
		ys := []int{y, y + height - 1}
		for _, cy := range ys {
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, cy, tview.BoxDrawingsLightHorizontal, nil, lineStyle)
			}
		}

		// Write the title text.
		tview.Print(
			screen,
			fmt.Sprintf(" %s ", text),
			x+titleLeftPad,
			y,
			width-2,
			tview.AlignLeft,
			textForeground,
		)

		return x + 1, y + 1, width - 2, height - 1
	}

	box := tview.NewBox().SetDrawFunc(drawf)

	return box
}

func newVerticalDivider(lineForeground tcell.Color) *tview.Box {

	style := tcell.StyleDefault.Foreground(lineForeground).Background(tcell.ColorBlack)

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		screen.SetContent(x, y, tview.BoxDrawingsLightDownAndHorizontal, nil, style)

		for cy := y + 1; cy < y+height-1; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, style)
		}

		screen.SetContent(x, y+height-1, tview.BoxDrawingsLightUpAndHorizontal, nil, style)

		return x + 1, y + 1, width - 2, height - 1
	}

	divider := tview.NewBox().SetBorder(false).SetDrawFunc(drawf)

	return divider
}
