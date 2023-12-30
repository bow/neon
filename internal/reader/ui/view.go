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
	FocusFeedsPane(*Display)
	FocusEntriesPane(*Display)
	FocusNextPane(*Display)
	FocusPreviousPane(*Display)
	FocusReadingPane(*Display)
	HideIntroPopup(*Display)
	NotifyInfof(text string, a ...any)
	NotifyErr(err error)
	NotifyErrf(text string, a ...any)
	NotifyWarnf(text string, a ...any)
	ToggleAboutPopup(*Display)
	ToggleFeedsInPane(*Display, <-chan *entity.Feed)
	ToggleHelpPopup(*Display)
	ToggleIntroPopup(*Display)
	ToggleStatsPopup(*Display, <-chan *entity.Stats)
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
func (v *View) ToggleAboutPopup(dsp *Display) {
	panic("ToggleAboutPopup is unimplemented")
}

//nolint:revive
func (v *View) ToggleFeedsInPane(dsp *Display, ch <-chan *entity.Feed) {
	panic("ToggleFeedsInPane is unimplemented")
}

//nolint:revive
func (v *View) ToggleHelpPopup(dsp *Display) {
	panic("ToggleHelpPopup is unimplemented")
}

//nolint:revive
func (v *View) ToggleIntroPopup(dsp *Display) {
	panic("ToggleIntroPopup is unimplemented")
}

//nolint:revive
func (v *View) ToggleStatsPopup(dsp *Display, ch <-chan *entity.Stats) {
	panic("ToggleStatsPopup is unimplemented")
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
