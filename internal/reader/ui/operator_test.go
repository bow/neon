// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const screenW, screenH = 210, 60

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
	r.Equal(dsp.helpPopup.grid, item)

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
	dsp.Init(func(ek *tcell.EventKey) *tcell.EventKey { return ek })
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
