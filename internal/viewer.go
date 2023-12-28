// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"github.com/gdamore/tcell/v2"

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
