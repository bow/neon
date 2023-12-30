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
	ClearStatusBar(*Display)
	CurrentFocus(*Display) tview.Primitive
	EntriesPane(*Display) tview.Primitive
	FeedsPane(*Display) tview.Primitive
	FocusFeedsPane(*Display)
	FocusEntriesPane(*Display)
	FocusNextPane(*Display)
	FocusPreviousPane(*Display)
	FocusReadingPane(*Display)
	HideIntroPopup(*Display)
	MainPage(*Display) tview.Primitive
	NotifyInfof(text string, a ...any)
	NotifyErr(err error)
	NotifyErrf(text string, a ...any)
	NotifyWarnf(text string, a ...any)
	ReadingPane(*Display) tview.Primitive
	ShowAboutPopup(*Display)
	ShowFeedsInPane(*Display, <-chan *entity.Feed)
	ShowHelpPopup(*Display)
	ShowIntroPopup(*Display)
	ShowStatsPopup(*Display, <-chan *entity.Stats)
	ToggleStatusBar(*Display)
	UnfocusPane(*Display)
}

//nolint:unused
type View struct {
	lang       *Lang
	focusStack tview.Primitive
}

//nolint:revive
func NewView() *View {
	view := View{lang: langEN}
	return &view
}

//nolint:revive
func (v *View) ClearStatusBar(dsp *Display) {
	panic("ClearStatusBar is unimplemented")
}

//nolint:revive
func (v *View) CurrentFocus(dsp *Display) tview.Primitive {
	return dsp.inner.GetFocus()
}

//nolint:revive
func (v *View) EntriesPane(dsp *Display) tview.Primitive {
	panic("EntriesPane is unimplemented")
}

//nolint:revive
func (v *View) FeedsPane(dsp *Display) tview.Primitive {
	panic("FeedsPane is unimplemented")
}

//nolint:revive
func (v *View) FocusFeedsPane(dsp *Display) {
	panic("FocusFeedsPane is unimplemented")
}

//nolint:revive
func (v *View) FocusEntriesPane(dsp *Display) {
	panic("FocusEntriesPane is unimplemented")
}

//nolint:revive
func (v *View) FocusNextPane(dsp *Display) {
	panic("FocusNextPane is unimplemented")
}

//nolint:revive
func (v *View) FocusPreviousPane(dsp *Display) {
	panic("FocusPreviousPane is unimplemented")
}

//nolint:revive
func (v *View) FocusReadingPane(dsp *Display) {
	panic("FocusReadingPane is unimplemented")
}

//nolint:revive
func (v *View) HideIntroPopup(dsp *Display) {
	panic("HideIntroPopup is unimplemented")
}

//nolint:revive
func (v *View) MainPage(dsp *Display) tview.Primitive {
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

//nolint:revive
func (v *View) ReadingPane(dsp *Display) tview.Primitive {
	panic("FeedsPane is unimplemented")
}

//nolint:revive
func (v *View) ShowAboutPopup(dsp *Display) {
	panic("ShowAboutPopup is unimplemented")
}

//nolint:revive
func (v *View) ShowFeedsInPane(dsp *Display, ch <-chan *entity.Feed) {
	panic("ShowFeedsInPane is unimplemented")
}

//nolint:revive
func (v *View) ShowHelpPopup(dsp *Display) {
	panic("ShowHelpPopup is unimplemented")
}

//nolint:revive
func (v *View) ShowIntroPopup(dsp *Display) {
	panic("ShowIntroPopup is unimplemented")
}

//nolint:revive
func (v *View) ShowStatsPopup(dsp *Display, ch <-chan *entity.Stats) {
	panic("ShowStatsPopup is unimplemented")
}

//nolint:revive
func (v *View) ToggleStatusBar(dsp *Display) {
	panic("ToggleStatusBar is unimplemented")
}

//nolint:revive
func (v *View) UnfocusPane(dsp *Display) {
	panic("UnfocusPane is unimplemented")
}

// Ensure View implements Viewer.
var _ Viewer = new(View)

//nolint:unused
type drawFunc func(screen tcell.Screen, x int, y int, w int, h int) (ix int, iy int, iw int, ih int)
