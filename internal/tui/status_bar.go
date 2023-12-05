package tui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

const pulledIcon = "▼"

type statusBar struct {
	theme     *Theme
	container *tview.Flex

	statusWidget       *tview.TextView
	lastPullTextWidget *tview.TextView
	lastPullIconWidget *tview.TextView

	visible bool
}

func newStatusBar(theme *Theme) *statusBar {

	statusWidget := tview.NewTextView().SetTextAlign(tview.AlignLeft)

	lastPullIconWidget := tview.NewTextView().SetTextColor(theme.LastPullForeground).
		SetTextAlign(tview.AlignCenter)

	lastPullTextWidget := tview.NewTextView().SetTextColor(theme.LastPullForeground).
		SetTextAlign(tview.AlignRight)

	lastPullWidget := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(lastPullIconWidget, 1, 0, false).
		AddItem(lastPullTextWidget, len(shortDateFormat)+1, 0, true)

	container := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statusWidget, 0, 1, false).
		AddItem(lastPullWidget, len(shortDateFormat)+2, 1, false)

	bar := statusBar{
		theme:              theme,
		container:          container,
		lastPullTextWidget: lastPullTextWidget,
		lastPullIconWidget: lastPullIconWidget,
		statusWidget:       statusWidget,
		visible:            true,
	}

	return &bar
}

func (b *statusBar) setChangedFunc(handler func()) *statusBar {
	b.statusWidget.SetChangedFunc(handler)
	b.lastPullIconWidget.SetChangedFunc(handler)
	b.lastPullTextWidget.SetChangedFunc(handler)
	return b
}

func (b *statusBar) setLastPullIcon() {
	b.lastPullIconWidget.SetText(pulledIcon)
}

func (b *statusBar) setLastPullTime(value *time.Time) {
	b.lastPullTextWidget.SetText(value.Local().Format(shortDateFormat))
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

func (b *statusBar) setNormalStatus(text string, a ...any) {
	b.statusWidget.
		SetTextColor(b.theme.StatusNormalForeground).
		Clear()
	if len(a) > 0 {
		fmt.Fprintf(b.statusWidget, "%s\n", fmt.Sprintf(text, a...))
	} else {
		fmt.Fprintf(b.statusWidget, "%s\n", text)
	}
}
