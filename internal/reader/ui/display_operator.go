// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"context"
	"time"

	"github.com/bow/neon/internal/reader/backend"
)

type DisplayOperator struct {
	ctx         context.Context
	callTimeout time.Duration
}

func NewDisplayOperator(
	ctx context.Context,
	callTimeout time.Duration,
) *DisplayOperator {
	do := DisplayOperator{
		ctx:         ctx,
		callTimeout: callTimeout,
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
func (do *DisplayOperator) ShowFeedsInPane(d *Display, b backend.Backend) {
	panic("ShowFeedsInPane is unimplemented")
}

func (do *DisplayOperator) ShowIntroPopup(d *Display) {
	d.showPopup(introPageName)
}

func (do *DisplayOperator) ToggleAboutPopup(d *Display, b backend.Backend) {
	if name := d.frontPageName(); name == aboutPageName {
		d.hidePopup(name)
	} else if name != introPageName {
		d.setAboutPopupText(b)
		d.switchPopup(aboutPageName, name)
	}
}

func (do *DisplayOperator) ToggleHelpPopup(d *Display) {
	if name := d.frontPageName(); name == helpPageName {
		d.hidePopup(name)
	} else {
		d.switchPopup(helpPageName, name)
	}
}

func (do *DisplayOperator) ToggleStatsPopup(d *Display, b backend.Backend) {
	if name := d.frontPageName(); name == statsPageName {
		d.hidePopup(name)
	} else if name != introPageName {
		ctx, cancel := do.callCtx()
		defer cancel()

		res := <-b.GetStats(ctx)
		if res.Err != nil {
			d.errEvent(res.Err)
			return
		}
		d.setStats(res.Value)
		d.switchPopup(statsPageName, name)
	}
}

func (do *DisplayOperator) ToggleStatusBar(d *Display) {
	d.toggleStatusBar()
}

func (do *DisplayOperator) UnfocusFront(d *Display) {
	name := d.frontPageName()
	if name == mainPageName || name == "" {
		d.inner.SetFocus(d.root)
	} else {
		d.hidePopup(name)
	}
}

func (do *DisplayOperator) callCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(do.ctx, do.callTimeout)
}

// Ensure DisplayOperator implements Operator.
var _ Operator = new(DisplayOperator)
