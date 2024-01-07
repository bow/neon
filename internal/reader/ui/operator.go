// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import "github.com/bow/neon/internal/reader/backend"

// Operator describes high-level UI operations.
type Operator interface {
	ClearStatusBar(d *Display)
	Draw(d *Display)
	FocusFeedsPane(d *Display)
	FocusEntriesPane(d *Display)
	FocusNextPane(d *Display)
	FocusPreviousPane(d *Display)
	FocusReadingPane(d *Display)
	HideIntroPopup(d *Display)
	NotifyInfof(d *Display, text string, a ...any)
	NotifyErr(d *Display, err error)
	NotifyErrf(d *Display, text string, a ...any)
	NotifyWarnf(d *Display, text string, a ...any)
	ShowIntroPopup(d *Display)
	ToggleAboutPopup(d *Display, b backend.Backend)
	ToggleFeedsInPane(d *Display, b backend.Backend)
	ToggleHelpPopup(d *Display)
	ToggleStatsPopup(d *Display, b backend.Backend)
	ToggleStatusBar(d *Display)
	UnfocusFront(d *Display)
}
