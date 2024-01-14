// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"context"
	"time"

	"github.com/bow/neon/internal/entity"
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

func (do *DisplayOperator) ShowAllFeeds(d *Display, f func() ([]*entity.Feed, error)) {
	feeds, err := f()
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

func (do *DisplayOperator) ToggleAboutPopup(d *Display, f func() string) {
	if name := d.frontPageName(); name == aboutPageName {
		d.hidePopup(name)
	} else if name != introPageName {
		d.setAboutPopupText(f())
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

func (do *DisplayOperator) ToggleStatsPopup(d *Display, f func() (*entity.Stats, error)) {
	if name := d.frontPageName(); name == statsPageName {
		d.hidePopup(name)
	} else if name != introPageName {
		stats, err := f()
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

// Ensure DisplayOperator implements Operator.
var _ Operator = new(DisplayOperator)
