// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package oreader

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/grpc"

	"github.com/bow/neon/api"
	"github.com/bow/neon/internal"
	"github.com/bow/neon/internal/entity"
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
	ctx    context.Context
	addr   string
	client api.NeonClient

	screen   tcell.Screen
	theme    *Theme
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

	feedsCh    chan *entity.Feed
	statsCache *entity.Stats
	focusStack tview.Primitive
}

type drawFunc func(screen tcell.Screen, x int, y int, w int, h int) (ix int, iy int, iw int, ih int)

type Builder struct {
	ctx      context.Context
	addr     string
	dopts    []grpc.DialOption
	initPath string
	theme    *Theme

	// Only for testing.
	cl  api.NeonClient
	scr tcell.Screen
}

func NewBuilder() *Builder {
	b := Builder{
		dopts: nil,
		ctx:   context.Background(),
		theme: DarkTheme,
	}
	return &b
}

func (b *Builder) Address(addr string) *Builder {
	b.addr = addr
	return b
}

func (b *Builder) DialOpts(dialOpts ...grpc.DialOption) *Builder {
	b.dopts = dialOpts
	return b
}

func (b *Builder) Context(ctx context.Context) *Builder {
	b.ctx = ctx
	return b
}

func (b *Builder) InitPath(path string) *Builder {
	b.initPath = path
	return b
}

func (b *Builder) Theme(theme *Theme) *Builder {
	b.theme = theme
	return b
}

func (b *Builder) client(cl api.NeonClient) *Builder {
	b.cl = cl
	return b
}

func (b *Builder) screen(scr tcell.Screen) *Builder {
	b.scr = scr
	return b
}

func (b *Builder) Build() (*Reader, error) {

	if b.addr == "" && b.cl == nil {
		return nil, fmt.Errorf("reader server address must be specified")
	}

	var (
		client api.NeonClient
		conn   *grpc.ClientConn
		err    error
	)
	if b.cl != nil {
		client = b.cl
	} else {
		conn, err = grpc.DialContext(b.ctx, b.addr, b.dopts...)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, fmt.Errorf("timeout when connecting to server %q", b.addr)
			}
			return nil, err
		}
		client = api.NewNeonClient(conn)
	}

	var screen tcell.Screen
	if b.scr != nil {
		screen = b.scr
	} else {
		screen, err = tcell.NewScreen()
		if err != nil {
			return nil, err
		}
	}

	rdr := Reader{
		ctx:     b.ctx,
		client:  client,
		addr:    b.addr,
		screen:  screen,
		theme:   b.theme,
		feedsCh: make(chan *entity.Feed),
	}
	rdr.setupLayout()

	return &rdr, nil
}

func (r *Reader) setupLayout() {
	r.app = tview.NewApplication().SetScreen(r.screen)
	r.root = tview.NewPages()

	r.setupMainPage()
	r.setupHelpPage()
	r.setupStatsPage()
	r.setupAboutPage()

	r.root.
		AddAndSwitchToPage(mainPageName, r.mainPage, true).
		AddPage(helpPageName, r.helpPage, true, false).
		AddPage(statsPageName, r.statsPage, true, false).
		AddPage(aboutPageName, r.aboutPage, true, false)

	r.root.SetInputCapture(r.globalKeyHandler())
	r.app.SetRoot(r.root, true).EnableMouse(true)
}

func (r *Reader) Show() error {

	if r.statsCache == nil {
		if err := r.getGlobalStats(); err != nil {
			return err
		}
	}

	if !r.isInitialized() {
		welcomeText := fmt.Sprintf(`Hello and welcome the %s reader.

For help, press [yellow]h[-] or go to [yellow]https://github.com/bow/neon[-].

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

	rsp, err := r.client.ListFeeds(r.ctx, &api.ListFeedsRequest{})
	if err != nil {
		panic(err)
	}
	for _, feed := range rsp.GetFeeds() {
		feed := feed
		go func() { r.feedsCh <- entity.FromFeedPb(feed) }()
	}

	stop := r.bar.startEventPoll()
	defer stop()

	return r.app.Run()
}

func (r *Reader) setupMainPage() {

	feedsPane := newFeedsPane(r.theme, r.feedsCh)
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

	r.bar = newStatusBar(r.ctx, r.theme).
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
		if err := r.getGlobalStats(); err != nil {
			panic(err)
		}
	}

	sc := r.statsCache
	var lpt string
	if sc.LastPullTime != nil {
		lpt = sc.LastPullTime.Format(longDateFormat)
	}

	statsText := fmt.Sprintf(`[aqua]Feeds[-]
[yellow]Total[-]: %d

[aqua]Entries[-]
[yellow]Unread[-]: %d
[yellow]Total[-] : %d

[aqua]Last pulled[-]
%s`,
		sc.NumFeeds,
		sc.NumEntries,
		sc.NumEntriesUnread,
		lpt,
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
[yellow]Build time[-]: %s
[yellow]Server[-]    : %s`,
		centerBanner(internal.Banner(), width),
		internal.Version(),
		commit,
		buildTime,
		r.addr,
	)

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
			case 'F', 'E', 'R':
				if front != mainPageName {
					r.root.HidePage(front)
					r.normalizeColors()
					front = mainPageName
				}
				if front == mainPageName {
					target := r.focusTarget(keyr)
					r.app.SetFocus(target)
					r.stashFocus()
				}
				return nil

			case 'S':
				if front == statsPageName {
					r.hidePopup(front)
				} else if front != welcomePageName {
					r.showPopup(statsPageName, front)
				}
				return nil

			case 'A':
				if front == aboutPageName {
					r.hidePopup(front)
				} else if front != welcomePageName {
					r.showPopup(aboutPageName, front)
				}
				return nil

			case 'b':
				r.bar.toggleFromMainPage(r.mainPage)
				return nil

			case 'c':
				r.bar.clearLatestEvent()
				return nil

			case 'h', '?':
				if front == helpPageName {
					r.hidePopup(front)
				} else {
					r.showPopup(helpPageName, front)
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
				target := r.adjacentFocusTarget(focused, reverse)
				r.app.SetFocus(target)
			}
			return nil

		case tcell.KeyEscape:
			switch front {
			case mainPageName, "":
				r.app.SetFocus(r.root)
			default:
				r.hidePopup(front)
			}
			return nil
		}

		return event
	}
}

