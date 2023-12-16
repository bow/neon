// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package tui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/lens/internal"
)

const (
	mainPageName    = "main"
	helpPageName    = "help"
	statsPageName   = "stats"
	aboutPageName   = "about"
	welcomePageName = "welcome"

	shortDateFormat = "02/Jan/06 15:04"
	longDateFormat  = "2 January 2006 - 15:04:05 MST"

	verticalPopupPadding = 4
)

type Reader struct {
	ctx      context.Context
	store    internal.FeedStore
	theme    *Theme
	dbPath   *string
	initPath string

	app *tview.Application

	root      *tview.Pages
	mainPage  *tview.Grid
	helpPage  *tview.Grid
	statsPage *tview.Grid
	aboutPage *tview.Grid

	feedsPane   *feedsPane
	entriesPane *tview.Box
	readingPane *tview.Box
	bar         *statusBar

	statsCache *internal.Stats
}

func NewReader(ctx context.Context, store internal.FeedStore, dbPath *string) *Reader {

	reader := Reader{
		ctx:    ctx,
		store:  store,
		theme:  DarkTheme,
		root:   tview.NewPages(),
		app:    tview.NewApplication(),
		dbPath: dbPath,
	}

	reader.setupMainPage()
	reader.setupHelpPage()
	reader.setupStatsPage()
	reader.setupAboutPage()

	reader.root.
		AddAndSwitchToPage(mainPageName, reader.mainPage, true).
		AddPage(helpPageName, reader.helpPage, true, false).
		AddPage(statsPageName, reader.statsPage, true, false).
		AddPage(aboutPageName, reader.aboutPage, true, false)

	reader.root.SetInputCapture(reader.globalKeyHandler())
	reader.app.SetRoot(reader.root, true).EnableMouse(true)

	return &reader
}

func (r *Reader) Show() error {
	stats, err := r.store.GetGlobalStats(r.ctx)
	if err != nil {
		return err
	}
	r.bar.updateFromStats(stats)
	if !r.isInitialized() {
		welcomeText := fmt.Sprintf(`Hello and welcome the %s reader.

For help, press [yellow]h[-] or go to [yellow]https://github.com/bow/lens[-].

To close this message, press [yellow]<Esc>[-].
`, internal.AppName())

		r.root.AddPage(
			welcomePageName,
			r.newPopup(
				r.theme.WelcomePopupTitle,
				tview.NewTextView().SetDynamicColors(true).SetText(welcomeText),
				1, 1,
				61,
				[]int{-1, calcPopupHeight(welcomeText), -3},
			),
			true,
			false,
		)
		r.dimColors()
		r.root.ShowPage(welcomePageName)
		defer r.initialize()
	}

	return r.app.Run()
}

func (r *Reader) WithInitPath(path string) *Reader {
	r.initPath = path
	return r
}

func (r *Reader) setupMainPage() {

	feedsPane := newFeedsPane(r.theme)
	feedsPane.SetInputCapture(r.feedsPaneKeyHandler())

	entriesPane := r.newPane(r.theme.EntriesPaneTitle, false)
	readingPane := r.newPane(r.theme.ReadingPaneTitle, true)

	narrowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(feedsPane, 0, 3, false).
		AddItem(entriesPane, 0, 4, false).
		AddItem(readingPane, 0, 5, false).
		AddItem(r.newNarrowStatusBarBorder(), 1, 0, false)

	wideFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(feedsPane, 45, 0, false).
				AddItem(r.newPaneDivider(), 1, 0, false).
				AddItem(
					tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(entriesPane, 0, 1, false).
						AddItem(readingPane, 0, 2, false),
					0, 1, false,
				),
			0, 1, false,
		).
		AddItem(r.newWideStatusBarBorder(45), 1, 0, false)

	mainPage := tview.NewGrid().
		SetRows(0).
		SetBorders(false).
		AddItem(narrowFlex, 0, 0, 1, 1, 0, 0, false).
		AddItem(wideFlex, 0, 0, 1, 1, 0, r.theme.WideViewMinWidth, false)

	r.feedsPane = feedsPane
	r.entriesPane = entriesPane
	r.readingPane = readingPane

	r.mainPage = mainPage

	r.bar = newStatusBar(r.theme).
		setChangedFunc(func() { r.app.Draw() }).
		addToMainPage(r.mainPage)
}

