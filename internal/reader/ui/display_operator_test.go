// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const screenW, screenH = 210, 60

func TestToggleAboutPopup(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	draw := setupDisplayOperatorTest(t)
	rpo := NewMockRepo(gomock.NewController(t))

	opr, dsp := draw()

	name, item := dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
	a.Nil(dsp.aboutPopup.content)

	backend1 := uuid.NewString()
	rpo.EXPECT().Backend().Return(backend1)

	opr.ToggleAboutPopup(dsp, rpo)
	name, item = dsp.root.GetFrontPage()
	a.Equal(aboutPageName, name)
	r.Equal(dsp.aboutPopup, item)
	r.NotNil(dsp.aboutPopup.content)
	c1, typeok1 := dsp.aboutPopup.content.(*tview.TextView)
	r.True(typeok1)
	a.Contains(c1.GetText(true), backend1)

	rpo.EXPECT().Backend().Times(0)

	opr.ToggleAboutPopup(dsp, rpo)
	name, item = dsp.root.GetFrontPage()
	a.Equal(mainPageName, name)
	r.Equal(dsp.mainPage, item)
	a.NotNil(dsp.aboutPopup.content)

	backend2 := uuid.NewString()
	rpo.EXPECT().Backend().Return(backend2)

	opr.ToggleAboutPopup(dsp, rpo)
	name, item = dsp.root.GetFrontPage()
	a.Equal(aboutPageName, name)
	r.Equal(dsp.aboutPopup, item)
	r.NotNil(dsp.aboutPopup.content)
	c2, typeok2 := dsp.aboutPopup.content.(*tview.TextView)
	r.True(typeok2)
	a.Contains(c2.GetText(true), backend2)
}

func TestToggleHelpPopup(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	draw := setupDisplayOperatorTest(t)

	opr, dsp := draw()

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

func setupDisplayOperatorTest(t *testing.T) func() (*DisplayOperator, *Display) {
	t.Helper()

	r := require.New(t)
	screen := tcell.NewSimulationScreen("UTF-8")
	r.NoError(screen.Init())
	screen.SetSize(screenW, screenH)

	var (
		stopWaiter sync.WaitGroup
		dsp        *Display
	)
	drawf := func() (*DisplayOperator, *Display) {
		dsp = newTestDisplay(t, screen)

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

		return NewDisplayOperator(), dsp
	}

	t.Cleanup(func() {
		if dsp != nil && dsp.inner != nil {
			dsp.inner.Stop()
			stopWaiter.Wait()
		}
	})

	return drawf
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

	for w := 0; w < screenW; w++ {
		for h := 0; h < screenH; h++ {
			pr, _, _, _ := screen.GetContent(w, h)
			if pr != ' ' {
				return true
			}
		}
	}
	return false
}
