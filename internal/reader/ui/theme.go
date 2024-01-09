// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func LoadTheme(name string) (*Theme, error) {
	if name == "dark" {
		return DarkTheme, nil
	}
	return nil, fmt.Errorf("theme %q does not exist", name)
}

type Theme struct {
	bg tcell.Color

	borderFG       tcell.Color
	borderNormalFG tcell.Color
	borderDimFG    tcell.Color

	titleFG       tcell.Color
	titleNormalFG tcell.Color
	titleDimFG    tcell.Color

	feedsGroup       tcell.Color
	feedsGroupNormal tcell.Color
	feedsGroupDim    tcell.Color

	statusBarFG       tcell.Color
	statusBarNormalFG tcell.Color
	statusBarDimFG    tcell.Color

	eventInfoFG       tcell.Color
	eventInfoNormalFG tcell.Color
	eventInfoDimFG    tcell.Color

	eventWarnFG       tcell.Color
	eventWarnNormalFG tcell.Color
	eventWarnDimFG    tcell.Color

	eventErrFG       tcell.Color
	eventErrNormalFG tcell.Color
	eventErrDimFG    tcell.Color

	popupTitleFG  tcell.Color
	popupBorderFG tcell.Color

	WideViewMinWidth int
}

// nolint:unused
func (theme *Theme) dim() {
	theme.borderFG = theme.borderDimFG
	theme.titleFG = theme.titleDimFG
	theme.statusBarFG = theme.statusBarDimFG
	theme.eventInfoFG = theme.eventInfoDimFG
	theme.feedsGroup = theme.feedsGroupDim
}

// nolint:unused
func (theme *Theme) normalize() {
	theme.borderFG = theme.borderNormalFG
	theme.titleFG = theme.titleNormalFG
	theme.statusBarFG = theme.statusBarNormalFG
	theme.eventInfoFG = theme.eventInfoNormalFG
	theme.feedsGroup = theme.feedsGroupNormal
}

//nolint:unused
func (theme *Theme) lineStyle() tcell.Style {
	return tcell.StyleDefault.
		Background(theme.bg).
		Foreground(theme.borderFG)
}

const darkForegroundDim = tcell.ColorDimGray

var DarkTheme = &Theme{
	bg: tcell.ColorBlack,

	borderFG:       tcell.ColorWhite,
	borderNormalFG: tcell.ColorWhite,
	borderDimFG:    darkForegroundDim,

	titleFG:       tcell.ColorYellow,
	titleNormalFG: tcell.ColorYellow,
	titleDimFG:    darkForegroundDim,

	feedsGroup:       tcell.ColorGrey,
	feedsGroupNormal: tcell.ColorGrey,
	feedsGroupDim:    darkForegroundDim,

	statusBarFG:       tcell.ColorGray,
	statusBarNormalFG: tcell.ColorGray,
	statusBarDimFG:    darkForegroundDim,

	eventInfoFG:       tcell.ColorOliveDrab,
	eventInfoNormalFG: tcell.ColorOliveDrab,
	eventInfoDimFG:    darkForegroundDim,

	eventWarnFG:       tcell.ColorDarkGoldenrod,
	eventWarnNormalFG: tcell.ColorDarkGoldenrod,
	eventWarnDimFG:    darkForegroundDim,

	eventErrFG:       tcell.ColorTomato,
	eventErrNormalFG: tcell.ColorTomato,
	eventErrDimFG:    darkForegroundDim,

	popupBorderFG: tcell.ColorGray,
	popupTitleFG:  tcell.ColorAqua,

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
