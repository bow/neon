// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const screenW, screenH = 210, 60

func TestToggleAboutPopupCalled(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().ToggleAboutPopup(rdr.display, tw.backend)
	tw.screen.InjectKey(tcell.KeyRune, 'A', tcell.ModNone)
}

func TestToggleHelpPopupCalled(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().
		ToggleHelpPopup(rdr.display).
		Times(2)

	tw.screen.InjectKey(tcell.KeyRune, '?', tcell.ModNone)
	tw.screen.InjectKey(tcell.KeyRune, 'h', tcell.ModNone)
}

func TestToggleStatsPopupCalled(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().ToggleStatsPopup(rdr.display, rdr.backend)

	tw.screen.InjectKey(tcell.KeyRune, 'S', tcell.ModNone)
}

func TestShowIntroPopupCalled(t *testing.T) {
	tw := setupReaderTest(t)

	tw.introSeen = false
	tw.state.EXPECT().MarkIntroSeen()
	tw.opr.EXPECT().ShowIntroPopup(gomock.Any())
	tw.draw()
}

func TestUnfocusFrontCalled(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().UnfocusFront(rdr.display)

	tw.screen.InjectKey(tcell.KeyEscape, ' ', tcell.ModNone)
}

func TestStartSmoke(t *testing.T) {

	tw := setupReaderTest(t)

	// Since draw states are hidden at this level, the test just checks that
	// precondition == all cells empty, postcondition == at least one cell non-empty
	empty := func() bool {
		for y := 0; y < screenH; y++ {
			for x := 0; x < screenW; x++ {
				pr, _, _, _ := tw.screen.GetContent(x, y)
				if pr != '\x00' {
					return false
				}
			}
		}
		return true
	}
	drawn := func() bool {
		for y := 0; y < screenH; y++ {
			for x := 0; x < screenW; x++ {
				pr, _, _, _ := tw.screen.GetContent(x, y)
				if pr != '\x00' && pr != ' ' {
					return true
				}
			}
		}
		return false
	}
	pollTimeout, tickFreq := 2*time.Second, 100*time.Millisecond

	assert.Eventually(t, empty, pollTimeout, tickFreq)

	tw.draw()
	assert.Eventually(t, drawn, pollTimeout, tickFreq)

	tw.screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
	assert.Eventually(t, empty, pollTimeout, tickFreq)
}

type testWrapper struct {
	screen  tcell.SimulationScreen
	opr     *MockOperator
	backend *MockBackend
	state   *MockState
	draw    func() *Reader

	introSeen bool
}

func setupReaderTest(t *testing.T) *testWrapper {
	t.Helper()

	var (
		r = require.New(t)

		screen = tcell.NewSimulationScreen("UTF-8")
		opr    = NewMockOperator(gomock.NewController(t))
		be     = NewMockBackend(gomock.NewController(t))
		stt    = NewMockState(gomock.NewController(t))
	)

	tw := &testWrapper{}
	tw.introSeen = true

	var wg sync.WaitGroup
	drawf := func() *Reader {
		rdr, err := NewBuilder().
			backend(be).
			screen(screen).
			operator(opr).
			state(stt).
			Build()
		r.NoError(err)
		r.NotNil(rdr)

		// This is called here because the underlying App calls screen.Init, which,
		// among other things, resets its size.
		screen.SetSize(screenW, screenH)

		wg.Add(1)
		go func() {
			defer wg.Done()
			stt.EXPECT().IntroSeen().Return(tw.introSeen)
			rerr := rdr.Start()
			r.NoError(rerr)
		}()

		t.Cleanup(func() {
			screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
			wg.Wait()
		})

		return rdr
	}

	tw.screen = screen
	tw.opr = opr
	tw.backend = be
	tw.state = stt
	tw.draw = drawf

	return tw
}