func (r *Reader) setupHelpPage() {

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
[yellow]1,F[-]     : Toggle feeds pane focus
[yellow]2,E[-]     : Toggle entries pane focus
[yellow]3,R[-]     : Toggle reading pane focus
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

	r.helpPage = r.newPopup(
		r.theme.HelpPopupTitle,
		helpWidget,
		1, 1,
		55,
		[]int{0, calcPopupHeight(helpText), 0},
	)
}

func (r *Reader) setupStatsPage() {

	if r.statsCache == nil {
		stats, err := r.store.GetGlobalStats(r.ctx)
		if err != nil {
			panic(err)
		}
		r.statsCache = stats
	}

	statsText := fmt.Sprintf(`[aqua]Feeds[-]
[yellow]Total[-]: %d

[aqua]Entries[-]
[yellow]Unread[-]: %d
[yellow]Total[-] : %d

[aqua]Last pulled[-]
%s`,
		r.statsCache.NumFeeds,
		r.statsCache.NumEntries,
		r.statsCache.NumEntriesUnread,
		r.statsCache.LastPullTime.Format(longDateFormat),
	)

	statsWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(statsText)

	r.statsPage = r.newPopup(
		r.theme.StatsPopupTitle,
		statsWidget,
		1, 1,
		37,
		[]int{-1, calcPopupHeight(statsText), -3},
	)
}

func (r *Reader) setupAboutPage() {

	commit := internal.GitCommit()

	var buildTime = internal.BuildTime()
	buildTimeVal, err := time.Parse(time.RFC3339, buildTime)
	if err == nil {
		buildTime = buildTimeVal.Format(longDateFormat)
	}

	width := len(commit) + 18

	aboutText := fmt.Sprintf(`%s

[yellow]Version[-]   : %s
[yellow]Git commit[-]: %s
[yellow]Build time[-]: %s`,
		centerBanner(internal.Banner(), width),
		internal.Version(),
		commit,
		buildTime,
	)
	if dbPath := r.dbPath; dbPath != nil {
		aboutText += fmt.Sprintf("\n[yellow]Database[-]  : %s", *dbPath)
	}

	aboutWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(aboutText)

	r.aboutPage = r.newPopup(
		r.theme.AboutPopupTitle,
		aboutWidget,
		0, 0,
		width,
		[]int{-1, calcPopupHeight(aboutText) - 1, -3},
	)
}

