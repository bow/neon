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
// TODO: Refactor and split UI components.
func Show(db internal.FeedStore) error { //nolint:revive

	theme := DefaultTheme

	root := tview.NewPages()

	topLeftBorderTip := tview.BoxDrawingsLightVerticalAndRight
	feedsPane := newPane(theme.FeedsPaneTitle, theme, nil)
	entriesPane := newPane(theme.EntriesPaneTitle, theme, nil)
	readingPane := newPane(theme.ReadingPaneTitle, theme, &topLeftBorderTip)

	narrowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(feedsPane, 0, 3, false).
		AddItem(entriesPane, 0, 4, false).
		AddItem(readingPane, 0, 5, false).
		AddItem(newNarrowFooterBorder(theme), 1, 0, false)

	wideFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(feedsPane, 45, 0, false).
				AddItem(newPaneDivider(theme), 1, 0, false).
				AddItem(
					tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(entriesPane, 0, 1, false).
						AddItem(readingPane, 0, 2, false),
					0, 1, false,
				),
			0, 1, false,
		).
		AddItem(newWideFooterBorder(theme, 45), 1, 0, false)

	stats, err := db.GetGlobalStats(context.Background())
	if err != nil {
		return err
	}

	unreadInfo := tview.NewTextView().
		SetTextColor(theme.StatsForeground).
		SetText(fmt.Sprintf("%d unread entries", stats.NumEntriesUnread))

	lastPullInfo := tview.NewTextView().
		SetTextColor(theme.LastPullForeground).
		SetText(
			fmt.Sprintf("Pulled %s", stats.LastPullTime.Local().Format("02/Jan/06 15:04")),
		)

	versionInfo := tview.NewTextView().
		SetTextColor(theme.VersionForeground).
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
	mainPage.AddItem(narrowFlex, 0, 0, 1, 1, 0, 0, false)

	// Wide layout.
	mainPage.AddItem(wideFlex, 0, 0, 1, 1, 0, theme.WideViewMinWidth, false)

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

[aqua]Reading pane[-]
[yellow]j/k[-]: Scroll down / up
[yellow]g[-]  : Go to top
[yellow]G[-]  : Go to bottom

