// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"sync"
	"testing"
	"time"

	"github.com/bow/lens/api"
	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShowSmoke(t *testing.T) {

	screen, client, draw := setupReaderTest(t)

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

	drawn := func() bool {
		return screenCellEqual(t, screen, 0, 0, tview.BoxDrawingsLightHorizontal)
	}

	assert.False(t, drawn())
	draw()
	assert.Eventually(t, drawn, 2*time.Second, 100*time.Millisecond)
}

func setupReaderTest(
	t *testing.T,
) (
	screen tcell.SimulationScreen,
	client *MockLensClient,
	drawf func(),
) {
	t.Helper()

	r := require.New(t)

	client = NewMockLensClient(gomock.NewController(t))

	screen = tcell.NewSimulationScreen("UTF-8")
	err := screen.Init()
	r.NoError(err)
	screen.SetSize(200, 45)

	var wg sync.WaitGroup
	drawf = func() {
		rdr, err := NewBuilder().
			client(client).
			screen(screen).
			Build()
		r.NoError(err)
		r.NotNil(rdr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			rerr := rdr.Show()
			r.NoError(rerr)
		}()
	}

	t.Cleanup(func() {
		screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
		wg.Wait()
	})

	return screen, client, drawf
}

func screenCell(t *testing.T, screen tcell.Screen, x, y int) rune {
	t.Helper()
	pr, _, _, _ := screen.GetContent(x, y)
	return pr
}

func screenCellEqual(t *testing.T, screen tcell.Screen, x, y int, expected rune) bool {
	t.Helper()
	return expected == screenCell(t, screen, x, y)
}