// nolint:revive,exhaustive
func (r *Reader) globalKeyHandler() func(event *tcell.EventKey) *tcell.EventKey {

	return func(event *tcell.EventKey) *tcell.EventKey {
		var (
			focused  = r.app.GetFocus()
			front, _ = r.root.GetFrontPage()
			key      = event.Key()
			keyr     = event.Rune()
		)

		switch key {

		case tcell.KeyRune:
			switch keyr {
			case '1', '2', '3', 'F', 'E', 'R':
				if front != mainPageName {
					r.root.HidePage(front)
					r.normalizeColors()
					front = mainPageName
				}
				if front == mainPageName {
					target := r.getFocusTarget(keyr)
					if target != r.app.GetFocus() {
						r.app.SetFocus(target)
					} else {
						r.app.SetFocus(r.root)
					}
				}
				return nil

			case 'S':
				if front == statsPageName {
					r.root.HidePage(front)
					r.normalizeColors()
				} else if front != welcomePageName {
					if front != mainPageName {
						r.root.HidePage(front)
						r.normalizeColors()
					}
					r.dimColors()
					r.root.ShowPage(statsPageName)
				}
				return nil

			case 'A':
				if front == aboutPageName {
					r.root.HidePage(front)
					r.normalizeColors()
				} else if front != welcomePageName {
					if front != mainPageName {
						r.root.HidePage(front)
						r.normalizeColors()
					}
					r.dimColors()
					r.bar.refreshColors()
					r.root.ShowPage(aboutPageName)
				}
				return nil

			case 'b':
				r.bar.Lock()
				defer r.bar.Unlock()
				r.bar.toggleFromMainPage(r.mainPage)
				return nil

			case 'c':
				r.bar.Lock()
				defer r.bar.Unlock()
				r.bar.clearActivity()
				return nil

			case 'h', '?':
				if front == helpPageName {
					r.root.HidePage(front)
					r.normalizeColors()
				} else {
					if front != mainPageName {
						r.root.HidePage(front)
						r.normalizeColors()
					}
					r.dimColors()
					r.root.ShowPage(helpPageName)
				}
				return nil

			case 'q':
				r.app.Stop()
				return nil
			}

		case tcell.KeyTab:
			if front != mainPageName {
				r.root.HidePage(front)
				front = mainPageName
			}
			if front == mainPageName {
				reverse := event.Modifiers()&tcell.ModAlt != 0
				target := r.getAdjacentFocusTarget(focused, reverse)
				r.app.SetFocus(target)
			}
			return nil

		case tcell.KeyEscape:
			switch front {
			case mainPageName, "":
				r.app.SetFocus(r.root)
			default:
				r.root.HidePage(front)
				r.normalizeColors()
			}
			return nil
		}

		return event
	}
}

// nolint:revive,exhaustive
func (r *Reader) feedsPaneKeyHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		keyr := event.Rune()

		if keyr == 'P' {

			go func() {
				r.bar.Lock()
				defer r.bar.Unlock()

				r.bar.showNormalActivity("Pulling feeds")

				var count int
				ch := r.store.PullFeeds(r.ctx, []internal.ID{})
				for pr := range ch {
					if err := pr.Error(); err != nil {
						// TODO: Add ok / fail status in ...?
						panic(err)
					}
					r.bar.showNormalActivity("Pulling: %s done", pr.URL())
					count++
				}
				switch count {
				case 0:
					r.bar.showNormalActivity("No feeds to pull")
				case 1:
					r.bar.showNormalActivity("Pulled %d feed successfully", count)
				default:
					r.bar.showNormalActivity("Pulled %d feeds successfully", count)
				}

				stats, err := r.store.GetGlobalStats(r.ctx)
				if err != nil {
					panic(err)
				}
				r.statsCache = stats
				r.bar.updateFromStats(stats)
			}()
			return nil
		}
		return event
	}
}

func (r *Reader) getFocusTarget(keyr rune) tview.Primitive {
	var target tview.Primitive
	switch keyr { // nolint:exhaustive
	case '1', 'F':
		target = r.feedsPane
	case '2', 'E':
		target = r.entriesPane
	case '3', 'R':
		target = r.readingPane
	default:
		panic(fmt.Sprintf("unexpected key: %c", keyr))
	}
	return target
}

func (r *Reader) getAdjacentFocusTarget(
	current tview.Primitive,
	reverse bool,
) tview.Primitive {
	targets := []tview.Primitive{r.feedsPane, r.entriesPane, r.readingPane}
	idx := 0
	if reverse {
		switch current {
		case r.entriesPane:
			idx = 0
		case r.readingPane:
			idx = 1
		default:
			idx = 2
		}
	} else {
		switch current {
		case r.entriesPane:
			idx = 2
		case r.feedsPane:
			idx = 1
		default:
			idx = 0
		}
	}
	return targets[idx]
}

func (r *Reader) dimColors() {
	r.theme.Dim()
	r.bar.refreshColors()
}

func (r *Reader) normalizeColors() {
	r.theme.Normalize()
	r.bar.refreshColors()
}

