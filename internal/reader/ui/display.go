// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/neon/internal"
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

	longDateFormat = "2 January 2006 - 15:04:05 MST"
)

func (d *Display) setRoot() {
	pages := tview.NewPages()
	d.setMainPage()
	d.setHelpPopup()
	d.setIntroPopup()

	d.aboutPopup = newPopup(d.lang.aboutPopupTitle, d.theme.popupTitleFG, 0, 0)
	d.statsPopup = newPopup(d.lang.statsPopupTitle, d.theme.popupTitleFG, 1, 1)

	pages.
		AddAndSwitchToPage(mainPageName, d.mainPage, true).
		AddPage(helpPageName, d.helpPopup, true, false).
		AddPage(aboutPageName, d.aboutPopup, true, false).
		AddPage(statsPageName, d.statsPopup, true, false).
		AddPage(introPageName, d.introPopup, true, false)

	// FIXME: Remove when we add mainPage proper, currently needed only to see if
	//		  app is drawn.
	pages.SetBorder(true)

	d.root = pages
	d.inner = d.inner.SetRoot(pages, true)
}

func (d *Display) setMainPage() {
	// TODO: Add inner flexes.
	grid := tview.NewGrid().
		SetRows(0).
		SetBorders(false)

	d.mainPage = grid
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
[yellow]h,?[-]     : Toggle this help
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

func (d *Display) setIntroPopup() {
	// TODO: Move some constants here into more commonly-accessible place.
	introText := fmt.Sprintf(`Hello and welcome the %s reader.

For help, press [yellow]h[-] or go to [yellow]https://github.com/bow/neon[-].

To close this message, press [yellow]<Esc>[-].
`, internal.AppName())

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
