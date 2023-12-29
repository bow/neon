// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/neon/internal/entity"
)

// Viewer describes the console reader.
type Viewer interface {
	ClearStatusBar()
	FocusFeedsPane()
	FocusEntriesPane()
	FocusNextPane()
	FocusPreviousPane()
	FocusReadingPane()
	HideIntroPopup()
	NotifyInfof(text string, a ...any)
	NotifyErr(err error)
	NotifyErrf(text string, a ...any)
	NotifyWarnf(text string, a ...any)
	SetFeedsPaneKeyHandler(handler func(*tcell.EventKey) *tcell.EventKey)
	SetGlobalKeyHandler(handler func(*tcell.EventKey) *tcell.EventKey)
	ShowAboutPopup()
	ShowFeedsInPane(<-chan *entity.Feed)
	ShowHelpPopup()
	ShowIntroPopup()
	ShowStatsPopup()
	Start() error
	Stop()
	ToggleStatusBar()
	UnfocusPane()
}

//nolint:unused
type View struct {
	screen tcell.Screen
	theme  *Theme
	lang   *Lang

	app  *tview.Application
	root *tview.Pages

	focusStack tview.Primitive
}

//nolint:revive
func New(scr tcell.Screen, theme string) (*View, error) {
	panic("New is unimplemented")
}

func (v *View) ClearStatusBar() {
	panic("ClearStatusBar is unimplemented")
}

func (v *View) FocusFeedsPane() {
	panic("FocusFeedsPane is unimplemented")
}

func (v *View) FocusEntriesPane() {
	panic("FocusEntriesPane is unimplemented")
}

func (v *View) FocusNextPane() {
	panic("FocusNextPane is unimplemented")
}

func (v *View) FocusPreviousPane() {
	panic("FocusPreviousPane is unimplemented")
}

func (v *View) FocusReadingPane() {
	panic("FocusReadingPane is unimplemented")
}

func (v *View) HideIntroPopup() {
	panic("HideIntroPopup is unimplemented")
}

//nolint:revive
func (v *View) NotifyInfof(text string, a ...any) {
	panic("NotifyInfof is unimplemented")
}

//nolint:revive
func (v *View) NotifyErr(err error) {
	panic("NotifyErr is unimplemented")
}

//nolint:revive
func (v *View) NotifyErrf(text string, a ...any) {
	panic("NotifyErrf is unimplemented")
}

//nolint:revive
func (v *View) NotifyWarnf(text string, a ...any) {
	panic("NotifyWarnf is unimplemented")
}

//nolint:revive
func (v *View) SetFeedsPaneKeyHandler(handler func(*tcell.EventKey) *tcell.EventKey) {
	panic("SetFeedsPaneKeyHandler(handler func is unimplemented")
}

//nolint:revive
func (v *View) SetGlobalKeyHandler(handler func(*tcell.EventKey) *tcell.EventKey) {
	panic("SetGlobalKeyHandler(handler func is unimplemented")
}

func (v *View) ShowAboutPopup() {
	panic("ShowAboutPopup is unimplemented")
}

//nolint:revive
func (v *View) ShowFeedsInPane(<-chan *entity.Feed) {
	panic("ShowFeedsInPane is unimplemented")
}

func (v *View) ShowHelpPopup() {
	panic("ShowHelpPopup is unimplemented")
}

func (v *View) ShowIntroPopup() {
	panic("ShowIntroPopup is unimplemented")
}

func (v *View) ShowStatsPopup() {
	panic("ShowStatsPopup is unimplemented")
}

func (v *View) Start() error {
	panic("Start is unimplemented")
}

func (v *View) Stop() {
	panic("Stop is unimplemented")
}

func (v *View) ToggleStatusBar() {
	panic("ToggleStatusBar is unimplemented")
}

func (v *View) UnfocusPane() {
	panic("UnfocusPane is unimplemented")
}

// Ensure View implements Viewer.
var _ Viewer = new(View)

//nolint:unused
type drawFunc func(screen tcell.Screen, x int, y int, w int, h int) (ix int, iy int, iw int, ih int)
