// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package tui

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
	VersionPopupTitle string
	WelcomePopupTitle string

	Background             tcell.Color
	BorderForeground       tcell.Color
	TitleForeground        tcell.Color
	VersionForeground      tcell.Color
	LastPullForeground     tcell.Color
	StatusNormalForeground tcell.Color
	PopupTitleForeground   tcell.Color
	PopupBorderForeground  tcell.Color

	NormalBorderForeground tcell.Color
	NormalTitleForeground  tcell.Color

	DimBorderForeground tcell.Color
	DimTitleForeground  tcell.Color

	WideViewMinWidth int
}

func (theme *Theme) Dim() {
	theme.BorderForeground = theme.DimBorderForeground
	theme.TitleForeground = theme.DimTitleForeground
}

func (theme *Theme) Normalize() {
	theme.BorderForeground = theme.NormalBorderForeground
	theme.TitleForeground = theme.NormalTitleForeground
}

func (theme *Theme) lineStyle() tcell.Style {
	return tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.BorderForeground)
}

var DarkTheme = &Theme{
	FeedsPaneTitle:    "Feeds",
	EntriesPaneTitle:  "Entries",
	ReadingPaneTitle:  "",
	HelpPopupTitle:    "Keys",
	StatsPopupTitle:   "Stats",
	VersionPopupTitle: "About",
	WelcomePopupTitle: "Welcome",

	Background:             tcell.ColorBlack,
	BorderForeground:       tcell.ColorWhite,
	TitleForeground:        tcell.ColorBlue,
	VersionForeground:      tcell.ColorGray,
	LastPullForeground:     tcell.ColorGray,
	StatusNormalForeground: tcell.ColorDarkGoldenrod,
	PopupBorderForeground:  tcell.ColorGray,
	PopupTitleForeground:   tcell.ColorAqua,

	// TODO: Add New() to ensure values are equal.
	NormalBorderForeground: tcell.ColorWhite,
	NormalTitleForeground:  tcell.ColorBlue,

	DimBorderForeground: tcell.ColorDimGray,
	DimTitleForeground:  tcell.ColorDimGray,

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
