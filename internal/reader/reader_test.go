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
	screen, opr, rpo, draw := setupReaderTest(t)

	rdr := draw()

	opr.EXPECT().ToggleAboutPopup(rdr.dsp, rpo)

	screen.InjectKey(tcell.KeyRune, 'A', tcell.ModNone)
}

func TestToggleHelpPopupCalled(t *testing.T) {
	screen, opr, _, draw := setupReaderTest(t)

	rdr := draw()

	opr.EXPECT().
		ToggleHelpPopup(rdr.dsp).
		Times(2)

	screen.InjectKey(tcell.KeyRune, '?', tcell.ModNone)
	screen.InjectKey(tcell.KeyRune, 'h', tcell.ModNone)
}

func TestUnfocusFrontCalled(t *testing.T) {
	screen, opr, _, draw := setupReaderTest(t)

	rdr := draw()

	opr.EXPECT().UnfocusFront(rdr.dsp)

	screen.InjectKey(tcell.KeyEscape, ' ', tcell.ModNone)
}

func TestStartSmoke(t *testing.T) {

	screen, _, _, draw := setupReaderTest(t)

	// Since draw states are hidden at this level, the test just checks that
	// precondition == all cells empty, postcondition == at least one cell non-empty
	cellEmpty := func(w, h int) bool {
		pr, _, _, _ := screen.GetContent(w, h)
		return pr == ' ' || pr == '\x00'
	}
	empty := func() bool {
		for w := 0; w < screenW; w++ {
			for h := 0; h < screenH; h++ {
				if !cellEmpty(w, h) {
					return false
				}
			}
		}
		return true
	}
	drawn := func() bool {
		for w := 0; w < screenW; w++ {
			for h := 0; h < screenH; h++ {
				if !cellEmpty(w, h) {
					return true
				}
			}
		}
		return false
	}

	assert.True(t, empty())

	draw()
	assert.Eventually(t, drawn, 2*time.Second, 100*time.Millisecond)

	// screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
	// assert.Eventually(t, empty, 2*time.Second, 100*time.Millisecond)
}

func setupReaderTest(
	t *testing.T,
) (
	tcell.SimulationScreen,
	*MockOperator,
	*MockRepo,
	func() *Reader,
) {
	t.Helper()

	var (
		r = require.New(t)

		screen = tcell.NewSimulationScreen("UTF-8")
		opr    = NewMockOperator(gomock.NewController(t))
		rpo    = NewMockRepo(gomock.NewController(t))
	)

	var wg sync.WaitGroup
	drawf := func() *Reader {
		rdr, err := NewBuilder().
			repo(rpo).
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

	return screen, opr, rpo, drawf
}
