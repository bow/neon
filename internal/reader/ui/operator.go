// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import "github.com/bow/neon/internal/reader/repo"

// Operator describes high-level UI operations.
type Operator interface {
	ClearStatusBar(d *Display)
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
	ToggleAboutPopup(d *Display, r repo.Repo)
	ToggleFeedsInPane(d *Display, r repo.Repo)
	ToggleHelpPopup(d *Display)
	ToggleIntroPopup(d *Display)
	ToggleStatsPopup(d *Display, r repo.Repo)
	ToggleStatusBar(d *Display)
	UnfocusPane(d *Display)
}
