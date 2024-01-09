// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"github.com/rivo/tview"
)

type statusBar struct {
	tview.Flex

	theme *Theme

	eventsWidget     *tview.TextView
	readStatusWidget *tview.TextView
	lastPullWidget   *tview.TextView
}

func newStatusBar(theme *Theme) *statusBar {

	var (
		eventsWidget     = tview.NewTextView().SetTextAlign(tview.AlignLeft)
		readStatusWidget = tview.NewTextView().SetTextAlign(tview.AlignCenter)
		lastPullWidget   = tview.NewTextView().SetTextAlign(tview.AlignRight)
	)

	quickStatusFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(readStatusWidget, 1, 0, false).
		AddItem(lastPullWidget, len(shortDateFormat)+1, 0, true)

	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	bar := statusBar{
		Flex:             *flex,
		theme:            theme,
		eventsWidget:     eventsWidget,
		readStatusWidget: readStatusWidget,
		lastPullWidget:   lastPullWidget,
	}
	bar.AddItem(eventsWidget, 0, 1, false).
		AddItem(quickStatusFlex, len(shortDateFormat)+2, 1, false)
	bar.refreshColors()

	return &bar
}

func (b *statusBar) refreshColors() {
	b.readStatusWidget.SetTextColor(b.theme.statusBarFG)
	b.lastPullWidget.SetTextColor(b.theme.statusBarFG)
}