// nolint:revive,exhaustive
func (r *Reader) feedsPaneKeyHandler() func(event *tcell.EventKey) *tcell.EventKey {
	pullLock := make(chan struct{}, 1)

	return func(event *tcell.EventKey) *tcell.EventKey {
		keyr := event.Rune()

		if keyr == 'P' {

			go func() {
				select {
				case pullLock <- struct{}{}:
					defer func() { <-pullLock }()
				default:
					return
				}

				r.bar.infoEventf("Pulling feeds")

				var okCount, errCount, totalCount int
				stream, err := r.client.PullFeeds(r.ctx, &api.PullFeedsRequest{})
				if err != nil {
					r.bar.errEvent(err)
					return
				}
				for {
					rsp, serr := stream.Recv()
					if serr == io.EOF {
						break
					}
					if serr != nil {
						r.bar.errEventf("Pull failed for %s: %s", rsp.GetUrl(), serr)
						errCount++
					} else {
						r.bar.infoEventf("Pulled %s", rsp.GetUrl())
						rsp.GetFeed()
						okCount++
					}
					totalCount++
				}
				if errCount == 0 {
					switch okCount {
					case 0:
						r.bar.infoEventf("No feeds to pull")
					case 1:
						r.bar.infoEventf("%d/%d feed pulled successfully", okCount, totalCount)
					default:
						r.bar.infoEventf("%d/%d feeds pulled successfully", okCount, totalCount)
					}
				} else {
					switch okCount {
					case 0:
						r.bar.errEventf("Failed to pull any feeds")
					default:
						r.bar.warnEventf("Only %d/%d feeds pulled successfully", okCount, totalCount)
					}
				}

				if err := r.getGlobalStats(); err != nil {
					r.bar.errEventf(fmt.Sprintf("Failed to refresh stats: %s", err))
				}
			}()
			return nil
		}
		return event
	}
}

func (r *Reader) focusTarget(keyr rune) tview.Primitive {
	var target tview.Primitive
	switch keyr { // nolint:exhaustive
	case 'F':
		target = r.feedsPane
	case 'E':
		target = r.entriesPane
	case 'R':
		target = r.readingPane
	default:
		panic(fmt.Sprintf("unexpected key: %c", keyr))
	}
	return target
}

func (r *Reader) adjacentFocusTarget(
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

func (r *Reader) showPopup(name string, currentFront string) {
	if currentFront == mainPageName {
		r.stashFocus()
	} else {
		r.root.HidePage(currentFront)
	}
	r.dimColors()
	r.root.ShowPage(name)
}

func (r *Reader) hidePopup(name string) {
	r.root.HidePage(name)
	r.normalizeColors()

	if r.focusStack != nil {
		r.app.SetFocus(r.focusStack)
	}
	r.focusStack = nil
}

func (r *Reader) stashFocus() {
	r.focusStack = r.app.GetFocus()
}

func (r *Reader) dimColors() {
	r.theme.Dim()
	r.bar.refreshColors()
	r.feedsPane.refreshColors()
}

func (r *Reader) normalizeColors() {
	r.theme.Normalize()
	r.bar.refreshColors()
	r.feedsPane.refreshColors()
}

func (r *Reader) newPane(title string, addTopLeftBorderTip bool) *tview.Box {

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

func (r *Reader) getGlobalStats() error {
	rsp, err := r.client.GetStats(r.ctx, &api.GetStatsRequest{})
	if err != nil {
		return err
	}

	stats := entity.FromStatsPb(rsp.GetGlobal())
	r.bar.updateFromStats(stats)
	r.statsCache = stats

	return nil
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

	leftPad := strings.Repeat(" ", ((width-maxLineWidth)/2)-(leftPopupMargin*2))
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