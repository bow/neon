// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package oreader

import (
	context "context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/neon/internal/entity"
)

const iconAllRead = "✔"

type statusBar struct {
	ctx context.Context

	theme     *Theme
	container *tview.Flex

	latestEventWidget *tview.TextView
	readWidget        *tview.TextView
	lastPullWidget    *tview.TextView

	visible bool

	eventsCh chan *event
	events   []*event
}

func newStatusBar(ctx context.Context, theme *Theme) *statusBar {

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
		ctx:               ctx,
		theme:             theme,
		container:         container,
		latestEventWidget: eventsWidget,
		readWidget:        readStatusWidget,
		lastPullWidget:    lastPullWidget,

		visible: true,

		events:   make([]*event, 0),
		eventsCh: make(chan *event),
	}

	return &bar
}

func (b *statusBar) refreshColors() {
	b.latestEventWidget.SetTextColor(b.theme.EventInfoForeground)
	b.readWidget.SetTextColor(b.theme.LastPullForeground)
	b.lastPullWidget.SetTextColor(b.theme.LastPullForeground)
}

func (b *statusBar) setChangedFunc(handler func()) *statusBar {
	b.latestEventWidget.SetChangedFunc(handler)
	b.readWidget.SetChangedFunc(handler)
	b.lastPullWidget.SetChangedFunc(handler)
	return b
}

func (b *statusBar) updateFromStats(stats *entity.Stats) {
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

func (b *statusBar) startEventPoll() (stop func()) {
	done := make(chan struct{})
	stop = func() {
		defer close(done)
		done <- struct{}{}
	}

	go func() {
		defer close(b.eventsCh)
		for {
			select {
			case <-done:
				return
			case ev := <-b.eventsCh:
				b.displayEvent(ev)
				b.events = append(b.events, ev)
			}
		}
	}()

	return stop
}

func (b *statusBar) displayEvent(ev *event) {
	var color tcell.Color
	switch ev.level {
	case eventLevelInfo:
		color = b.theme.EventInfoForeground
	case eventLevelWarn:
		color = b.theme.EventWarnForeground
	case eventLevelErr:
		color = b.theme.EventErrForeground
	default:
		panic(fmt.Sprintf("unsupported event level: %v", ev.level))
	}
	b.latestEventWidget.SetTextColor(color).Clear()
	fmt.Fprintf(b.latestEventWidget, "%s\n", ev.text)
}

func (b *statusBar) clearLatestEvent() {
	b.latestEventWidget.Clear()
}

func (b *statusBar) infoEventf(text string, a ...any) { b.eventf(eventLevelInfo, text, a...) }
func (b *statusBar) warnEventf(text string, a ...any) { b.eventf(eventLevelWarn, text, a...) }
func (b *statusBar) errEventf(text string, a ...any) {
	b.eventf(eventLevelErr, fmt.Sprintf(text, a...))
}
func (b *statusBar) errEvent(err error) {
	b.eventf(eventLevelErr, fmt.Sprintf("%s", err))
}

func (b *statusBar) eventf(level eventLevel, text string, a ...any) {
	ev := event{level: level, timestamp: time.Now(), text: fmt.Sprintf(text, a...)}
	go func() { b.eventsCh <- &ev }()
}

type event struct {
	level     eventLevel
	timestamp time.Time
	text      string
}

type eventLevel uint8

const (
	eventLevelInfo eventLevel = iota
	eventLevelWarn
	eventLevelErr
)
