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

func TestStartOkSmoke(t *testing.T) {

	screen, draw := setupReaderTest(t)

	// Since draw states are hidden at this level, the test just checks that
	// precondition == all cells empty, postcondition == at least one cell non-empty
	cellEmpty := func(w, h int) bool {
		pr, _, _, _ := screen.GetContent(w, h)
		return pr == ' '
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
}

func setupReaderTest(
	t *testing.T,
) (
	screen tcell.SimulationScreen,
	drawf func() *Reader,
) {
	t.Helper()

	r := require.New(t)

	repo := NewMockRepo(gomock.NewController(t))

	screen = tcell.NewSimulationScreen("UTF-8")
	r.NoError(screen.Init())
	screen.SetSize(screenW, screenH)

	viewer := NewMockViewer(gomock.NewController(t))

	var wg sync.WaitGroup
	drawf = func() *Reader {
		rdr, err := NewBuilder().
			repo(repo).
			screen(screen).
			viewer(viewer).
			Build()
		r.NoError(err)
		r.NotNil(rdr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			rerr := rdr.Start()
			r.NoError(rerr)
		}()

		return rdr
	}

	t.Cleanup(func() {
		screen.InjectKey(tcell.KeyCtrlC, ' ', tcell.ModNone)
		wg.Wait()
	})

	return screen, drawf
}
