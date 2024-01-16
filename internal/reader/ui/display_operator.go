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

func (do *DisplayOperator) RefreshStats(d *Display, f func() (*entity.Stats, error)) {
	stats, err := f()
	if err != nil {
		d.errEvent(err)
		return
	}
	d.setStats(stats)
}

func (do *DisplayOperator) RefreshFeeds(d *Display, f func() (<-chan entity.PullResult, error)) {
	d.infoEventf("Pulling feeds")

	var okc, errc, totalc int
	ch, err := f()
	if err != nil {
		d.errEvent(err)
		return
	}
	for pr := range ch {
		if perr := pr.Error(); perr != nil {
			d.errEventf("Pull failed for %s: %s", pr.URL(), perr)
			errc++
		} else {
			d.infoEventf("Pulled %s", pr.URL())
			// TODO: Actually process the pulled feed here.
			okc++
		}
		totalc++
	}
	if errc == 0 {
		switch okc {
		case 0:
			d.infoEventf("No feeds to pull")
		case 1:
			d.infoEventf("%d/%d feed pulled successfully", okc, totalc)
		default:
			d.infoEventf("%d/%d feeds pulled successfully", okc, totalc)
		}
	} else {
		switch okc {
		case 0:
			d.errEventf("Failed to pull any feeds")
		default:
			d.warnEventf("Only %d/%d feeds pulled successfully", okc, totalc)
		}
	}
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

func (do *DisplayOperator) ToggleAboutPopup(d *Display, backend string) {
	if name := d.frontPageName(); name == aboutPageName {
		d.hidePopup(name)
	} else if name != introPageName {
		d.setAboutPopupText(backend)
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
		do.RefreshStats(d, f)
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
