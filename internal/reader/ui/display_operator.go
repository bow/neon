// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"context"
	"time"

	"github.com/rivo/tview"

	"github.com/bow/neon/internal/reader/backend"
)

type DisplayOperator struct {
	ctx            context.Context
	focusStack     tview.Primitive
	backendTimeout time.Duration
}

func NewDisplayOperator() *DisplayOperator {
	do := DisplayOperator{
		// FIXME: Pass in these values instead of hard-coding.
		ctx:            context.Background(),
		backendTimeout: 2 * time.Second,
	}
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
func (do *DisplayOperator) NotifyInfof(d *Display, text string, a ...any) {
	panic("NotifyInfof is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) NotifyErr(d *Display, err error) {
	panic("NotifyErr is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) NotifyErrf(d *Display, text string, a ...any) {
	panic("NotifyErrf is unimplemented")
}

//nolint:revive
func (do *DisplayOperator) NotifyWarnf(d *Display, text string, a ...any) {
	panic("NotifyWarnf is unimplemented")
}

func (do *DisplayOperator) ShowIntroPopup(d *Display) {
	do.showPopup(d, introPageName)
}

func (do *DisplayOperator) ToggleAboutPopup(d *Display, b backend.Backend) {
	if name := do.frontPageName(d); name == aboutPageName {
		do.hidePopup(d, name)
	} else if name != introPageName {
		d.setAboutPopupText(b)
		do.switchPopup(d, aboutPageName, name)
	}
}

//nolint:revive
func (do *DisplayOperator) ToggleFeedsInPane(d *Display, b backend.Backend) {
	panic("ToggleFeedsInPane is unimplemented")
}

func (do *DisplayOperator) ToggleHelpPopup(d *Display) {
	if name := do.frontPageName(d); name == helpPageName {
		do.hidePopup(d, name)
	} else {
		do.switchPopup(d, helpPageName, name)
	}
}

//nolint:revive
func (do *DisplayOperator) ToggleStatsPopup(d *Display, b backend.Backend) {
	if name := do.frontPageName(d); name == statsPageName {
		do.hidePopup(d, name)
	} else if name != introPageName {
		ctx, cancel := context.WithTimeout(do.ctx, do.backendTimeout)
		defer cancel()

		res := <-b.GetStats(ctx)
		if res.Err != nil {
			// FIXME: Show error in status bar.
			panic(res.Err)
		} else {
			d.setStatsPopupValues(res.Value)
			do.switchPopup(d, statsPageName, name)
			// FIXME: Move to a more generic place.
			d.inner.Draw()
		}
	}
}

//nolint:revive
func (do *DisplayOperator) ToggleStatusBar(d *Display) {
	panic("ToggleStatusBar is unimplemented")
}

func (do *DisplayOperator) UnfocusFront(d *Display) {
	name := do.frontPageName(d)
	if name == mainPageName || name == "" {
		d.inner.SetFocus(d.root)
	} else {
		do.hidePopup(d, name)
	}
}

func (do *DisplayOperator) frontPageName(d *Display) string {
	name, _ := d.root.GetFrontPage()
	return name
}

func (do *DisplayOperator) switchPopup(d *Display, name string, currentFront string) {
	if currentFront == mainPageName {
		do.stashFocus(d)
	} else {
		d.root.HidePage(currentFront)
	}
	do.showPopup(d, name)
}

func (do *DisplayOperator) showPopup(d *Display, name string) {
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

func (do *DisplayOperator) stashFocus(d *Display) {
	do.focusStack = d.inner.GetFocus()
}

// Ensure DisplayOperator implements Operator.
var _ Operator = new(DisplayOperator)
