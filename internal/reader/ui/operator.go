// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import "github.com/bow/neon/internal/entity"

// Operator describes high-level UI operations.
type Operator interface {
	ClearStatusBar(*Display)
	FocusFeedsPane(*Display)
	FocusEntriesPane(*Display)
	FocusNextPane(*Display)
	FocusPreviousPane(*Display)
	FocusReadingPane(*Display)
	RefreshFeeds(*Display, func() (<-chan entity.PullResult, error))
	RefreshStats(*Display, func() (*entity.Stats, error))
	PopulateFeedsPane(*Display, func() ([]*entity.Feed, error))
	ShowIntroPopup(*Display)
	ToggleAboutPopup(*Display, string)
	ToggleAllFeedsFold(*Display)
	ToggleCurrentFeedFold(*Display)
	ToggleHelpPopup(*Display)
	ToggleStatsPopup(*Display, func() (*entity.Stats, error))
	ToggleStatusBar(*Display)
	UnfocusFront(*Display)
}
