// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"bufio"
	"fmt"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	leftPopupMargin      = 2
	rightPopupMargin     = 2
	verticalPopupPadding = 4
)

type popup struct {
	tview.Grid

	content tview.Primitive
	frame   *tview.Frame
}

func (p *popup) setContent(prim tview.Primitive) {
	p.frame.SetPrimitive(prim)
	p.content = prim
}

func (p *popup) setWidth(w int) {
	p.Grid.SetColumns(0, w, 0)
}

func (p *popup) setGridRows(rows []int) {
	p.Grid.SetRows(rows...)
}

func newPopup(
	title string,
	titleColorFG tcell.Color,
	topPadding, bottomPadding int,
) *popup {
	var content tview.Primitive = nil

	frame := tview.NewFrame(content).
		SetBorders(topPadding, bottomPadding, 0, 0, leftPopupMargin, rightPopupMargin)

	frame.SetBorder(true).
		SetTitle(fmt.Sprintf(" %s ", title)).
		SetTitleColor(titleColorFG)

	grid := tview.NewGrid().
		AddItem(frame, 1, 1, 1, 1, 0, 0, true)

	p := popup{Grid: *grid, frame: frame, content: content}

	return &p
}

func newFilledPopup(
	title string,
	content tview.Primitive,
	titleColorFG tcell.Color,
	topPadding, bottomPadding int,
	width int,
	gridRows []int,
) *popup {

	frame := tview.NewFrame(content).
		SetBorders(topPadding, bottomPadding, 0, 0, leftPopupMargin, rightPopupMargin)

	frame.SetBorder(true).
		SetTitle(fmt.Sprintf(" %s ", title)).
		SetTitleColor(titleColorFG)

	grid := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(gridRows...).
		AddItem(frame, 1, 1, 1, 1, 0, 0, true)

	p := popup{Grid: *grid, frame: frame, content: content}

	return &p
}

func centerBanner(text string, width int) string {
	if width == 0 {
		return text
	}
	maxLineWidth, lines := textWidth(text)
	if maxLineWidth > width {
		return text
	}

	leftPad := strings.Repeat(" ", ((width-maxLineWidth)/2)-leftPopupMargin-1)
	paddedLines := make([]string, len(lines))
	for i, line := range lines {
		paddedLines[i] = fmt.Sprintf("%s%s", leftPad, line)
	}

	sep := "\n"
	if runtime.GOOS == "windows" {
		sep = "\r\n"
	}

	return strings.Join(paddedLines, sep)
}

func textWidth(text string) (int, []string) {
	var (
		maxLineWidth = 0
		lines        = make([]string, 0)
		sc           = bufio.NewScanner(strings.NewReader(text))
	)
	for sc.Scan() {
		line := sc.Text()
		if ncols := len(line); ncols > maxLineWidth {
			maxLineWidth = ncols
		}
		lines = append(lines, line)
	}

	return maxLineWidth, lines
}

func popupWidth(text string) (cols int) {
	tw, _ := textWidth(text)
	// +2 to returned value, to account for left + right borders
	return tw + leftPopupMargin + rightPopupMargin + 2
}

func popupHeight(text string) (rows int) {
	sc := bufio.NewScanner(strings.NewReader(text))
	for sc.Scan() {
		rows++
	}
	return rows + verticalPopupPadding
}
