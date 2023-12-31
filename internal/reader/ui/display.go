// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Display struct {
	theme  *Theme
	lang   *Lang
	screen tcell.Screen

	inner      *tview.Application
	root       *tview.Pages
	mainPage   *tview.Grid
	aboutPopup *popup
	helpPopup  *popup

	initialized bool
}

func NewDisplay(screen tcell.Screen, theme string) (*Display, error) {
	th, err := LoadTheme(theme)
	if err != nil {
		return nil, err
	}

	d := Display{
		theme:  th,
		lang:   langEN,
		screen: screen,
		inner: tview.NewApplication().
			EnableMouse(true).
			SetScreen(screen),
	}
	d.setRoot()

	return &d, nil
}

func (d *Display) Init(globalKeyHandler KeyHandler) {
	d.inner = d.inner.SetInputCapture(globalKeyHandler)
	d.initialized = true
}

func (d *Display) Start() error {
	if !d.initialized {
		return fmt.Errorf("display must be initialized before starting")
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

	longDateFormat = "2 January 2006 - 15:04:05 MST"
)

func (d *Display) setRoot() {
	pages := tview.NewPages()
	d.setMainPage()
	d.setHelpPopup()
	d.setAboutPopup()

	pages.
		AddAndSwitchToPage(mainPageName, d.mainPage, true).
		AddPage(helpPageName, d.helpPopup.grid, true, false).
		AddPage(aboutPageName, d.aboutPopup.grid, true, false)

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

func (d *Display) setAboutPopup() {
	d.aboutPopup = newEmptyPopup(d.lang.aboutPopupTitle, d.theme.popupTitleFG, 0, 0)
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

	d.helpPopup = newPopup(
		d.lang.helpPopupTitle,
		helpWidget,
		d.theme.popupTitleFG,
		1, 1,
		55,
		[]int{0, popupHeight(helpText), 0},
	)
}
