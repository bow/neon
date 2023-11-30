// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/iris/internal"
)

const (
	mainPageName    = "main"
	helpPageName    = "help"
	statsPageName   = "stats"
	versionPageName = "version"

	pulledIcon = "▼"

	shortDateFormat = "02/Jan/06 15:04"
	longDateFormat  = "2 January 2006 - 15:04:05 MST"
)

type Reader struct {
	ctx   context.Context
	store internal.FeedStore
	theme *Theme

	app *tview.Application

	root        *tview.Pages
	mainPage    *tview.Grid
	helpPage    *tview.Grid
	statsPage   *tview.Grid
	versionPage *tview.Grid

	feedsPane   *tview.Box
	entriesPane *tview.Box
	readingPane *tview.Box

	statusWidget       *tview.TextView
	lastPullTextWidget *tview.TextView
	lastPullIconWidget *tview.TextView

	statsCache *internal.Stats
}

func NewReader(ctx context.Context, store internal.FeedStore, theme *Theme) *Reader {

	if theme == nil {
		theme = DefaultTheme
	}

	reader := Reader{
		ctx:   ctx,
		store: store,
		theme: theme,
		root:  tview.NewPages(),
		app:   tview.NewApplication(),
	}

	reader.setupMainPage()
	reader.setupHelpPage()
	reader.setupStatsPage()
	reader.setupVersionPage()

	reader.root.
		AddAndSwitchToPage(mainPageName, reader.mainPage, true).
		AddPage(helpPageName, reader.helpPage, true, false).
		AddPage(statsPageName, reader.statsPage, true, false).
		AddPage(versionPageName, reader.versionPage, true, false)

	reader.root.SetInputCapture(reader.keyHandler())
	reader.app.SetRoot(reader.root, true).EnableMouse(true)

	return &reader
}

func (r *Reader) Show() error {
	stats, err := r.store.GetGlobalStats(r.ctx)
	if err != nil {
		return err
	}
	r.setLastPullTime(stats.LastPullTime)
	r.setLastPullIcon()

	return r.app.Run()
}

func (r *Reader) setupMainPage() {

	feedsPane := r.newPane(r.theme.FeedsPaneTitle, false)
	entriesPane := r.newPane(r.theme.EntriesPaneTitle, false)
	readingPane := r.newPane(r.theme.ReadingPaneTitle, true)

	narrowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(feedsPane, 0, 3, false).
		AddItem(entriesPane, 0, 4, false).
		AddItem(readingPane, 0, 5, false).
		AddItem(r.newNarrowFooterBorder(), 1, 0, false)

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
		AddItem(r.newWideFooterBorder(45), 1, 0, false)

	statusWidget := tview.NewTextView().SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() { r.app.Draw() })

	lastPullIconWidget := tview.NewTextView().SetTextColor(r.theme.LastPullForeground).
		SetTextAlign(tview.AlignCenter)

	lastPullTextWidget := tview.NewTextView().SetTextColor(r.theme.LastPullForeground).
		SetTextAlign(tview.AlignRight).
		SetChangedFunc(func() { r.app.Draw() })

	lastPullFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(lastPullIconWidget, 1, 0, false).
		AddItem(lastPullTextWidget, len(shortDateFormat)+1, 0, true)

	footer := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statusWidget, 0, 1, false).
		AddItem(lastPullFlex, len(shortDateFormat)+2, 1, false)

	mainPage := tview.NewGrid().
		SetColumns(0).
		SetRows(0, 1).
		SetBorders(false).
		AddItem(narrowFlex, 0, 0, 1, 1, 0, 0, false).
		AddItem(wideFlex, 0, 0, 1, 1, 0, r.theme.WideViewMinWidth, false).
		AddItem(footer, 1, 0, 1, 1, 0, 0, false)

	r.feedsPane = feedsPane
	r.entriesPane = entriesPane
	r.readingPane = readingPane

	r.statusWidget = statusWidget
	r.lastPullTextWidget = lastPullTextWidget
	r.lastPullIconWidget = lastPullIconWidget

	r.mainPage = mainPage
}

