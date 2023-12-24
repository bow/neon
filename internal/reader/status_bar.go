// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"fmt"
	"sync"
	"time"

	"github.com/rivo/tview"

	"github.com/bow/lens/internal"
)

const iconAllRead = "✔"

type statusBar struct {
	sync.RWMutex

	theme     *Theme
	container *tview.Flex

	eventsWidget   *tview.TextView
	readWidget     *tview.TextView
	lastPullWidget *tview.TextView

	events []*event

	visible bool
}

func newStatusBar(theme *Theme) *statusBar {

	eventsWidget := tview.NewTextView().SetTextAlign(tview.AlignLeft)

	readStatusWidget := tview.NewTextView().SetTextColor(theme.LastPullForeground).
		SetTextAlign(tview.AlignCenter)

	lastPullWidget := tview.NewTextView().SetTextColor(theme.LastPullForeground).
		SetTextAlign(tview.AlignRight)

	quickStatusWidget := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(readStatusWidget, 1, 0, false).
		AddItem(lastPullWidget, len(shortDateFormat)+1, 0, true)

	container := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(eventsWidget, 0, 1, false).
		AddItem(quickStatusWidget, len(shortDateFormat)+2, 1, false)

	bar := statusBar{
		theme:          theme,
		container:      container,
		eventsWidget:   eventsWidget,
		readWidget:     readStatusWidget,
		lastPullWidget: lastPullWidget,
		visible:        true,
	}

	return &bar
}

func (b *statusBar) refreshColors() {
	b.eventsWidget.SetTextColor(b.theme.EventNormalForeground)
	b.readWidget.SetTextColor(b.theme.LastPullForeground)
	b.lastPullWidget.SetTextColor(b.theme.LastPullForeground)
}

func (b *statusBar) setChangedFunc(handler func()) *statusBar {
	b.eventsWidget.SetChangedFunc(handler)
	b.readWidget.SetChangedFunc(handler)
	b.lastPullWidget.SetChangedFunc(handler)
	return b
}

func (b *statusBar) updateFromStats(stats *internal.Stats) {
	if stats.NumFeeds < 1 {
		return
	}
	b.setLastPullTime(stats.LastPullTime)
	if stats.NumEntriesUnread == 0 {
		b.setAllRead()
	}
}

func (b *statusBar) setAllRead() {
	b.readWidget.SetText(iconAllRead)
}

func (b *statusBar) setLastPullTime(value *time.Time) {
	if value != nil {
		b.lastPullWidget.SetText(value.Local().Format(shortDateFormat))
	}
}

func (b *statusBar) toggleFromMainPage(page *tview.Grid) {
	if b.visible {
		b.removeFromMainPage(page)
	} else {
		b.addToMainPage(page)
	}
	b.visible = !b.visible
}

func (b *statusBar) addToMainPage(page *tview.Grid) *statusBar {
	page.SetRows(0, 1).AddItem(b.container, 1, 0, 1, 1, 0, 0, false)
	return b
}

func (b *statusBar) removeFromMainPage(page *tview.Grid) *statusBar {
	page.RemoveItem(b.container).SetRows(0)
	return b
}

func (b *statusBar) showNormalEvent(text string, a ...any) {
	b.eventsWidget.
		SetTextColor(b.theme.EventNormalForeground).
		Clear()
	if len(a) > 0 {
		fmt.Fprintf(b.eventsWidget, "%s\n", fmt.Sprintf(text, a...))
	} else {
		fmt.Fprintf(b.eventsWidget, "%s\n", text)
	}
}

func (b *statusBar) clearEvents() {
	b.eventsWidget.Clear()
}