func (r *Reader) newPane(title string, addTopLeftBorderTip bool) *tview.Box {

	var unfocused, focused string
	if title != "" {
		unfocused = fmt.Sprintf(" %s ", title)
		focused = fmt.Sprintf("[::b]» %s[::-] ", title)
	} else {
		focused = "[::b]»[::-] "
	}

	makedrawf := func(
		title string,
		leftPad int,
	) func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			lineStyle := r.theme.lineStyle()
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, lineStyle)
			}
			if addTopLeftBorderTip {
				screen.SetContent(x-1, y, tview.BoxDrawingsLightVerticalAndRight, nil, lineStyle)
			}

			// Write the title text.
			tview.Print(
				screen,
				title,
				x+leftPad,
				y,
				width-2,
				tview.AlignLeft,
				r.theme.TitleForeground,
			)

			return x + 2, y + 1, width - 2, height - 1
		}
	}

	box := tview.NewBox().SetDrawFunc(makedrawf(unfocused, 1))

	box.SetFocusFunc(func() { box.SetDrawFunc(makedrawf(focused, 0)) })
	box.SetBlurFunc(func() { box.SetDrawFunc(makedrawf(unfocused, 1)) })

	return box
}

const (
	leftPopupMargin  = 2
	rightPopupMargin = 2
)

func (r *Reader) newPopup(
	title string,
	contents *tview.TextView,
	top, bottom int,
	ncols int,
	gridRows []int,
) *tview.Grid {

	frame := tview.NewFrame(contents).
		SetBorders(top, bottom, 0, 0, leftPopupMargin, rightPopupMargin)

	frame.SetBorder(true).
		SetTitle(r.makeTitle(title)).
		SetTitleColor(r.theme.PopupTitleForeground)

	return tview.NewGrid().
		SetColumns(0, ncols, 0).
		SetRows(gridRows...).
		AddItem(frame, 1, 1, 1, 1, 0, 0, true)
}

func (r *Reader) newPaneDivider() *tview.Box {

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		style := r.theme.lineStyle()
		screen.SetContent(x, y, tview.BoxDrawingsLightDownAndHorizontal, nil, style)
		for cy := y + 1; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, style)
		}
		return x + 1, y + 1, width - 2, height - 1
	}

	return tview.NewBox().SetBorder(false).SetDrawFunc(drawf)
}

func (r *Reader) newNarrowStatusBarBorder() *tview.Box {

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		style := r.theme.lineStyle()
		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
		}
		return x + 1, y + 1, width - 2, height - 1
	}

	return tview.NewBox().SetBorder(false).SetDrawFunc(drawf)
}

func (r *Reader) newWideStatusBarBorder(branchPoint int) *tview.Box {

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		style := r.theme.lineStyle()
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

func (r *Reader) makeTitle(text string) string {
	return fmt.Sprintf(" %s ", text)
}

func (r *Reader) isInitialized() bool {
	// Reader default is to assume already initialized.
	if r.initPath == "" {
		return true
	}
	exists := true
	_, err := os.Stat(r.initPath)
	if err != nil && os.IsNotExist(err) {
		exists = false
	}
	return exists
}

func (r *Reader) initialize() {
	if p := r.initPath; p != "" {
		os.Create(r.initPath) // nolint:errcheck
	}
}

func centerBanner(text string, width int) string {
	if width == 0 {
		return text
	}

	var (
		maxLineWidth = 0
		lines        = make([]string, 0)
		sc           = bufio.NewScanner(strings.NewReader(text))
	)
	for sc.Scan() {
		line := sc.Text()
		if ncols := len(line); ncols > maxLineWidth {
			maxLineWidth = ncols
		}
		lines = append(lines, line)
	}

	if maxLineWidth > width {
		return text
	}

	leftPad := strings.Repeat(" ", ((width-maxLineWidth)/2)-leftPopupMargin)
	paddedLines := make([]string, len(lines))
	for i, line := range lines {
		paddedLines[i] = fmt.Sprintf("%s%s", leftPad, line)
	}

	sep := "\n"
	if runtime.GOOS == "windows" {
		sep = "\r\n"
	}

	return strings.Join(paddedLines, sep)
}

func calcPopupHeight(text string) (rows int) {
	sc := bufio.NewScanner(strings.NewReader(text))
	for sc.Scan() {
		rows++
	}
	return rows + verticalPopupPadding
}
