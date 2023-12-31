// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import "github.com/bow/neon/internal/entity"

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
	Start(d *Display) error
	Stop(d *Display)
	ToggleAboutPopup(d *Display, backend string)
	ToggleFeedsInPane(d *Display, feeds <-chan *entity.Feed)
	ToggleHelpPopup(d *Display)
	ToggleIntroPopup(d *Display)
	ToggleStatsPopup(d *Display, stats <-chan *entity.Stats)
	ToggleStatusBar(d *Display)
	UnfocusPane(d *Display)
}
