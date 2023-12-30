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

type KeyHandler = func(*tcell.EventKey) *tcell.EventKey

//nolint:unused
type View struct {
	focusStack tview.Primitive
}

//nolint:revive
func NewView() *View {
	view := View{}
	return &view
}

//nolint:revive
func (v *View) ClearStatusBar(d *Display) {
	panic("ClearStatusBar is unimplemented")
}

//nolint:revive
func (v *View) FocusFeedsPane(d *Display) {
	panic("FocusFeedsPane is unimplemented")
}

//nolint:revive
func (v *View) FocusEntriesPane(d *Display) {
	panic("FocusEntriesPane is unimplemented")
}

//nolint:revive
func (v *View) FocusNextPane(d *Display) {
	panic("FocusNextPane is unimplemented")
}

//nolint:revive
func (v *View) FocusPreviousPane(d *Display) {
	panic("FocusPreviousPane is unimplemented")
}

//nolint:revive
func (v *View) FocusReadingPane(d *Display) {
	panic("FocusReadingPane is unimplemented")
}

//nolint:revive
func (v *View) HideIntroPopup(d *Display) {
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
func (v *View) ToggleAboutPopup(d *Display) {
	panic("ToggleAboutPopup is unimplemented")
}

//nolint:revive
func (v *View) ToggleFeedsInPane(d *Display, ch <-chan *entity.Feed) {
	panic("ToggleFeedsInPane is unimplemented")
}

func (v *View) ToggleHelpPopup(d *Display) {
	if name := v.frontPageName(d); name == helpPageName {
		v.hidePopup(d, name)
	} else {
		v.showPopup(d, helpPageName, name)
	}
}

//nolint:revive
func (v *View) ToggleIntroPopup(d *Display) {
	panic("ToggleIntroPopup is unimplemented")
}

//nolint:revive
func (v *View) ToggleStatsPopup(d *Display, ch <-chan *entity.Stats) {
	panic("ToggleStatsPopup is unimplemented")
}

//nolint:revive
func (v *View) ToggleStatusBar(d *Display) {
	panic("ToggleStatusBar is unimplemented")
}

//nolint:revive
func (v *View) UnfocusPane(d *Display) {
	panic("UnfocusPane is unimplemented")
}

func (v *View) frontPageName(d *Display) string {
	name, _ := d.root.GetFrontPage()
	return name
}

func (v *View) showPopup(d *Display, name string, currentFront string) {
	if currentFront == mainPageName {
		v.stashFocus(d)
	} else {
		d.root.HidePage(currentFront)
	}
	d.dimMainPage()
	d.root.ShowPage(name)
}

func (v *View) hidePopup(d *Display, name string) {
	d.root.HidePage(name)
	d.normalizeMainPage()
	if top := v.focusStack; top != nil {
		d.inner.SetFocus(top)
	}
	v.focusStack = nil
}

func (v *View) stashFocus(d *Display) { v.focusStack = d.inner.GetFocus() }

// Ensure View implements Viewer.
var _ Viewer = new(View)

//nolint:unused
type drawFunc func(screen tcell.Screen, x int, y int, w int, h int) (ix int, iy int, iw int, ih int)
