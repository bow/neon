// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/entity"
)

type Display struct {
	theme *Theme
	lang  *Lang

	inner      *tview.Application
	root       *tview.Pages
	mainPage   *tview.Grid
	aboutPopup *popup
	helpPopup  *popup
	introPopup *popup
	statsPopup *popup

	bar *statusBar

	handlersSet bool
}

func NewDisplay(screen tcell.Screen, theme string) (*Display, error) {
	th, err := LoadTheme(theme)
	if err != nil {
		return nil, err
	}

	d := Display{
		theme: th,
		lang:  langEN,
		inner: tview.NewApplication().
			EnableMouse(true).
			SetScreen(screen),
	}
	d.setRoot()

	return &d, nil
}

func (d *Display) SetHandlers(globalKeyHandler KeyHandler) {
	d.inner = d.inner.SetInputCapture(globalKeyHandler)
	d.handlersSet = true
}

func (d *Display) Start() error {
	if !d.handlersSet {
		return fmt.Errorf("display key handlers must be set before starting")
	}
	return d.inner.Run()
}

func (d *Display) Stop() {
	d.inner.Stop()
}

func (d *Display) dimMainPage() {
	d.theme.dim()
}

func (d *Display) normalizeMainPage() {
	d.theme.normalize()
}

const (
	mainPageName  = "main"
	aboutPageName = "about"
	helpPageName  = "help"
	introPageName = "intro"
	statsPageName = "stats"

	longDateFormat  = "2 January 2006 - 15:04:05 MST"
	shortDateFormat = "02/Jan/06 15:04"
)

func (d *Display) setRoot() {
	pages := tview.NewPages()
	d.setMainPage()
	d.setHelpPopup()
	d.setIntroPopup()

	d.bar = newStatusBar(d.theme)
	d.addStatusBar()

	d.aboutPopup = newPopup(d.lang.aboutPopupTitle, d.theme.popupTitleFG, 0, 0)
	d.statsPopup = newPopup(d.lang.statsPopupTitle, d.theme.popupTitleFG, 1, 1)

	pages.
		AddAndSwitchToPage(mainPageName, d.mainPage, true).
		AddPage(helpPageName, d.helpPopup, true, false).
		AddPage(aboutPageName, d.aboutPopup, true, false).
		AddPage(statsPageName, d.statsPopup, true, false).
		AddPage(introPageName, d.introPopup, true, false)

	d.root = pages
	d.inner = d.inner.SetRoot(pages, true)
}

func (d *Display) setMainPage() {

	feedsPane := newPane(d.lang.feedsPaneTitle, d.theme, false)
	entriesPane := newPane(d.lang.entriesPaneTitle, d.theme, false)
	readingPane := newPane(d.lang.readingPaneTitle, d.theme, true)

	narrowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(feedsPane, 0, 3, false).
		AddItem(entriesPane, 0, 4, false).
		AddItem(readingPane, 0, 5, false).
		AddItem(newNarrowStatusBarBorder(d.theme), 1, 0, false)

	wideFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(feedsPane, 45, 0, false).
				AddItem(newPaneDivider(d.theme), 1, 0, false).
				AddItem(
					tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(entriesPane, 0, 1, false).
						AddItem(readingPane, 0, 2, false),
					0, 1, false,
				),
			0, 1, false,
		).
		AddItem(newWideStatusBarBorder(d.theme, 45), 1, 0, false)

	grid := tview.NewGrid().
		SetRows(0).
		SetBorders(false).
		AddItem(narrowFlex, 0, 0, 1, 1, 0, 0, false).
		AddItem(wideFlex, 0, 0, 1, 1, 0, d.theme.WideViewMinWidth, false)

	d.mainPage = grid
}

func (d *Display) addStatusBar() {
	d.mainPage.SetRows(0, 1).AddItem(d.bar, 1, 0, 1, 1, 0, 0, false)
}

func (d *Display) setHelpPopup() {
	// TODO: Consider moving the content out into where handlers are defined.
	helpText := `[aqua]Feeds pane[-]
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
[yellow]F[-]       : Set focus to feeds pane
[yellow]E[-]       : Set focus to entries pane
[yellow]R[-]       : Set focus to reading pane
[yellow]Tab[-]     : Switch to next pane
[yellow]Alt-Tab[-] : Switch to previous pane
[yellow]b[-]       : Toggle status bar
[yellow]c[-]       : Clear status bar
[yellow]X[-]       : Export feeds to OPML
[yellow]I[-]       : Import feeds from OPML
[yellow]Esc[-]     : Unset current focus or close open frame
[yellow]S[-]       : Toggle stats popup
[yellow]A[-]       : Toggle 'about' popup
[yellow]H,?[-]     : Toggle this help
[yellow]q,Ctrl-C[-]: Quit reader`

	helpWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(helpText)

	d.helpPopup = newFilledPopup(
		d.lang.helpPopupTitle,
		helpWidget,
		d.theme.popupTitleFG,
		1, 1,
		popupWidth(helpWidget.GetText(true)),
		[]int{0, popupHeight(helpText), 0},
	)
}

