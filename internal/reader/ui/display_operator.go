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

func (do *DisplayOperator) ClearStatusBar(d *Display) {
	d.clearEvent()
}

func (do *DisplayOperator) FocusFeedsPane(d *Display) {
	d.focusPane(d.feedsPane)
}

func (do *DisplayOperator) FocusEntriesPane(d *Display) {
	d.focusPane(d.entriesPane)
}

func (do *DisplayOperator) FocusNextPane(d *Display) {
	d.focusAdjacentPane(false)
}

func (do *DisplayOperator) FocusPreviousPane(d *Display) {
	d.focusAdjacentPane(true)
}

func (do *DisplayOperator) FocusReadingPane(d *Display) {
	d.focusPane(d.readingPane)
}

func (do *DisplayOperator) ShowFeedsInPane(d *Display, b backend.Backend) {
	ctx, cancel := do.callCtx()
	defer cancel()

	feeds, err := b.ListFeeds(ctx)
	if err != nil {
		d.errEvent(err)
		return
	}
	go func() {
		for _, feed := range feeds {
			d.feedsCh <- feed
		}
	}()
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

		stats, err := b.GetStats(ctx)
		if err != nil {
			d.errEvent(err)
			return
		}
		d.setStats(stats)
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