[aqua]Global[-]
[yellow]1|F[-]  : Switch to the feeds pane
[yellow]2|E[-]  : Switch to the entries pane
[yellow]3|R[-]  : Switch to the reading pane
[yellow]Tab[-]  : Switch to next pane
[yellow]A-Tab[-]: Switch to previous pane
[yellow]X[-]    : Export feeds to OPML
[yellow]I[-]    : Import feeds from OPML
[yellow]Esc[-]  : Unset current focus or close open frame
[yellow]h|?[-]  : Toggle this help
[yellow]q[-]    : Quit reader`)

	helpPage := tview.NewFrame(help1).
		SetBorders(1, 1, 0, 0, 2, 2)

	helpPage.SetBorder(true).
		SetBorderColor(theme.PopupBorderForeground).
		SetTitle(fmt.Sprintf(" %s ", theme.HelpPopupTitle)).
		SetTitleColor(theme.PopupTitleForeground)

	const (
		mainPageName = "main"
		helpPageName = "help"
	)

	root.
		AddAndSwitchToPage(mainPageName, mainPage, true).
		AddPage(
			helpPageName,
			tview.NewGrid().
				SetColumns(0, 55, 0).
				SetRows(0, 36, 0).
				AddItem(helpPage, 1, 1, 1, 1, 0, 0, true),
			true,
			false,
		)

	app := tview.NewApplication()

	panesMap := map[rune]*tview.Box{
		'1': feedsPane,
		'F': feedsPane,
		'2': entriesPane,
		'E': entriesPane,
		'3': readingPane,
		'R': readingPane,
	}
	panesOrder := []*tview.Box{feedsPane, entriesPane, readingPane}

	// nolint:exhaustive
	eventHandler := func(event *tcell.EventKey) *tcell.EventKey {

		var (
			focused  = app.GetFocus()
			front, _ = root.GetFrontPage()
			key      = event.Key()
			keyr     = event.Rune()
		)

		switch focused {

		case feedsPane:
			if keyr == 'P' {
				// TODO: Add animation in footer?
				ch := db.PullFeeds(context.Background(), []internal.ID{})
				// TODO: Add ok / fail status in ...?
				for pr := range ch {
					if err := pr.Error(); err != nil {
						panic(err)
					}
				}
				stats, err := db.GetGlobalStats(context.Background())
				if err != nil {
					panic(err)
				}
				unreadInfo.SetText(fmt.Sprintf("%d unread entries", stats.NumEntriesUnread))
				lastPullInfo.SetText(
					fmt.Sprintf(
						"Pulled %s", stats.LastPullTime.Local().Format("02/Jan/06 15:04"),
					),
				)
				return nil
			}

		default:
			switch key {

			case tcell.KeyRune:
				switch keyr {
				case '1', '2', '3', 'F', 'E', 'R':
					if front == mainPageName {
						app.SetFocus(panesMap[keyr])
					}
					return nil

				case 'h', '?':
					if front == helpPageName {
						root.HidePage(helpPageName)
					} else {
						root.ShowPage(helpPageName)
					}
					return nil

				case 'q':
					app.Stop()
					return nil
				}

			case tcell.KeyTab:
				if front == mainPageName {
					target := 0
					if event.Modifiers()&tcell.ModAlt != 0 {
						switch focused {
						case nil, mainPage, feedsPane:
							target = 2
						case entriesPane:
							target = 0
						case readingPane:
							target = 1
						}
					} else {
						switch focused {
						case nil, mainPage, readingPane:
							target = 0
						case entriesPane:
							target = 2
						case feedsPane:
							target = 1
						}
					}
					app.SetFocus(panesOrder[target])
				}
				return nil

			case tcell.KeyEscape:
				switch front {
				case helpPageName:
					root.HidePage(helpPageName)
					return nil
				case mainPageName:
					app.SetFocus(nil)
					return nil
				}
			}
		}

		return event
	}

	app.SetInputCapture(eventHandler)

	if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}

type Theme struct {
	FeedsPaneTitle   string
	EntriesPaneTitle string
	ReadingPaneTitle string
	HelpPopupTitle   string

	Background            tcell.Color
	BorderForeground      tcell.Color
	TitleForeground       tcell.Color
	VersionForeground     tcell.Color
	LastPullForeground    tcell.Color
	StatsForeground       tcell.Color
	PopupTitleForeground  tcell.Color
	PopupBorderForeground tcell.Color

	WideViewMinWidth int
}

var DefaultTheme = &Theme{
	FeedsPaneTitle:   "Feeds",
	EntriesPaneTitle: "Entries",
	ReadingPaneTitle: "",
	HelpPopupTitle:   "Keys",

	Background:            tcell.ColorBlack,
	BorderForeground:      tcell.ColorWhite,
	TitleForeground:       tcell.ColorBlue,
	VersionForeground:     tcell.ColorGray,
	LastPullForeground:    tcell.ColorGray,
	StatsForeground:       tcell.ColorDarkGoldenrod,
	PopupBorderForeground: tcell.ColorGray,
	PopupTitleForeground:  tcell.ColorAqua,

	WideViewMinWidth: 150,
}

// TODO: Consider moving to theme.
func init() {
	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical
	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight
}

func newPane(
	title string,
	theme *Theme,
	topLeftBorderTip *rune,
) *tview.Box {

	lineStyle := tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.BorderForeground)

	var unfocused, focused string
	if title != "" {
		unfocused = fmt.Sprintf(" %s ", title)
		focused = fmt.Sprintf(" • %s ", title)
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
				theme.TitleForeground,
			)

			return x + 1, y + 1, width - 2, height - 1
		}
	}

	box := tview.NewBox().SetDrawFunc(makedrawf(unfocused, 3))

	box.SetFocusFunc(func() { box.SetDrawFunc(makedrawf(focused, 1)) })
	box.SetBlurFunc(func() { box.SetDrawFunc(makedrawf(unfocused, 3)) })

	return box
}

func newPaneDivider(theme *Theme) *tview.Box {

	style := tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.BorderForeground)

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

func newNarrowFooterBorder(theme *Theme) *tview.Box {

	style := tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.BorderForeground)

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
		}

		return x + 1, y + 1, width - 2, height - 1
	}

	divider := tview.NewBox().SetBorder(false).SetDrawFunc(drawf)

	return divider
}

func newWideFooterBorder(theme *Theme, branch int) *tview.Box {

	style := tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.BorderForeground)

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
