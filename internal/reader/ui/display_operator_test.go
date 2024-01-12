// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bow/neon/internal/entity"
)

const screenW, screenH = 210, 60

func TestToggleAboutPopup(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	draw, opr, dsp := setupDisplayOperatorTest(t)
	be := NewMockBackend(gomock.NewController(t))

	draw()

	name, item := dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
	a.Nil(dsp.aboutPopup.content)

	bn1 := uuid.NewString()
	be.EXPECT().String().Return(bn1)

	opr.ToggleAboutPopup(dsp, be)
	name, item = dsp.root.GetFrontPage()
	a.Equal(aboutPageName, name)
	r.Equal(dsp.aboutPopup, item)
	r.NotNil(dsp.aboutPopup.content)
	c1, typeok1 := dsp.aboutPopup.content.(*tview.TextView)
	r.True(typeok1)
	a.Contains(c1.GetText(true), bn1)

	be.EXPECT().String().Times(0)

	opr.ToggleAboutPopup(dsp, be)
	name, item = dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
	a.NotNil(dsp.aboutPopup.content)

	bn2 := uuid.NewString()
	be.EXPECT().String().Return(bn2)

	opr.ToggleAboutPopup(dsp, be)
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

	be := NewMockBackend(gomock.NewController(t))
	stats := entity.Stats{
		NumFeeds:         40,
		NumEntries:       236,
		NumEntriesUnread: 5,
		// TODO: Test using non-nil value.
		LastPullTime:         nil,
		MostRecentUpdateTime: nil,
	}
	be.EXPECT().
		GetStats(gomock.Any()).
		Return(&stats, nil).
		Times(1)

	opr.ToggleStatsPopup(dsp, be)
	name, item = dsp.root.GetFrontPage()
	a.Equal(statsPageName, name)
	r.Equal(dsp.statsPopup, item)
	c, typeok := dsp.statsPopup.content.(*tview.TextView)
	r.True(typeok)
	a.Contains(c.GetText(true), "Total: 40")
	a.Contains(c.GetText(true), "Total : 236")
	a.Contains(c.GetText(true), "Unread: 5")

	opr.ToggleStatsPopup(dsp, be)
	name, item = dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
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
				250*time.Millisecond,
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
