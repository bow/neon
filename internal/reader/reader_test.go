// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/bow/neon/internal/entity"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const screenW, screenH = 210, 60

func TestClearStatusBar(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().ClearStatusBar(rdr.display)
	tw.screen.InjectKey(tcell.KeyRune, 'c', tcell.ModNone)
}

func TestToggleAboutPopupCalled(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	bn := "caecd498-0493-4ebe-827c-ccc9f041c218"
	tw.backend.EXPECT().String().Return(bn)
	tw.opr.EXPECT().ToggleAboutPopup(rdr.display, bn)

	tw.screen.InjectKey(tcell.KeyRune, 'A', tcell.ModNone)
}

func TestFocusEntriesPane(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().FocusEntriesPane(rdr.display)
	tw.screen.InjectKey(tcell.KeyRune, 'E', tcell.ModNone)
}

func TestFocusFeedsPane(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().FocusFeedsPane(rdr.display)
	tw.screen.InjectKey(tcell.KeyRune, 'F', tcell.ModNone)
}

func TestFocusNextPane(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().FocusNextPane(rdr.display)
	tw.screen.InjectKey(tcell.KeyTab, ' ', tcell.ModNone)
}

func TestFocusPreviousPane(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().FocusPreviousPane(rdr.display)
	tw.screen.InjectKey(tcell.KeyTab, ' ', tcell.ModAlt)
}

func TestFocusReadingPane(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().FocusReadingPane(rdr.display)
	tw.screen.InjectKey(tcell.KeyRune, 'R', tcell.ModNone)
}

func TestToggleHelpPopupCalled(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().
		ToggleHelpPopup(rdr.display).
		Times(2)

	tw.screen.InjectKey(tcell.KeyRune, '?', tcell.ModNone)
	tw.screen.InjectKey(tcell.KeyRune, 'H', tcell.ModNone)
}

func TestToggleStatsPopupCalled(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.backend.EXPECT().GetStatsF(gomock.Any()).
		Return(func() (*entity.Stats, error) { return nil, nil })
	tw.opr.EXPECT().ToggleStatsPopup(rdr.display, gomock.Any())

	tw.screen.InjectKey(tcell.KeyRune, 'S', tcell.ModNone)
}

func TestToggleStatusBarCalled(t *testing.T) {
	tw := setupReaderTest(t)

	rdr := tw.draw()

	tw.opr.EXPECT().ToggleStatusBar(rdr.display)

	tw.screen.InjectKey(tcell.KeyRune, 'b', tcell.ModNone)
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
		rdr, err := NewBuilder(context.Background()).
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

			be.EXPECT().GetStatsF(gomock.Any()).
				Return(func() (*entity.Stats, error) { return nil, nil })
			opr.EXPECT().RefreshStats(gomock.Any(), gomock.Any())

			be.EXPECT().ListFeedsF(gomock.Any()).
				Return(func() ([]*entity.Feed, error) { return nil, nil })
			opr.EXPECT().ShowAllFeeds(gomock.Any(), gomock.Any())

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
