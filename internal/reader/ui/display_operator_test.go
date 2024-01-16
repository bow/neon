// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bow/neon/internal/entity"
)

const screenW, screenH = 210, 60

func TestClearStatusBar(t *testing.T) {
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	w := dsp.bar.eventsWidget

	r.Empty(w.GetText(true))
	fmt.Fprintf(&w.TextView, "foobar")

	r.NotEmpty(w.GetText(true))
	opr.ClearStatusBar(dsp)
	r.Empty(w.GetText(true))
}

func TestFocusEntriesPane(t *testing.T) {
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	r.Equal(dsp.mainPage, dsp.inner.GetFocus())

	opr.FocusEntriesPane(dsp)
	r.Equal(dsp.entriesPane, dsp.inner.GetFocus())
}

func TestFocusNextPane(t *testing.T) {
	r := require.New(t)

	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	r.Equal(dsp.mainPage, dsp.inner.GetFocus())

	opr.FocusNextPane(dsp)
	r.Equal(dsp.feedsPane, dsp.inner.GetFocus())

	opr.FocusNextPane(dsp)
	r.Equal(dsp.entriesPane, dsp.inner.GetFocus())

	opr.FocusNextPane(dsp)
	r.Equal(dsp.readingPane, dsp.inner.GetFocus())

	opr.FocusNextPane(dsp)
	r.Equal(dsp.feedsPane, dsp.inner.GetFocus())
}

func TestFocusPreviousPane(t *testing.T) {
	r := require.New(t)

	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	r.Equal(dsp.mainPage, dsp.inner.GetFocus())

	opr.FocusPreviousPane(dsp)
	r.Equal(dsp.readingPane, dsp.inner.GetFocus())

	opr.FocusPreviousPane(dsp)
	r.Equal(dsp.entriesPane, dsp.inner.GetFocus())

	opr.FocusPreviousPane(dsp)
	r.Equal(dsp.feedsPane, dsp.inner.GetFocus())

	opr.FocusPreviousPane(dsp)
	r.Equal(dsp.readingPane, dsp.inner.GetFocus())
}

func TestFocusReadingPane(t *testing.T) {
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	r.Equal(dsp.mainPage, dsp.inner.GetFocus())

	opr.FocusReadingPane(dsp)
	r.Equal(dsp.readingPane, dsp.inner.GetFocus())
}

func TestFocusFeedPane(t *testing.T) {
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	r.Equal(dsp.mainPage, dsp.inner.GetFocus())

	opr.FocusFeedsPane(dsp)
	r.Equal(dsp.feedsPane, dsp.inner.GetFocus())
}

func TestShowAllFeedsErr(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	r.Empty(dsp.feedsPane.GetRoot().GetChildren())
	opr.ShowAllFeeds(
		dsp,
		func() ([]*entity.Feed, error) { return nil, fmt.Errorf("fail") },
	)
	a.Empty(dsp.feedsPane.GetRoot().GetChildren())
	a.Eventually(
		func() bool { return strings.Contains(dsp.bar.eventsWidget.GetText(true), "fail") },
		2*time.Second,
		500*time.Millisecond,
	)
}

func TestShowAllFeedsOk(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	groupNodes := func() []*tview.TreeNode {
		return dsp.feedsPane.GetRoot().GetChildren()
	}

	feedNodes := func() []*tview.TreeNode {
		fns := make([]*tview.TreeNode, 0)
		for _, gn := range groupNodes() {
			fns = append(fns, gn.GetChildren()...)
		}
		return fns
	}

	draw()

	r.Empty(groupNodes())
	opr.ShowAllFeeds(
		dsp,
		func() ([]*entity.Feed, error) {
			feeds := []*entity.Feed{
				{
					ID:         entity.ID(1),
					Title:      "Feed W",
					FeedURL:    "http://w.com/feed.xml",
					Subscribed: twoWeeksAgo,
					LastPulled: twoWeeksAgo,
					Updated:    &twoWeeksAgo,
				},
				{
					ID:         entity.ID(2),
					Title:      "Feed D",
					FeedURL:    "http://d.com/feed.xml",
					Subscribed: threeDaysAgo,
					LastPulled: yesterday,
					Updated:    &threeDaysAgo,
				},
				{
					ID:         entity.ID(6),
					Title:      "Feed Y",
					FeedURL:    "http://y.com/feed.xml",
					Subscribed: yesterday,
					LastPulled: yesterday,
					Updated:    &yesterday,
				},
				{
					ID:         entity.ID(8),
					Title:      "Feed N",
					FeedURL:    "http://n.com/feed.xml",
					Subscribed: yesterday,
					LastPulled: now,
					Updated:    &now,
				},
			}
			return feeds, nil
		},
	)

	a.Eventually(
		func() bool { return len(groupNodes()) == 3 },
		2*time.Second,
		500*time.Millisecond,
	)
	a.Len(feedNodes(), 4)
}

func TestShowIntroPopup(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	// This call before draw() is intentional and simulates expected use.
	opr.ShowIntroPopup(dsp)

	draw()

	name, item := dsp.root.GetFrontPage()
	a.Equal(introPageName, name)
	r.Equal(dsp.introPopup, item)

	opr.UnfocusFront(dsp)
	name, item = dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
}

