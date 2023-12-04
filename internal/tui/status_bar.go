package tui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

const pulledIcon = "▼"

type statusBar struct {
	reader    *Reader
	container *tview.Flex

	statusWidget       *tview.TextView
	lastPullTextWidget *tview.TextView
	lastPullIconWidget *tview.TextView

	visible bool
}

func newStatusBar(r *Reader) *statusBar {

	statusWidget := tview.NewTextView().SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() { r.app.Draw() })

	lastPullIconWidget := tview.NewTextView().SetTextColor(r.theme.LastPullForeground).
		SetTextAlign(tview.AlignCenter)

	lastPullTextWidget := tview.NewTextView().SetTextColor(r.theme.LastPullForeground).
		SetTextAlign(tview.AlignRight).
		SetChangedFunc(func() { r.app.Draw() })

	lastPullWidget := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(lastPullIconWidget, 1, 0, false).
		AddItem(lastPullTextWidget, len(shortDateFormat)+1, 0, true)

	container := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statusWidget, 0, 1, false).
		AddItem(lastPullWidget, len(shortDateFormat)+2, 1, false)

	bar := statusBar{
		reader:             r,
		container:          container,
		lastPullTextWidget: lastPullTextWidget,
		lastPullIconWidget: lastPullIconWidget,
		statusWidget:       statusWidget,
		visible:            true,
	}

	r.statusBar = &bar

	return &bar
}

func (b *statusBar) setLastPullIcon() {
	b.lastPullIconWidget.SetText(pulledIcon)
}

func (b *statusBar) setLastPullTime(value *time.Time) {
	b.lastPullTextWidget.SetText(value.Local().Format(shortDateFormat))
}

func (b *statusBar) toggle() {
	if b.visible {
		b.reader.removeStatusBar(b.reader.mainPage, b)
	} else {
		b.reader.addStatusBar(b.reader.mainPage, b)
	}
	b.visible = !b.visible
}

func (b *statusBar) setNormalStatus(text string, a ...any) {
	b.statusWidget.
		SetTextColor(b.reader.theme.StatusNormalForeground).
		Clear()
	if len(a) > 0 {
		fmt.Fprintf(b.statusWidget, "%s\n", fmt.Sprintf(text, a...))
	} else {
		fmt.Fprintf(b.statusWidget, "%s\n", text)
	}
}
