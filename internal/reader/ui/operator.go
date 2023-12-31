// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"github.com/rivo/tview"

	"github.com/bow/neon/internal/entity"
)

// Operator describes high-level UI operations.
type Operator interface {
	ClearStatusBar(d *Display)
	FocusFeedsPane(d *Display)
	FocusEntriesPane(d *Display)
	FocusNextPane(d *Display)
	FocusPreviousPane(d *Display)
	FocusReadingPane(d *Display)
	HideIntroPopup(d *Display)
	NotifyInfof(text string, a ...any)
	NotifyErr(err error)
	NotifyErrf(text string, a ...any)
	NotifyWarnf(text string, a ...any)
	Start(d *Display) error
	Stop(d *Display)
	ToggleAboutPopup(d *Display)
	ToggleFeedsInPane(d *Display, feeds <-chan *entity.Feed)
	ToggleHelpPopup(d *Display)
	ToggleIntroPopup(d *Display)
	ToggleStatsPopup(d *Display, stats <-chan *entity.Stats)
	ToggleStatusBar(d *Display)
	UnfocusPane(d *Display)
}

type DisplayOperator struct {
	focusStack tview.Primitive
}

//nolint:revive
func NewDisplayOperator() *DisplayOperator {
	do := DisplayOperator{}
	return &do
}

//nolint:revive
func (do *DisplayOperator) ClearStatusBar(d *Display) {
	panic("ClearStatusBar is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) FocusFeedsPane(d *Display) {
	panic("FocusFeedsPane is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) FocusEntriesPane(d *Display) {
	panic("FocusEntriesPane is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) FocusNextPane(d *Display) {
	panic("FocusNextPane is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) FocusPreviousPane(d *Display) {
	panic("FocusPreviousPane is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) FocusReadingPane(d *Display) {
	panic("FocusReadingPane is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) HideIntroPopup(d *Display) {
	panic("HideIntroPopup is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) NotifyInfof(text string, a ...any) {
	panic("NotifyInfof is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) NotifyErr(err error) {
	panic("NotifyErr is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) NotifyErrf(text string, a ...any) {
	panic("NotifyErrf is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) NotifyWarnf(text string, a ...any) {
	panic("NotifyWarnf is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) Start(d *Display) error {
	return d.inner.Run()
}

//nolint:revive
func (do *DisplayOperator) Stop(d *Display) {
	d.inner.Stop()
}

//nolint:revive
func (do *DisplayOperator) ToggleAboutPopup(d *Display) {
	panic("ToggleAboutPopup is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) ToggleFeedsInPane(d *Display, ch <-chan *entity.Feed) {
	panic("ToggleFeedsInPane is unimplemented")
}

func (do *DisplayOperator) ToggleHelpPopup(d *Display) {
	if name := do.frontPageName(d); name == helpPageName {
		do.hidePopup(d, name)
	} else {
		do.showPopup(d, helpPageName, name)
	}
}

//nolint:revive
func (do *DisplayOperator) ToggleIntroPopup(d *Display) {
	panic("ToggleIntroPopup is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) ToggleStatsPopup(d *Display, ch <-chan *entity.Stats) {
	panic("ToggleStatsPopup is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) ToggleStatusBar(d *Display) {
	panic("ToggleStatusBar is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) UnfocusPane(d *Display) {
	panic("UnfocusPane is unimplemented")
}

func (do *DisplayOperator) frontPageName(d *Display) string {
	name, _ := d.root.GetFrontPage()
	return name
}

func (do *DisplayOperator) showPopup(d *Display, name string, currentFront string) {
	if currentFront == mainPageName {
		do.stashFocus(d)
	} else {
		d.root.HidePage(currentFront)
	}
	d.dimMainPage()
	d.root.ShowPage(name)
}

func (do *DisplayOperator) hidePopup(d *Display, name string) {
	d.root.HidePage(name)
	d.normalizeMainPage()
	if top := do.focusStack; top != nil {
		d.inner.SetFocus(top)
	}
	do.focusStack = nil
}

func (do *DisplayOperator) stashFocus(d *Display) { do.focusStack = d.inner.GetFocus() }

// Ensure DisplayOperator implements Operator.
var _ Operator = new(DisplayOperator)
