// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/term"

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

	topLeftBorderTip := tview.BoxDrawingsLightVerticalAndRight
	feedsPane := newPane("Feeds", titleForeground, lineForeground, nil)
	entriesPane := newPane("Entries", titleForeground, lineForeground, nil)
	contentPane := newPane("", titleForeground, lineForeground, &topLeftBorderTip)

	narrowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(feedsPane, 0, 3, false).
		AddItem(entriesPane, 0, 4, false).
		AddItem(contentPane, 0, 5, false).
		AddItem(newNarrowFooterBorder(lineForeground), 1, 0, false)

	wideFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(feedsPane, 45, 0, false).
				AddItem(newPaneDivider(lineForeground), 1, 0, false).
				AddItem(
					tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(entriesPane, 0, 1, false).
						AddItem(contentPane, 0, 2, false),
					0, 1, false,
				),
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

	// FIXME: Remove this workaround that makes display unresponsive
	//        when https://github.com/rivo/tview/issues/921 is resolved.
	width, _, err := term.GetSize(0)
	if err != nil {
		return err
	}
	if width < wideViewMinWidth {
		// Narrow layout.
		mainPage.AddItem(narrowFlex, 0, 0, 1, 1, 0, 0, false)
	} else {
		// Wide layout.
		mainPage.AddItem(wideFlex, 0, 0, 1, 1, 0, 0, false)
	}

	help1 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[aqua]Feeds pane[-]
[yellow]j/k[-]: Next / previous item
[yellow]p[-]  : Pull current feed
[yellow]P[-]  : Pull all feeds
[yellow]R[-]  : Mark all entries in current feed read
[yellow]s[-]  : Star / unstar feed
[yellow]a[-]  : Add feed
[yellow]e[-]  : Edit feed
[yellow]d[-]  : Delete feed

[aqua]Entries pane[-]
[yellow]j/k[-]: Next / previous entry
[yellow]r[-]  : Mark current entry read
[yellow]u[-]  : Mark current entry unread
[yellow]b[-]  : Add / remove current entry from bookmarks

[aqua]Content pane[-]
[yellow]j/k[-]: Scroll down / up
[yellow]g[-]  : Go to top
[yellow]G[-]  : Go to bottom

[aqua]Global[-]
[yellow]1[-]    : Switch to the feeds pane
[yellow]2[-]    : Switch to the entries pane
[yellow]3[-]    : Switch to the content pane
[yellow]Tab[-]  : Switch to next pane
[yellow]S-Tab[-]: Switch to previous pane
[yellow]X[-]    : Export feeds to OPML
[yellow]I[-]    : Import feeds from OPML
[yellow]h|?[-]  : Toggle this help
[yellow]q[-]    : Quit reader`)

	helpPage := tview.NewFrame(help1).
		SetBorders(1, 1, 0, 0, 2, 2)

	helpPage.SetBorder(true).
		SetBorderColor(helpBorderLineForeground).
		SetTitle(" Keys ").
		SetTitleColor(helpTitleForeground)

	root.
		AddAndSwitchToPage("main", mainPage, true).
		AddPage(
			"help",
			tview.NewGrid().
				SetColumns(0, 55, 0).
				SetRows(0, 35, 0).
				AddItem(helpPage, 1, 1, 1, 1, 0, 0, true),
			true,
			false,
		)

	app := tview.NewApplication()
	app.
		SetInputCapture(
			func(event *tcell.EventKey) *tcell.EventKey {
				er := event.Rune()
				if er == 'h' || er == '?' {
					if fp, _ := root.GetFrontPage(); fp == "help" {
						root.HidePage("help")
					} else {
						root.ShowPage("help")
					}
					return nil
				} else if er == 'q' {
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
	topLeftBorderTip *rune,
) *tview.Box {

	lineStyle := tcell.StyleDefault.Foreground(lineForeground).Background(tcell.ColorBlack)

	var unfocused, focused string
	if text != "" {
		unfocused = fmt.Sprintf(" %s ", text)
		focused = fmt.Sprintf(" • %s ", text)
	} else {
		focused = " • "
	}

	makedrawf := func(
		title string,
		leftPad int,
	) func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, lineStyle)
			}
			if topLeftBorderTip != nil {
				screen.SetContent(x-1, y, *topLeftBorderTip, nil, lineStyle)
			}

			// Write the title text.
			tview.Print(
				screen,
				title,
				x+leftPad,
				y,
				width-2,
				tview.AlignLeft,
				textForeground,
			)

			return x + 1, y + 1, width - 2, height - 1
		}
	}

	box := tview.NewBox().SetDrawFunc(makedrawf(unfocused, 3))

	box.SetFocusFunc(func() { box.SetDrawFunc(makedrawf(focused, 1)) })
	box.SetBlurFunc(func() { box.SetDrawFunc(makedrawf(unfocused, 3)) })

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
