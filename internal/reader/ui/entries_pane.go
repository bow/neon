// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"
	"time"

	"github.com/bow/neon/internal/entity"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type entriesPane struct {
	tview.Table

	theme *Theme
	lang  *Lang

	store *entriesStore

	readingPane *readingPane
}

func newEntriesPane(theme *Theme, lang *Lang, rp *readingPane) *entriesPane {
	ep := entriesPane{
		theme: theme,
		lang:  lang,

		store: newEntriesStore(),

		readingPane: rp,
	}

	ep.initTable()

	focusf, unfocusf := ep.makeDrawFuncs()
	ep.SetDrawFunc(unfocusf)
	ep.SetFocusFunc(func() { ep.SetDrawFunc(focusf) })
	ep.SetBlurFunc(func() { ep.SetDrawFunc(unfocusf) })

	return &ep
}

func (ep *entriesPane) setEntries(entries []*entity.Entry) {
	ep.store.set(entries)
	ep.refreshEntries()
}

func (ep *entriesPane) refreshEntries() {
	rowf := ep.makeRowFuncs()

	ep.Clear()
	for i, entry := range ep.store.all() {

		colIdx := 0
		addCell := func(cell *tview.TableCell) {
			ep.SetCell(i, colIdx, cell)
			colIdx++
		}

		for _, cell := range rowf(entry) {
			addCell(cell)
		}

		ep.ScrollToBeginning()
	}
}

func (ep *entriesPane) initTable() {
	table := tview.NewTable().SetSelectable(true, false)
	ep.Table = *table
}

var year = time.Now().Year()

func (ep *entriesPane) makeRowFuncs() func(*entity.Entry) []*tview.TableCell {
	var (
		_, _, w, _      = ep.GetInnerRect()
		rowW            = w - 1 // account for padding
		timeFormat      = shortDateFormat
		timeTruncFormat = shortTruncDateFormat
		timeW           = shortDateWidth
	)

	if float32(timeW) > 0.2*float32(rowW) {
		timeFormat = compactDateFormat
		timeTruncFormat = compactTruncDateFormat
		timeW = compactDateWidth
	}

	titleW := rowW - timeW

	return func(entry *entity.Entry) []*tview.TableCell {

		titleCol := tview.NewTableCell(fmt.Sprintf("%-*s", titleW, entry.Title)).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(titleW)

		pubTS := ""
		if pubTime := entry.Published; pubTime != nil {
			tf := timeFormat
			if pubTime.Year() == year {
				tf = timeTruncFormat
			}
			pubTS = pubTime.Local().Format(tf)
		}
		pubDateCol := tview.NewTableCell(fmt.Sprintf("%*s", timeW, pubTS)).
			SetAlign(tview.AlignRight).
			SetMaxWidth(timeW)

		return []*tview.TableCell{titleCol, pubDateCol}
	}
}

// nolint:dupl
func (ep *entriesPane) makeDrawFuncs() (focusf, unfocusf drawFunc) {

	titleUF, titleF := fmtPaneTitle(ep.lang.entriesPaneTitle)

	drawf := func(
		focused bool,
	) func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		var title string
		if focused {
			title = titleF
		} else {
			title = titleUF
		}

		return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			style := ep.theme.lineStyle()
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, style)
			}

			// Write the title text.
			tview.Print(
				screen,
				title,
				x,
				y,
				width-2,
				tview.AlignLeft,
				ep.theme.titleFG,
			)

			return x + 1, y + 1, width - 2, height - 1
		}
	}

	focusf = drawf(true)
	unfocusf = drawf(false)

	return focusf, unfocusf
}

type entriesStore struct {
	items []*entity.Entry
}

func newEntriesStore() *entriesStore {
	les := entriesStore{items: make([]*entity.Entry, 0)}
	return &les
}

func (les *entriesStore) set(entries []*entity.Entry) {
	les.items = entries
}

func (les *entriesStore) all() []*entity.Entry {
	return les.items
}