func (r *Reader) setupHelpPage() {

	helpWidget := tview.NewTextView().
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
[yellow]1,F[-]     : Toggle feeds pane focus
[yellow]2,E[-]     : Toggle entries pane focus
[yellow]3,R[-]     : Toggle reading pane focus
[yellow]Tab[-]     : Switch to next pane
[yellow]Alt-Tab[-] : Switch to previous pane
[yellow]X[-]       : Export feeds to OPML
[yellow]I[-]       : Import feeds from OPML
[yellow]Esc[-]     : Unset current focus or close open frame
[yellow]S[-]       : Toggle stats info
[yellow]V[-]       : Toggle version info
[yellow]h,?[-]     : Toggle this help
[yellow]q,Ctrl-C[-]: Quit reader`)

	helpFrame := tview.NewFrame(helpWidget).SetBorders(1, 1, 0, 0, 2, 2)

	helpFrame.SetBorder(true).
		SetBorderColor(r.theme.PopupBorderForeground).
		SetTitle(r.makeTitle(r.theme.HelpPopupTitle)).
		SetTitleColor(r.theme.PopupTitleForeground)

	helpPage := tview.NewGrid().
		SetColumns(0, 55, 0).
		SetRows(0, 38, 0).
		AddItem(helpFrame, 1, 1, 1, 1, 0, 0, true)

	r.helpPage = helpPage
}

func (r *Reader) setupStatsPage() {

	if r.statsCache == nil {
		stats, err := r.store.GetGlobalStats(r.ctx)
		if err != nil {
			panic(err)
		}
		r.statsCache = stats
	}

	statsWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(
			fmt.Sprintf(`[aqua]Feeds[-]
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
			),
		)

	statsFrame := tview.NewFrame(statsWidget).SetBorders(1, 1, 0, 0, 2, 2)

	statsFrame.SetBorder(true).
		SetBorderColor(r.theme.PopupBorderForeground).
		SetTitle(r.makeTitle("Stats")).
		SetTitleColor(r.theme.PopupTitleForeground)

	statsPage := tview.NewGrid().
		SetColumns(0, 37, 0).
		SetRows(-1, 13, -3).
		AddItem(statsFrame, 1, 1, 1, 1, 0, 0, true)

	r.statsPage = statsPage
}

func (r *Reader) setupVersionPage() {

	commit := internal.GitCommit()

	var buildTime = internal.BuildTime()
	buildTimeVal, err := time.Parse(time.RFC3339, buildTime)
	if err == nil {
		buildTime = buildTimeVal.Format(longDateFormat)
	}

	versionWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf(`[yellow]Version[-]   : %s
[yellow]Git commit[-]: %s
[yellow]Build time[-]: %s
`,
			internal.Version(),
			commit,
			buildTime,
		))

	versionFrame := tview.NewFrame(versionWidget).SetBorders(1, 1, 0, 0, 2, 2)

	versionFrame.SetBorder(true).
		SetBorderColor(r.theme.PopupBorderForeground).
		SetTitle(r.makeTitle("iris feed reader")).
		SetTitleColor(r.theme.PopupTitleForeground)

	versionPage := tview.NewGrid().
		SetColumns(0, len(commit)+len("Git commit:")+7, 0).
		SetRows(-1, 7, -3).
		AddItem(versionFrame, 1, 1, 1, 1, 0, 0, true)

	r.versionPage = versionPage
}

