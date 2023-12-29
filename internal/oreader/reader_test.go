// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package oreader

import (
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bow/neon/api"
)

const (
	// how long to wait by default for delayed test events.
	pollTimeout = 2 * time.Second
	// how frequent to check delayed test events by default.
	tickFreq = 100 * time.Millisecond
)

func TestShowSmoke(t *testing.T) {

	screen, draw := setupReaderTest(t)

	drawn := func() bool {
		return screenCellEqual(t, screen, 0, 0, tview.BoxDrawingsLightHorizontal)
	}

	assert.False(t, drawn())
	draw()
	assert.Eventually(t, drawn, pollTimeout, tickFreq)
}

func TestHelpPopupSmoke(t *testing.T) {

	screen, draw := setupReaderTest(t)

	draw()

	dimmed := func() bool {
		return screenForegroundEqual(t, screen, 0, 0, darkForegroundDim)
	}

	assert.False(t, dimmed())
	screen.InjectKey(tcell.KeyRune, 'h', tcell.ModNone)
	assert.Eventually(t, dimmed, pollTimeout, tickFreq)
}

func TestFocusFeedsPane(t *testing.T) {
	screen, draw := setupReaderTest(t)
	reader := draw()

	marked := func() bool { return screenCellEqual(t, screen, 0, 0, '▶') }
	unmarked := func() bool {
		return screenCellEqual(t, screen, 0, 0, tview.BoxDrawingsLightHorizontal)
	}

	focused := func() bool { return reader.app.GetFocus() == reader.feedsPane }
	unfocused := func() bool { return !focused() }

	assert.Eventually(t, unfocused, pollTimeout, tickFreq)
	assert.True(t, unmarked())

	screen.InjectKey(tcell.KeyRune, 'F', tcell.ModNone)
	assert.Eventually(t, focused, pollTimeout, tickFreq)
	assert.True(t, marked())

	screen.InjectKey(tcell.KeyEsc, ' ', tcell.ModNone)
	assert.Eventually(t, unfocused, pollTimeout, tickFreq)
	assert.True(t, unmarked())
}

func setupReaderTest(
	t *testing.T,
) (
	screen tcell.SimulationScreen,
	drawf func() *Reader,
) {
	t.Helper()

	r := require.New(t)

	client := NewMockNeonClient(gomock.NewController(t))
	// Needed since we call the list feeds endpoint prior to Show.
	client.EXPECT().
		ListFeeds(gomock.Any(), gomock.Any()).
		Return(
			&api.ListFeedsResponse{Feeds: nil},
			nil,
		)
	// Needed since we call the stats endpoint prior to Show.
	client.EXPECT().
		GetStats(gomock.Any(), gomock.Any()).
		Return(
			&api.GetStatsResponse{
				Global: &api.GetStatsResponse_Stats{
					NumFeeds:             2,
					NumEntries:           5,
					NumEntriesUnread:     5,
					LastPullTime:         nil,
					MostRecentUpdateTime: nil,
				},
			},
			nil,
		)

	screen = tcell.NewSimulationScreen("UTF-8")
	err := screen.Init()
	r.NoError(err)
	screen.SetSize(200, 45)

	var wg sync.WaitGroup
	drawf = func() *Reader {
		irdr, ierr := NewBuilder().
			client(client).
			screen(screen).
			Build()
		r.NoError(ierr)
		r.NotNil(irdr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			rerr := irdr.Show()
			r.NoError(rerr)
		}()

		return irdr
	}

	t.Cleanup(func() {
		screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
		wg.Wait()
	})

	return screen, drawf
}

func screenCell(t *testing.T, screen tcell.Screen, x, y int) (rune, tcell.Style) {
	t.Helper()
	pr, _, st, _ := screen.GetContent(x, y)
	return pr, st
}

func screenForegroundEqual(
	t *testing.T,
	screen tcell.Screen,
	x, y int,
	expected tcell.Color,
) bool {
	t.Helper()
	_, st := screenCell(t, screen, x, y)
	fg, _, _ := st.Decompose()
	return expected == fg
}

func screenCellEqual(t *testing.T, screen tcell.Screen, x, y int, expected rune) bool {
	t.Helper()
	pr, _ := screenCell(t, screen, x, y)
	return expected == pr
}
