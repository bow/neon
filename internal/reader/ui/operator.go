// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import "github.com/bow/neon/internal/reader/backend"

// Operator describes high-level UI operations.
type Operator interface {
	ClearStatusBar(*Display)
	FocusFeedsPane(*Display)
	FocusEntriesPane(*Display)
	FocusNextPane(*Display)
	FocusPreviousPane(*Display)
	FocusReadingPane(*Display)
	ShowAllFeeds(*Display, backend.Backend)
	ShowIntroPopup(*Display)
	ToggleAboutPopup(*Display, backend.Backend)
	ToggleHelpPopup(*Display)
	ToggleStatsPopup(*Display, backend.Backend)
	ToggleStatusBar(*Display)
	UnfocusFront(*Display)
}