// nolint:revive,exhaustive
func (r *Reader) keyHandler() func(event *tcell.EventKey) *tcell.EventKey {

	defaultHandler := func(
		event *tcell.EventKey,
		key tcell.Key,
		keyr rune,
		front string,
		focused tview.Primitive,
	) *tcell.EventKey {

		switch key {

		case tcell.KeyRune:
			switch keyr {
			case '1', '2', '3', 'F', 'E', 'R':
				if front != mainPageName {
					r.root.HidePage(front)
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
				} else {
					if front != mainPageName {
						r.root.HidePage(front)
					}
					r.root.ShowPage(statsPageName)
				}
				return nil

			case 'V':
				if front == versionPageName {
					r.root.HidePage(front)
				} else {
					if front != mainPageName {
						r.root.HidePage(front)
					}
					r.root.ShowPage(versionPageName)
				}
				return nil

			case 'h', '?':
				if front == helpPageName {
					r.root.HidePage(front)
				} else {
					if front != mainPageName {
						r.root.HidePage(front)
					}
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
			case helpPageName, statsPageName, versionPageName:
				r.root.HidePage(front)
			case mainPageName, "":
				r.app.SetFocus(r.root)
			}
			return nil
		}

		return event
	}

	return func(event *tcell.EventKey) *tcell.EventKey {
		var (
			focused  = r.app.GetFocus()
			front, _ = r.root.GetFrontPage()
			key      = event.Key()
			keyr     = event.Rune()
		)
		switch focused {

		case r.feedsPane:
			switch keyr {
			case 'P':

				// TODO: Consider wrapping footer as own primitive.
				r.setNormalStatus("Pulling feeds")
				go func() {
					var count int
					ch := r.store.PullFeeds(r.ctx, []internal.ID{})
					for pr := range ch {
						if err := pr.Error(); err != nil {
							// TODO: Add ok / fail status in ...?
							panic(err)
						}
						r.setNormalStatus("Pulling: %s done", pr.URL())
						count++
					}
					if count > 1 {
						r.setNormalStatus("Pulled %d feeds successfully", count)
					} else {
						r.setNormalStatus("Pulled %d feed successfully", count)
					}

					stats, err := r.store.GetGlobalStats(r.ctx)
					if err != nil {
						panic(err)
					}
					r.statsCache = stats
					r.setLastPullTime(stats.LastPullTime)
				}()
				return nil

			default:
				return defaultHandler(event, key, keyr, front, focused)
			}

		default:
			return defaultHandler(event, key, keyr, front, focused)
		}
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

func (r *Reader) setNormalStatus(text string, a ...any) {
	r.statusWidget.
		SetTextColor(r.theme.StatusNormalForeground).
		Clear()
	if len(a) > 0 {
		fmt.Fprintf(r.statusWidget, "%s\n", fmt.Sprintf(text, a...))
	} else {
		fmt.Fprintf(r.statusWidget, "%s\n", text)
	}
}

func (r *Reader) setLastPullIcon() {
	r.lastPullIconWidget.SetText(pulledIcon)
}

func (r *Reader) setLastPullTime(value *time.Time) {
	r.lastPullTextWidget.SetText(value.Local().Format(shortDateFormat))
}

func (r *Reader) getAdjacentFocusTarget(
	current tview.Primitive,
	reverse bool,
) tview.Primitive {
	targets := []*tview.Box{r.feedsPane, r.entriesPane, r.readingPane}
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

func (r *Reader) newPane(title string, addTopLeftBorderTip bool) *tview.Box {

	var (
		unfocused, focused string
		lineStyle          = r.theme.lineStyle()
	)

	if title != "" {
		unfocused = r.makeTitle(title)
		focused = r.makeTitle(fmt.Sprintf("• %s", title))
	} else {
		focused = r.makeTitle("•")
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

			return x + 1, y + 1, width - 2, height - 1
		}
	}

	box := tview.NewBox().SetDrawFunc(makedrawf(unfocused, 3))

	box.SetFocusFunc(func() { box.SetDrawFunc(makedrawf(focused, 1)) })
	box.SetBlurFunc(func() { box.SetDrawFunc(makedrawf(unfocused, 3)) })

	return box
}

func (r *Reader) newPaneDivider() *tview.Box {

	style := r.theme.lineStyle()
	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		screen.SetContent(x, y, tview.BoxDrawingsLightDownAndHorizontal, nil, style)
		for cy := y + 1; cy < y+height; cy++ {
			screen.SetContent(x, cy, tview.BoxDrawingsLightVertical, nil, style)
		}
		return x + 1, y + 1, width - 2, height - 1
	}

	return tview.NewBox().SetBorder(false).SetDrawFunc(drawf)
}

func (r *Reader) newNarrowFooterBorder() *tview.Box {

	style := r.theme.lineStyle()
	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
		}
		return x + 1, y + 1, width - 2, height - 1
	}

	return tview.NewBox().SetBorder(false).SetDrawFunc(drawf)
}

func (r *Reader) newWideFooterBorder(branchPoint int) *tview.Box {

	style := r.theme.lineStyle()

	drawf := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

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

type Theme struct {
	FeedsPaneTitle   string
	EntriesPaneTitle string
	ReadingPaneTitle string
	HelpPopupTitle   string

	Background             tcell.Color
	BorderForeground       tcell.Color
	TitleForeground        tcell.Color
	VersionForeground      tcell.Color
	LastPullForeground     tcell.Color
	StatusNormalForeground tcell.Color
	PopupTitleForeground   tcell.Color
	PopupBorderForeground  tcell.Color

	WideViewMinWidth int
}

func (theme *Theme) lineStyle() tcell.Style {
	return tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.BorderForeground)
}

var DefaultTheme = &Theme{
	FeedsPaneTitle:   "Feeds",
	EntriesPaneTitle: "Entries",
	ReadingPaneTitle: "",
	HelpPopupTitle:   "Keys",

	Background:             tcell.ColorBlack,
	BorderForeground:       tcell.ColorWhite,
	TitleForeground:        tcell.ColorBlue,
	VersionForeground:      tcell.ColorGray,
	LastPullForeground:     tcell.ColorGray,
	StatusNormalForeground: tcell.ColorDarkGoldenrod,
	PopupBorderForeground:  tcell.ColorGray,
	PopupTitleForeground:   tcell.ColorAqua,

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