func (d *Display) setIntroPopup() {
	// TODO: Move some constants here into more commonly-accessible place.
	introText := fmt.Sprintf(`Hello and welcome the %s reader.

For help, press [yellow]?[-] or go to [yellow]%s[-].

To close this message, press [yellow]<Esc>[-].
`, internal.AppName(), internal.AppHomepage())

	introWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(introText)

	d.introPopup = newFilledPopup(
		d.lang.introPopupTitle,
		introWidget,
		d.theme.popupTitleFG,
		1, 1,
		popupWidth(introWidget.GetText(true)),
		[]int{-1, popupHeight(introText), -3},
	)
}

func (d *Display) setAboutPopupText(name fmt.Stringer) {
	commit := internal.GitCommit()

	var buildTime = internal.BuildTime()
	buildTimeVal, err := time.Parse(time.RFC3339, buildTime)
	if err == nil {
		buildTime = buildTimeVal.Format(longDateFormat)
	}

	infoText := fmt.Sprintf(`[yellow]Version[-]   : %s
[yellow]Git commit[-]: %s
[yellow]Build time[-]: %s
[yellow]Backend[-]   : %s`,
		internal.Version(),
		commit,
		buildTime,
		name,
	)

	aboutWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(infoText)

	// NOTE: We assume the banner's width is less than the one computed here.
	width := popupWidth(aboutWidget.GetText(true))
	banner := centerBanner(internal.Banner(), width)
	aboutText := fmt.Sprintf("%s\n\n%s", banner, infoText)

	aboutWidget.SetText(aboutText)

	height := popupHeight(aboutText) - 1

	d.aboutPopup.setWidth(width)
	d.aboutPopup.setGridRows([]int{-1, height, -3})
	d.aboutPopup.setContent(aboutWidget)
}

func (d *Display) setStatsPopupValues(values *entity.Stats) {

	var lpt string
	if values.LastPullTime != nil {
		lpt = values.LastPullTime.Format(longDateFormat)
	}

	statsText := fmt.Sprintf(`[aqua]Feeds[-]
[yellow]Total[-]: %d

[aqua]Entries[-]
[yellow]Unread[-]: %d
[yellow]Total[-] : %d

[aqua]Last pulled[-]
%s`,
		values.NumFeeds,
		values.NumEntriesUnread,
		values.NumEntries,
		lpt,
	)

	statsWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(statsText)

	width := popupWidth(statsWidget.GetText(true))
	height := popupHeight(statsText)

	d.statsPopup.setWidth(width)
	d.statsPopup.setGridRows([]int{-1, height, -3})
	d.statsPopup.setContent(statsWidget)
}

func newPane(title string, theme *Theme, addTopLeftBorderTip bool) *tview.Box {

	var unfocused, focused string
	if title != "" {
		unfocused = fmt.Sprintf(" %s ", title)
		focused = fmt.Sprintf("[::b]▶ %s[::-] ", title)
	} else {
		focused = "[::b]▶[::-] "
	}

	makedrawf := func(
		title string,
		leftPad int,
	) func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			style := theme.lineStyle()
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
			}
			if addTopLeftBorderTip {
				screen.SetContent(x-1, y, tview.BoxDrawingsLightVerticalAndRight, nil, style)
			}

			// Write the title text.
			tview.Print(
				screen,
				title,
				x+leftPad,
				y,
				width-2,
				tview.AlignLeft,
				theme.titleFG,
			)

			return x + 2, y + 1, width - 2, height - 1
		}
	}

	box := tview.NewBox().SetDrawFunc(makedrawf(unfocused, 1))

	box.SetFocusFunc(func() { box.SetDrawFunc(makedrawf(focused, 0)) })
	box.SetBlurFunc(func() { box.SetDrawFunc(makedrawf(unfocused, 1)) })

	return box
}

func newNarrowStatusBarBorder(theme *Theme) *tview.Box {

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		style := theme.lineStyle()
		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
		}
		return x + 1, y + 1, width - 2, height - 1
	}

	return tview.NewBox().SetBorder(false).SetDrawFunc(drawf)
}

func newWideStatusBarBorder(theme *Theme, branchPoint int) *tview.Box {

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		style := theme.lineStyle()
		for cx := x; cx < x+width; cx++ {
			if cx == branchPoint {
				screen.SetContent(cx, y, tview.BoxDrawingsLightUpAndHorizontal, nil, style)
			} else {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
			}
		}

		return x + 1, y + 1, width - 2, height - 1
	}

	return tview.NewBox().SetBorder(false).SetDrawFunc(drawf)
}

func newPaneDivider(theme *Theme) *tview.Box {

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		style := theme.lineStyle()
		screen.SetContent(x, y, tview.BoxDrawingsLightDownAndHorizontal, nil, style)
		for cy := y + 1; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, style)
		}
		return x + 1, y + 1, width - 2, height - 1
	}

	return tview.NewBox().SetBorder(false).SetDrawFunc(drawf)
}
