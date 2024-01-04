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
	screen, opr, be, draw := setupReaderTest(t)

	rdr := draw()

	opr.EXPECT().ToggleAboutPopup(rdr.display, be)

	screen.InjectKey(tcell.KeyRune, 'A', tcell.ModNone)
}

func TestToggleHelpPopupCalled(t *testing.T) {
	screen, opr, _, draw := setupReaderTest(t)

	rdr := draw()

	opr.EXPECT().
		ToggleHelpPopup(rdr.display).
		Times(2)

	screen.InjectKey(tcell.KeyRune, '?', tcell.ModNone)
	screen.InjectKey(tcell.KeyRune, 'h', tcell.ModNone)
}

func TestUnfocusFrontCalled(t *testing.T) {
	screen, opr, _, draw := setupReaderTest(t)

	rdr := draw()

	opr.EXPECT().UnfocusFront(rdr.display)

	screen.InjectKey(tcell.KeyEscape, ' ', tcell.ModNone)
}

func TestStartSmoke(t *testing.T) {

	screen, _, _, draw := setupReaderTest(t)

	// Since draw states are hidden at this level, the test just checks that
	// precondition == all cells empty, postcondition == at least one cell non-empty
	empty := func() bool {
		for y := 0; y < screenH; y++ {
			for x := 0; x < screenW; x++ {
				pr, _, _, _ := screen.GetContent(x, y)
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
				pr, _, _, _ := screen.GetContent(x, y)
				if pr != '\x00' && pr != ' ' {
					return true
				}
			}
		}
		return false
	}
	pollTimeout, tickFreq := 2*time.Second, 100*time.Millisecond

	assert.Eventually(t, empty, pollTimeout, tickFreq)

	draw()
	assert.Eventually(t, drawn, pollTimeout, tickFreq)

	screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
	assert.Eventually(t, empty, pollTimeout, tickFreq)
}

func setupReaderTest(
	t *testing.T,
) (
	tcell.SimulationScreen,
	*MockOperator,
	*MockBackend,
	func() *Reader,
) {
	t.Helper()

	var (
		r = require.New(t)

		screen = tcell.NewSimulationScreen("UTF-8")
		opr    = NewMockOperator(gomock.NewController(t))
		be     = NewMockBackend(gomock.NewController(t))
	)

	var wg sync.WaitGroup
	drawf := func() *Reader {
		rdr, err := NewBuilder().
			backend(be).
			screen(screen).
			operator(opr).
			Build()
		r.NoError(err)
		r.NotNil(rdr)

		// This is called here because the underlying App calls screen.Init, which,
		// among other things, resets its size.
		screen.SetSize(screenW, screenH)

		wg.Add(1)
		go func() {
			defer wg.Done()
			rerr := rdr.Start()
			r.NoError(rerr)
		}()

		t.Cleanup(func() {
			screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
			wg.Wait()
		})

		return rdr
	}

	return screen, opr, be, drawf
}
