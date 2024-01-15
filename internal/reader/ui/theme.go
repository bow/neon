// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Theme struct {
	bg tcell.Color

	lineFG       tcell.Color
	lineNormalFG tcell.Color
	lineDimFG    tcell.Color

	titleFG       tcell.Color
	titleNormalFG tcell.Color
	titleDimFG    tcell.Color

	feedNode       tcell.Color
	feedNodeNormal tcell.Color
	feedNodeDim    tcell.Color

	feedGroupNode       tcell.Color
	feedGroupNodeNormal tcell.Color
	feedGroupNodeDim    tcell.Color

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

	wideViewMinWidth int
}

// nolint:unused
func (t *Theme) dim() {
	t.lineFG = t.lineDimFG
	t.titleFG = t.titleDimFG

	t.statusBarFG = t.statusBarDimFG
	t.eventInfoFG = t.eventInfoDimFG
	t.eventWarnFG = t.eventWarnDimFG
	t.eventErrFG = t.eventErrDimFG

	t.feedNode = t.feedNodeDim
	t.feedGroupNode = t.feedGroupNodeDim
}

// nolint:unused
func (t *Theme) normalize() {
	t.lineFG = t.lineNormalFG
	t.titleFG = t.titleNormalFG

	t.statusBarFG = t.statusBarNormalFG
	t.eventInfoFG = t.eventInfoNormalFG
	t.eventWarnFG = t.eventWarnNormalFG
	t.eventErrFG = t.eventErrNormalFG

	t.feedNode = t.feedNodeNormal
	t.feedGroupNode = t.feedGroupNodeNormal
}

//nolint:unused
func (t *Theme) lineStyle() tcell.Style {
	return tcell.StyleDefault.
		Background(t.bg).
		Foreground(t.lineFG)
}

func loadTheme(name string) (*Theme, error) {
	if name == "dark" {
		return DarkTheme, nil
	}
	return nil, fmt.Errorf("theme %q does not exist", name)
}

const darkForegroundDim = tcell.ColorDimGray

var DarkTheme = &Theme{
	bg: tcell.ColorBlack,

	lineFG:       tcell.ColorWhite,
	lineNormalFG: tcell.ColorWhite,
	lineDimFG:    darkForegroundDim,

	titleFG:       tcell.ColorYellowGreen,
	titleNormalFG: tcell.ColorYellowGreen,
	titleDimFG:    darkForegroundDim,

	feedNode:       tcell.ColorWhite,
	feedNodeNormal: tcell.ColorWhite,
	feedNodeDim:    darkForegroundDim,

	feedGroupNode:       tcell.ColorGrey,
	feedGroupNodeNormal: tcell.ColorGrey,
	feedGroupNodeDim:    darkForegroundDim,

	statusBarFG:       tcell.ColorGray,
	statusBarNormalFG: tcell.ColorGray,
	statusBarDimFG:    darkForegroundDim,

	eventInfoFG:       tcell.ColorYellowGreen,
	eventInfoNormalFG: tcell.ColorYellowGreen,
	eventInfoDimFG:    darkForegroundDim,

	eventWarnFG:       tcell.ColorDarkGoldenrod,
	eventWarnNormalFG: tcell.ColorDarkGoldenrod,
	eventWarnDimFG:    darkForegroundDim,

	eventErrFG:       tcell.ColorTomato,
	eventErrNormalFG: tcell.ColorTomato,
	eventErrDimFG:    darkForegroundDim,

	popupBorderFG: tcell.ColorGray,
	popupTitleFG:  tcell.ColorAqua,

	wideViewMinWidth: 150,
}

func init() {
	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical
	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight
}
