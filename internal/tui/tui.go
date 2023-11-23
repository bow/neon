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
	helpTitleForeground := tcell.ColorAqua
	helpBorderLineForeground := tcell.ColorGray

	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical
	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight

	wideViewMinWidth := 150

	root := tview.NewPages()

	feedsPane := newPane("Feeds", titleForeground, lineForeground, 1)
	entriesPane := newPane("Entries", titleForeground, lineForeground, 1)

	narrowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(feedsPane, 0, 1, false).
		AddItem(entriesPane, 0, 2, false).
		AddItem(newNarrowFooterBorder(lineForeground), 1, 0, false)

	wideFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(feedsPane, 45, 0, false).
				AddItem(newPaneDivider(lineForeground), 1, 0, false).
				AddItem(entriesPane, 0, 1, false),
			0, 1, false,
		).
		AddItem(newWideFooterBorder(lineForeground, 45), 1, 0, false)

	stats, err := db.GetGlobalStats(context.Background())
	if err != nil {
		return err
	}

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

	mainPage := tview.NewGrid().
		SetColumns(0).
		SetRows(0, 1).
		SetBorders(false).
		AddItem(footer, 1, 0, 1, 1, 0, 0, false)

	// Narrow layout.
	mainPage.
		AddItem(narrowFlex, 0, 0, 1, 1, 0, 0, false)

	// Wide layout.
	mainPage.
		AddItem(wideFlex, 0, 0, 1, 1, 0, wideViewMinWidth, false)

	helpPage := tview.NewFrame(nil).
		SetBorder(true).
		SetBorderColor(helpBorderLineForeground).
		SetTitle(" [::b]Help[::-] ").
		SetTitleColor(helpTitleForeground)

	helpPage.
		SetFocusFunc(func() { helpPage.SetTitle(" [::b]Help[::-] ") }).
		SetBlurFunc(func() { helpPage.SetTitle(" Help ") })

	root.
		AddAndSwitchToPage("main", mainPage, true).
		AddPage(
			"help",
			tview.NewGrid().
				SetColumns(0, 64, 0).
				SetRows(0, 22, 0).
				AddItem(helpPage, 1, 1, 1, 1, 0, 0, true),
			true,
			false,
		)

	app := tview.NewApplication()
	app.
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {
				if event.Rune() == 'h' {
					if fp, _ := root.GetFrontPage(); fp == "help" {
						root.HidePage("help")
					} else {
						root.ShowPage("help")
					}
					return nil
				} else if event.Rune() == 'q' {
					app.Stop()
				}
				return event
			},
		)

	if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}

func newPane(
	text string,
	textForeground, lineForeground tcell.Color,
	titleLeftPad int,
) *tview.Box {

	lineStyle := tcell.StyleDefault.Foreground(lineForeground).Background(tcell.ColorBlack)

	var unfocused, focused string
	if text != "" {
		unfocused = fmt.Sprintf(" %s ", text)
		focused = fmt.Sprintf(" [::b]%s[::-] ", text)
	}

	makedrawf := func(
		title string,
	) func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, lineStyle)
			}

			// Write the title text.
			tview.Print(
				screen,
				title,
				x+titleLeftPad,
				y,
				width-2,
				tview.AlignLeft,
				textForeground,
			)

			return x + 1, y + 1, width - 2, height - 1
		}
	}

	box := tview.NewBox().SetDrawFunc(makedrawf(unfocused))

	box.SetFocusFunc(func() { box.SetDrawFunc(makedrawf(focused)) })
	box.SetBlurFunc(func() { box.SetDrawFunc(makedrawf(unfocused)) })

	return box
}

func newPaneDivider(lineForeground tcell.Color) *tview.Box {

	style := tcell.StyleDefault.Foreground(lineForeground).Background(tcell.ColorBlack)

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		screen.SetContent(x, y, tview.BoxDrawingsLightDownAndHorizontal, nil, style)

		for cy := y + 1; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, style)
		}

		return x + 1, y + 1, width - 2, height - 1
	}

	divider := tview.NewBox().SetBorder(false).SetDrawFunc(drawf)

	return divider
}

func newNarrowFooterBorder(lineForeground tcell.Color) *tview.Box {

	style := tcell.StyleDefault.Foreground(lineForeground).Background(tcell.ColorBlack)

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
		}

		return x + 1, y + 1, width - 2, height - 1
	}

	divider := tview.NewBox().SetBorder(false).SetDrawFunc(drawf)

	return divider
}

func newWideFooterBorder(lineForeground tcell.Color, branch int) *tview.Box {

	style := tcell.StyleDefault.Foreground(lineForeground).Background(tcell.ColorBlack)

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		for cx := x; cx < x+width; cx++ {
			if cx == branch {
				screen.SetContent(cx, y, tview.BoxDrawingsLightUpAndHorizontal, nil, style)
			} else {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
			}
		}

		return x + 1, y + 1, width - 2, height - 1
	}

	divider := tview.NewBox().SetBorder(false).SetDrawFunc(drawf)

	return divider
}