func TestToggleAboutPopup(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	name, item := dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
	a.Nil(dsp.aboutPopup.content)

	bn1 := "6541de58-41ed-4117-920b-24d472057c2c"
	opr.ToggleAboutPopup(dsp, bn1)
	name, item = dsp.root.GetFrontPage()
	a.Equal(aboutPageName, name)
	r.Equal(dsp.aboutPopup, item)
	r.NotNil(dsp.aboutPopup.content)
	c1, typeok1 := dsp.aboutPopup.content.(*tview.TextView)
	r.True(typeok1)
	a.Contains(c1.GetText(true), bn1)

	opr.ToggleAboutPopup(dsp, "")
	name, item = dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
	a.NotNil(dsp.aboutPopup.content)

	bn2 := "411068b3-51e5-4565-b768-d53e17af98e6"
	opr.ToggleAboutPopup(dsp, bn2)
	name, item = dsp.root.GetFrontPage()
	a.Equal(aboutPageName, name)
	r.Equal(dsp.aboutPopup, item)
	r.NotNil(dsp.aboutPopup.content)
	c2, typeok2 := dsp.aboutPopup.content.(*tview.TextView)
	r.True(typeok2)
	a.Contains(c2.GetText(true), bn2)
}

func TestToggleHelpPopup(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	name, item := dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)

	opr.ToggleHelpPopup(dsp)
	name, item = dsp.root.GetFrontPage()
	a.Equal(helpPageName, name)
	r.Equal(dsp.helpPopup, item)

	opr.ToggleHelpPopup(dsp)
	name, item = dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
}

func TestToggleStatsPopup(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	name, item := dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)

	stats := entity.Stats{
		NumFeeds:         40,
		NumEntries:       236,
		NumEntriesUnread: 5,
		// TODO: Test using non-nil value.
		LastPullTime:         nil,
		MostRecentUpdateTime: nil,
	}
	f := func() (*entity.Stats, error) { return &stats, nil }

	opr.ToggleStatsPopup(dsp, f)
	name, item = dsp.root.GetFrontPage()
	a.Equal(statsPageName, name)
	r.Equal(dsp.statsPopup, item)
	c, typeok := dsp.statsPopup.content.(*tview.TextView)
	r.True(typeok)
	a.Contains(c.GetText(true), "Total: 40")
	a.Contains(c.GetText(true), "Total : 236")
	a.Contains(c.GetText(true), "Unread: 5")

	opr.ToggleStatsPopup(dsp, f)
	name, item = dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
}

func TestToggleStatusBar(t *testing.T) {
	a := assert.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	a.True(dsp.barVisible)

	opr.ToggleStatusBar(dsp)
	a.False(dsp.barVisible)

	opr.ToggleStatusBar(dsp)
	a.True(dsp.barVisible)
}

func TestUnfocusFront(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)

	draw()

	opr.UnfocusFront(dsp)
	name, item := dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)

	opr.ToggleHelpPopup(dsp)
	name, item = dsp.root.GetFrontPage()
	a.Equal(helpPageName, name)
	r.Equal(dsp.helpPopup, item)

	opr.UnfocusFront(dsp)
	name, item = dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
}

func setupDisplayOperatorTest(t *testing.T) (
	func(),
	*DisplayOperator,
	*Display,
) {
	t.Helper()

	var (
		r = require.New(t)

		screen = tcell.NewSimulationScreen("UTF-8")
		dsp    = newTestDisplay(t, screen)
	)
	var stopWaiter sync.WaitGroup
	drawf := func() {
		// This is called here because the underlying App calls screen.Init, which,
		// among other things, resets its size.
		screen.SetSize(screenW, screenH)

		stopWaiter.Add(1)
		go func() {
			defer stopWaiter.Done()
			rerr := dsp.Start()
			r.NoError(rerr)
		}()

		var startWaiter sync.WaitGroup
		startWaiter.Add(1)
		go func() {
			defer startWaiter.Done()
			r.Eventually(
				func() bool { return screenDrawn(t, screen) },
				2*time.Second,
				100*time.Millisecond,
			)
		}()
		startWaiter.Wait()
	}

	t.Cleanup(func() {
		if dsp != nil && dsp.inner != nil {
			dsp.inner.Stop()
			stopWaiter.Wait()
		}
	})

	return drawf, NewDisplayOperator(context.Background(), 1*time.Second), dsp
}

func newTestDisplay(t *testing.T, screen tcell.Screen) *Display {
	t.Helper()

	r := require.New(t)
	dsp, err := NewDisplay(screen, "dark")
	r.NoError(err)
	r.NotNil(dsp)
	dsp.SetHandlers(func(ek *tcell.EventKey) *tcell.EventKey { return ek })
	return dsp
}

func screenDrawn(t *testing.T, screen tcell.Screen) bool {
	t.Helper()

	for y := 0; y < screenH; y++ {
		for x := 0; x < screenW; x++ {
			pr, _, _, _ := screen.GetContent(y, x)
			// \x00 is when cell is invalid, ' ' is when cell is drawn as empty
			if pr != '\x00' && pr != ' ' {
				return true
			}
		}
	}
	return false
}
