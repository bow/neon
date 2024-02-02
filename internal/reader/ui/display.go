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

	inner *tview.Application
	root  *tview.Pages

	mainPage *tview.Grid

	feedsCh   chan *entity.Feed
	feedsPane *feedsPane

	entriesPane *tview.Box
	readingPane *tview.Box

	bar        *statusBar
	barVisible bool
	eventsCh   chan *event

	aboutPopup *popup
	helpPopup  *popup
	introPopup *popup
	statsPopup *popup

	handlersSet bool

	focusStack tview.Primitive
	counter    int
}

func NewDisplay(screen tcell.Screen, theme string) (*Display, error) {
	th, err := loadTheme(theme)
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
	d.eventsCh = make(chan *event)
	d.feedsCh = make(chan *entity.Feed)

	return &d, nil
}

func (d *Display) SetHandlers(
	globalKeyHandler KeyHandler,
	feedsPaneKeyHandler KeyHandler,
) {
	d.inner.SetInputCapture(globalKeyHandler)
	d.feedsPane.SetInputCapture(feedsPaneKeyHandler)
	d.handlersSet = true
}

func (d *Display) Start() error {
	if !d.handlersSet {
		return fmt.Errorf("display key handlers must be set before starting")
	}
	stop := d.startEventPoll()
	defer stop()
	// TODO: Consider making this similar to event polling.
	go d.feedsPane.startFeedsPoll(d.feedsCh)

	return d.inner.Run()
}

func (d *Display) Draw() {
	d.inner.Draw()
}

func (d *Display) Stop() {
	d.inner.Stop()
}

func (d *Display) dimMainPage() {
	d.theme.dim()
	d.feedsPane.refreshColors()
	d.bar.refreshColors()
}

func (d *Display) normalizeMainPage() {
	d.theme.normalize()
	d.feedsPane.refreshColors()
	d.bar.refreshColors()
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
	d.bar.setChangedFunc(func() { d.inner.Draw() })
	d.addStatusBar()

	d.aboutPopup = newPopup(
		d.lang.aboutPopupTitle,
		d.theme.popupTitleFG,
		0, 0,
		-1, -3,
	)
	d.statsPopup = newPopup(
		d.lang.statsPopupTitle,
		d.theme.popupTitleFG,
		1, 1,
		-1, -3,
	)

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

	d.feedsCh = make(chan *entity.Feed)
	feedsPane := newFeedsPane(d.theme, d.lang)

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
		AddItem(wideFlex, 0, 0, 1, 1, 0, d.theme.wideViewMinWidth, false)

	d.mainPage = grid
	d.feedsPane = feedsPane
	d.entriesPane = entriesPane
	d.readingPane = readingPane
}

func (d *Display) startEventPoll() (stop func()) {
	done := make(chan struct{})
	stop = func() {
		defer close(done)
		done <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case ev := <-d.eventsCh:
				d.bar.showEvent(ev)
			}
		}
	}()

	return stop
}

func (d *Display) addStatusBar() {
	d.mainPage.SetRows(0, 1).AddItem(d.bar, 1, 0, 1, 1, 0, 0, false)
	d.barVisible = true
}

func (d *Display) removeStatusBar() {
	d.mainPage.RemoveItem(d.bar).SetRows(0)
	d.barVisible = false
}

func (d *Display) toggleStatusBar() {
	if d.barVisible {
		d.removeStatusBar()
	} else {
		d.addStatusBar()
	}
}

func (d *Display) clearEvent() {
	d.bar.clearLatestEvent()
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
[yellow]S[-]       : Toggle stats popup and show latest values
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
		popupWidth(helpWidget.GetText(true)),
		popupHeight(helpText),
		1, 1,
		0, 0,
	)
}

func (d *Display) setIntroPopup() {
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
		popupWidth(introWidget.GetText(true)),
		popupHeight(introText),
		1, 1,
		-1, -3,
	)
}

func (d *Display) setAboutPopupText(name string) {
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

	d.aboutPopup.setWidth(width)
	d.aboutPopup.setHeight(popupHeight(aboutText) - 1)
	d.aboutPopup.setContent(aboutWidget)
}

func (d *Display) setStats(stats *entity.Stats) {
	d.setStatsPopupValues(stats)
	d.bar.setStats(stats)
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

	d.statsPopup.setWidth(popupWidth(statsWidget.GetText(true)))
	d.statsPopup.setHeight(popupHeight(statsText))
	d.statsPopup.setContent(statsWidget)
}

func (d *Display) frontPageName() string {
	name, _ := d.root.GetFrontPage()
	return name
}

func (d *Display) switchPopup(name string, currentFront string) {
	if currentFront == mainPageName {
		d.stashFocus()
	} else {
		d.root.HidePage(currentFront)
	}
	d.showPopup(name)
}

func (d *Display) showPopup(name string) {
	d.dimMainPage()
	d.root.ShowPage(name)
}

func (d *Display) hidePopup(name string) {
	d.root.HidePage(name)
	d.normalizeMainPage()
	if top := d.focusStack; top != nil {
		d.inner.SetFocus(top)
	}
	d.focusStack = nil
}

func (d *Display) stashFocus() {
	d.focusStack = d.inner.GetFocus()
}

func (d *Display) focusPane(pane tview.Primitive) {
	front := d.frontPageName()
	if front != mainPageName {
		d.root.HidePage(front)
		d.normalizeMainPage()
	}
	d.inner.SetFocus(pane)
	d.stashFocus()
}

func (d *Display) focusAdjacentPane(reverse bool) {
	d.counter++
	if front := d.frontPageName(); front != mainPageName {
		d.hidePopup(front)
	}
	targets := []tview.Primitive{d.feedsPane, d.entriesPane, d.readingPane}
	current := d.inner.GetFocus()
	idx := 0
	if reverse {
		switch current {
		case d.entriesPane:
			idx = 0
		case d.readingPane:
			idx = 1
		default:
			idx = 2
		}
	} else {
		switch current {
		case d.entriesPane:
			idx = 2
		case d.feedsPane:
			idx = 1
		default:
			idx = 0
		}
	}
	d.inner.SetFocus(targets[idx])
}

func (d *Display) infoEventf(text string, a ...any) { d.eventf(eventLevelInfo, text, a...) }

func (d *Display) warnEventf(text string, a ...any) { d.eventf(eventLevelWarn, text, a...) }

func (d *Display) errEventf(text string, a ...any) {
	d.eventf(eventLevelErr, fmt.Sprintf(text, a...))
}

func (d *Display) errEvent(err error) {
	d.eventf(eventLevelErr, fmt.Sprintf("%s", err))
}

func (d *Display) eventf(level eventLevel, text string, a ...any) {
	ev := event{level: level, timestamp: time.Now(), text: fmt.Sprintf(text, a...)}
	go func() { d.eventsCh <- &ev }()
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

type event struct {
	level     eventLevel
	timestamp time.Time
	text      string
}

type eventLevel uint8

const (
	eventLevelInfo eventLevel = iota
	eventLevelWarn
	eventLevelErr
)
