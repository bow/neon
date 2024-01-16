// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"
	"time"

	"github.com/bow/neon/internal/entity"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const iconAllRead = "✔"

type statusBar struct {
	tview.Flex

	theme *Theme

	eventsWidget     *eventsTextView
	readStatusWidget *tview.TextView
	lastPullWidget   *tview.TextView
}

func newStatusBar(theme *Theme) *statusBar {

	var (
		readStatusWidget = tview.NewTextView().SetTextAlign(tview.AlignCenter)
		lastPullWidget   = tview.NewTextView().SetTextAlign(tview.AlignRight)
	)
	eventsWidget := newEventsTextView(theme)
	eventsWidget.SetTextAlign(tview.AlignLeft)

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

func (b *statusBar) setChangedFunc(f func()) {
	b.eventsWidget.SetChangedFunc(f)
	b.readStatusWidget.SetChangedFunc(f)
	b.lastPullWidget.SetChangedFunc(f)
}

func (b *statusBar) setStats(stats *entity.Stats) {
	if stats.NumFeeds < 1 {
		return
	}
	b.setLastPullTime(stats.LastPullTime)
	if stats.NumEntriesUnread == 0 {
		b.setAllRead()
	}
}

func (b *statusBar) refreshColors() {
	b.eventsWidget.refreshColors()
	b.readStatusWidget.SetTextColor(b.theme.statusBarFG)
	b.lastPullWidget.SetTextColor(b.theme.statusBarFG)
}

func (b *statusBar) setAllRead() {
	b.readStatusWidget.SetText(iconAllRead)
}

func (b *statusBar) setLastPullTime(value *time.Time) {
	if value != nil {
		b.lastPullWidget.SetText(value.Local().Format(shortDateFormat))
	}
}

func (b *statusBar) showEvent(ev *event) {
	b.eventsWidget.show(ev)
}

func (b *statusBar) clearLatestEvent() {
	b.eventsWidget.Clear()
}

type eventsTextView struct {
	tview.TextView
	theme   *Theme
	current *event
}

func newEventsTextView(theme *Theme) *eventsTextView {
	etv := eventsTextView{TextView: *tview.NewTextView(), theme: theme}
	return &etv
}

func (etv *eventsTextView) show(ev *event) {
	etv.current = ev
	etv.refreshColors()
	etv.Clear()
	fmt.Fprintf(&etv.TextView, "%s\n", ev.text)
}

func (etv *eventsTextView) refreshColors() {
	if etv.current == nil {
		return
	}
	var color tcell.Color
	switch etv.current.level {
	case eventLevelInfo:
		color = etv.theme.eventInfoFG
	case eventLevelWarn:
		color = etv.theme.eventWarnFG
	case eventLevelErr:
		color = etv.theme.eventErrFG
	default:
		panic(fmt.Sprintf("unsupported event level: %v", etv.current.level))
	}
	etv.SetTextColor(color)
}
