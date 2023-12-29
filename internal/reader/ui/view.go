// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/neon/internal/entity"
)

// Viewer describes the console reader.
type Viewer interface {
	ClearStatusBar()
	CurrentFocus() tview.Primitive
	EntriesPane() tview.Primitive
	FeedsPane() tview.Primitive
	FocusFeedsPane()
	FocusEntriesPane()
	FocusNextPane()
	FocusPreviousPane()
	FocusReadingPane()
	HideIntroPopup()
	MainPage() tview.Primitive
	NotifyInfof(text string, a ...any)
	NotifyErr(err error)
	NotifyErrf(text string, a ...any)
	NotifyWarnf(text string, a ...any)
	ReadingPane() tview.Primitive
	Show() error
	ShowAboutPopup()
	ShowFeedsInPane(<-chan *entity.Feed)
	ShowHelpPopup()
	ShowIntroPopup()
	ShowStatsPopup(<-chan *entity.Stats)
	ToggleStatusBar()
	UnfocusPane()
}

//nolint:unused
type View struct {
	theme *Theme
	lang  *Lang
	app   *tview.Application

	root *tview.Pages

	focusStack tview.Primitive
}

//nolint:revive
func NewView(theme string) (*View, error) {
	th, err := LoadTheme(theme)
	if err != nil {
		return nil, err
	}

	root := tview.NewPages()
	app := tview.NewApplication().
		EnableMouse(true).
		SetRoot(root, true)

	view := View{
		app:   app,
		theme: th,
		lang:  langEN,

		root: root,
	}

	return &view, nil
}

func (v *View) ClearStatusBar() {
	panic("ClearStatusBar is unimplemented")
}

func (v *View) CurrentFocus() tview.Primitive {
	return v.app.GetFocus()
}

func (v *View) EntriesPane() tview.Primitive {
	panic("EntriesPane is unimplemented")
}

func (v *View) FeedsPane() tview.Primitive {
	panic("FeedsPane is unimplemented")
}

func (v *View) FocusFeedsPane() {
	panic("FocusFeedsPane is unimplemented")
}

func (v *View) FocusEntriesPane() {
	panic("FocusEntriesPane is unimplemented")
}

func (v *View) FocusNextPane() {
	panic("FocusNextPane is unimplemented")
}

func (v *View) FocusPreviousPane() {
	panic("FocusPreviousPane is unimplemented")
}

func (v *View) FocusReadingPane() {
	panic("FocusReadingPane is unimplemented")
}

func (v *View) HideIntroPopup() {
	panic("HideIntroPopup is unimplemented")
}

func (v *View) MainPage() tview.Primitive {
	panic("MainPage is unimplemented")
}

//nolint:revive
func (v *View) NotifyInfof(text string, a ...any) {
	panic("NotifyInfof is unimplemented")
}

//nolint:revive
func (v *View) NotifyErr(err error) {
	panic("NotifyErr is unimplemented")
}

//nolint:revive
func (v *View) NotifyErrf(text string, a ...any) {
	panic("NotifyErrf is unimplemented")
}

//nolint:revive
func (v *View) NotifyWarnf(text string, a ...any) {
	panic("NotifyWarnf is unimplemented")
}

func (v *View) ReadingPane() tview.Primitive {
	panic("FeedsPane is unimplemented")
}

func (v *View) Show() error {
	return v.app.Run()
}

func (v *View) ShowAboutPopup() {
	panic("ShowAboutPopup is unimplemented")
}

//nolint:revive
func (v *View) ShowFeedsInPane(ch <-chan *entity.Feed) {
	panic("ShowFeedsInPane is unimplemented")
}

func (v *View) ShowHelpPopup() {
	panic("ShowHelpPopup is unimplemented")
}

func (v *View) ShowIntroPopup() {
	panic("ShowIntroPopup is unimplemented")
}

//nolint:revive
func (v *View) ShowStatsPopup(ch <-chan *entity.Stats) {
	panic("ShowStatsPopup is unimplemented")
}

func (v *View) ToggleStatusBar() {
	panic("ToggleStatusBar is unimplemented")
}

func (v *View) UnfocusPane() {
	panic("UnfocusPane is unimplemented")
}

// Ensure View implements Viewer.
var _ Viewer = new(View)

//nolint:unused
type drawFunc func(screen tcell.Screen, x int, y int, w int, h int) (ix int, iy int, iw int, ih int)
