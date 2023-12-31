// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"bufio"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/neon/internal"
)

const (
	leftPopupMargin      = 2
	rightPopupMargin     = 2
	verticalPopupPadding = 4
)

type popup struct {
	content tview.Primitive
	frame   *tview.Frame
	grid    *tview.Grid
}

func (p *popup) setContent(prim tview.Primitive) {
	p.frame.SetPrimitive(prim)
	p.content = prim
}

func (p *popup) setWidth(w int) {
	p.grid.SetColumns(0, w, 0)
}

func (p *popup) setGridRows(rows []int) {
	p.grid.SetRows(rows...)
}

func newEmptyPopup(
	title string,
	titleColorFG tcell.Color,
	top, bottom int,
) *popup {
	var content tview.Primitive = nil

	frame := tview.NewFrame(content).
		SetBorders(top, bottom, 0, 0, leftPopupMargin, rightPopupMargin)

	frame.SetBorder(true).
		SetTitle(fmt.Sprintf(" %s ", title)).
		SetTitleColor(titleColorFG)

	grid := tview.NewGrid().
		AddItem(frame, 1, 1, 1, 1, 0, 0, true)

	p := popup{grid: grid, frame: frame, content: content}

	return &p
}

func newPopup(
	title string,
	content tview.Primitive,
	titleColorFG tcell.Color,
	top, bottom int,
	width int,
	gridRows []int,
) *popup {

	frame := tview.NewFrame(content).
		SetBorders(top, bottom, 0, 0, leftPopupMargin, rightPopupMargin)

	frame.SetBorder(true).
		SetTitle(fmt.Sprintf(" %s ", title)).
		SetTitleColor(titleColorFG)

	grid := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(gridRows...).
		AddItem(frame, 1, 1, 1, 1, 0, 0, true)

	p := popup{grid: grid, frame: frame, content: content}

	return &p
}

func setAboutPopupText(p *popup, backend string) {
	commit := internal.GitCommit()

	var buildTime = internal.BuildTime()
	buildTimeVal, err := time.Parse(time.RFC3339, buildTime)
	if err == nil {
		buildTime = buildTimeVal.Format(longDateFormat)
	}

	infoText := fmt.Sprintf(`[yellow]Version[-]   : %s
[yellow]Git commit[-]: %s
[yellow]Build time[-]: %s
[yellow]Backend[-]   : %s`,
		internal.Version(),
		commit,
		buildTime,
		backend,
	)

	aboutWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetText(infoText)

	// NOTE: We assume the banner's width is less than the one computed here.
	width := popupWidth(aboutWidget.GetText(true))
	banner := centerBanner(internal.Banner(), width)
	aboutText := fmt.Sprintf("%s\n\n%s", banner, infoText)

	aboutWidget.SetText(aboutText)

	height := popupHeight(aboutText) - 1

	p.setWidth(width)
	p.setGridRows([]int{-1, height, -3})
	p.setContent(aboutWidget)
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
