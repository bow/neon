// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	FeedsPaneTitle    string
	EntriesPaneTitle  string
	ReadingPaneTitle  string
	HelpPopupTitle    string
	StatsPopupTitle   string
	AboutPopupTitle   string
	WelcomePopupTitle string

	Background tcell.Color

	BorderForeground       tcell.Color
	BorderForegroundNormal tcell.Color
	BorderForegroundDim    tcell.Color

	TitleForeground       tcell.Color
	TitleForegroundNormal tcell.Color
	TitleForegroundDim    tcell.Color

	FeedsGroup       tcell.Color
	FeedsGroupNormal tcell.Color
	FeedsGroupDim    tcell.Color

	LastPullForeground       tcell.Color
	LastPullForegroundNormal tcell.Color
	LastPullForegroundDim    tcell.Color

	ActivityNormalForeground       tcell.Color
	ActivityNormalForegroundNormal tcell.Color
	ActivityNormalForegroundDim    tcell.Color

	PopupTitleForeground  tcell.Color
	PopupBorderForeground tcell.Color

	WideViewMinWidth int
}

func (theme *Theme) Dim() {
	theme.BorderForeground = theme.BorderForegroundDim
	theme.TitleForeground = theme.TitleForegroundDim
	theme.LastPullForeground = theme.LastPullForegroundDim
	theme.ActivityNormalForeground = theme.ActivityNormalForegroundDim
	theme.FeedsGroup = theme.FeedsGroupDim
}

func (theme *Theme) Normalize() {
	theme.BorderForeground = theme.BorderForegroundNormal
	theme.TitleForeground = theme.TitleForegroundNormal
	theme.LastPullForeground = theme.LastPullForegroundNormal
	theme.ActivityNormalForeground = theme.ActivityNormalForegroundNormal
	theme.FeedsGroup = theme.FeedsGroupNormal
}

func (theme *Theme) lineStyle() tcell.Style {
	return tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.BorderForeground)
}

const darkForegroundDim = tcell.ColorDimGray

var DarkTheme = &Theme{
	FeedsPaneTitle:    "Feeds",
	EntriesPaneTitle:  "Entries",
	ReadingPaneTitle:  "",
	HelpPopupTitle:    "Keys",
	StatsPopupTitle:   "Stats",
	AboutPopupTitle:   "About",
	WelcomePopupTitle: "Welcome",

	Background: tcell.ColorBlack,

	BorderForeground:       tcell.ColorWhite,
	BorderForegroundNormal: tcell.ColorWhite,
	BorderForegroundDim:    darkForegroundDim,

	TitleForeground:       tcell.ColorYellow,
	TitleForegroundNormal: tcell.ColorYellow,
	TitleForegroundDim:    darkForegroundDim,

	FeedsGroup:       tcell.ColorGrey,
	FeedsGroupNormal: tcell.ColorGrey,
	FeedsGroupDim:    darkForegroundDim,

	LastPullForeground:       tcell.ColorGray,
	LastPullForegroundNormal: tcell.ColorGray,
	LastPullForegroundDim:    darkForegroundDim,

	ActivityNormalForeground:       tcell.ColorDarkGoldenrod,
	ActivityNormalForegroundNormal: tcell.ColorDarkGoldenrod,
	ActivityNormalForegroundDim:    darkForegroundDim,

	PopupBorderForeground: tcell.ColorGray,
	PopupTitleForeground:  tcell.ColorAqua,

	WideViewMinWidth: 150,
}

func init() {
	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical
	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight
}
