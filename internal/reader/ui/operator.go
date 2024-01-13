// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import "github.com/bow/neon/internal/reader/backend"

// Operator describes high-level UI operations.
type Operator interface {
	ClearStatusBar(d *Display)
	FocusFeedsPane(d *Display)
	FocusEntriesPane(d *Display)
	FocusNextPane(d *Display)
	FocusPreviousPane(d *Display)
	FocusReadingPane(d *Display)
	ShowFeedsInPane(d *Display, b backend.Backend)
	ShowIntroPopup(d *Display)
	ToggleAboutPopup(d *Display, b backend.Backend)
	ToggleHelpPopup(d *Display)
	ToggleStatsPopup(d *Display, b backend.Backend)
	ToggleStatusBar(d *Display)
	UnfocusFront(d *Display)
}
