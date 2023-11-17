// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/iris/internal"
)

// Show displays a reader for the given datastore.
func Show(db internal.FeedStore) error {

	lineForeground := tcell.ColorWhite
	titleForeground := tcell.ColorBlue
	versionForeground := tcell.ColorGray
	lastPullForeground := tcell.ColorGray
	statsForeground := tcell.ColorDarkGoldenrod
	wideViewMinWidth := 150

	newFeedsPane := func(withBottomBorder bool) *tview.Box {
		return newPane("Feeds", titleForeground, lineForeground, 1, withBottomBorder)
	}

	newEntriesPane := func() *tview.Box {
		return newPane("Entries", titleForeground, lineForeground, 1, true)
	}

	narrowReadingGrid := tview.NewGrid().
		SetRows(-1, -2).
		SetBorders(false).
		AddItem(newFeedsPane(false), 0, 0, 1, 1, 0, 0, false).
		AddItem(newEntriesPane(), 1, 0, 1, 1, 0, 0, false)

	wideReadingGrid := tview.NewGrid().
		SetColumns(45, 1, 0).
		SetBorders(false).
		AddItem(newFeedsPane(true), 0, 0, 1, 1, 0, wideViewMinWidth, false).
		AddItem(newVerticalDivider(lineForeground), 0, 1, 1, 1, 0, wideViewMinWidth, false).
		AddItem(newEntriesPane(), 0, 2, 1, 1, 0, wideViewMinWidth, false)

	stats, err := db.GetGlobalStats(context.Background())
	if err != nil {
		return err
	}

	header := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignRight), 0, 1, false)

	unreadInfo := tview.NewTextView().
		SetTextColor(statsForeground).
		SetText(fmt.Sprintf("%d unread entries", stats.NumEntriesUnread))

	lastPullInfo := tview.NewTextView().
		SetTextColor(lastPullForeground).
		SetText(
			fmt.Sprintf("Pulled %s", stats.LastPullTime.Local().Format("02/Jan/06 15:04")),
		)

	versionInfo := tview.NewTextView().
		SetTextColor(versionForeground).
		SetText(fmt.Sprintf("iris v%s", internal.Version()))

	footer := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		// TODO: Refresh values when requested.
		AddItem(lastPullInfo.SetTextAlign(tview.AlignLeft), 0, 1, false).
		AddItem(unreadInfo.SetTextAlign(tview.AlignCenter), 0, 1, false).
		AddItem(versionInfo.SetTextAlign(tview.AlignRight), 0, 1, false)

	root := tview.NewGrid().
		SetRows(1, 0, 1).
		SetBorders(false)

	// Narrow layout, less than 100px wide.
	root.
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(narrowReadingGrid, 1, 0, 1, 1, 0, 0, false).
		AddItem(footer, 2, 0, 1, 1, 0, 0, false)

	// Wide layout, width of 100px or more.
	root.
		AddItem(header, 0, 0, 1, 1, 0, wideViewMinWidth, false).
		AddItem(wideReadingGrid, 1, 0, 1, 1, 0, wideViewMinWidth, false).
		AddItem(footer, 2, 0, 1, 1, 0, wideViewMinWidth, false)

	app := tview.NewApplication()

	if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}

func newPane(
	text string,
	textForeground, lineForeground tcell.Color,
	titleLeftPad int,
	withBottomBorder bool,
) *tview.Box {

	lineStyle := tcell.StyleDefault.Foreground(lineForeground).Background(tcell.ColorBlack)

	hBorderF := func(screen tcell.Screen, cx, cy int) {
		screen.SetContent(cx, cy, tview.BoxDrawingsLightHorizontal, nil, lineStyle)
	}

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		// Draw top and optionally bottom borders.
		for cx := x; cx < x+width; cx++ {
			hBorderF(screen, cx, y)
		}
		if withBottomBorder {
			for cx := x; cx < x+width; cx++ {
				hBorderF(screen, cx, y+height-1)
			}
		}

		var displayed string
		if text != "" {
			displayed = fmt.Sprintf(" %s ", text)
		}

		// Write the title text.
		tview.Print(
			screen,
			displayed,
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
