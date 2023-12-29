// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestShowOkSmoke(t *testing.T) {

	screen, draw := setupReaderTest(t)

	drawn := func() bool {
		pr, _, _, _ := screen.GetContent(0, 0)
		return pr != ' '
	}

	assert.False(t, drawn())
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

	model := NewMockModel(gomock.NewController(t))

	screen = tcell.NewSimulationScreen("UTF-8")
	r.NoError(screen.Init())
	screen.SetSize(210, 60)

	app := tview.NewApplication().
		EnableMouse(true).
		SetRoot(tview.NewBox().SetBorder(true), true).
		SetScreen(screen)

	viewer := NewMockViewer(gomock.NewController(t))
	viewer.EXPECT().
		Show().
		DoAndReturn(func() error { return app.Run() })

	var wg sync.WaitGroup
	drawf = func() *Reader {
		rdr, err := NewBuilder().
			model(model).
			viewer(viewer).
			Build()
		r.NoError(err)
		r.NotNil(rdr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			rerr := rdr.Show()
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
