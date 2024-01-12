// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

type Lang struct {
	feedsPaneTitle   string
	entriesPaneTitle string
	readingPaneTitle string

	aboutPopupTitle string
	helpPopupTitle  string
	statsPopupTitle string
	introPopupTitle string

	updatedTodayText     string
	updatedThisWeekText  string
	updatedThisMonthText string
	updatedEarlierText   string
	updatedUnknownText   string
}

var langEN = &Lang{
	feedsPaneTitle:   "Feeds",
	entriesPaneTitle: "Entries",
	readingPaneTitle: "",

	aboutPopupTitle: "About",
	helpPopupTitle:  "Keys",
	statsPopupTitle: "Stats",
	introPopupTitle: "Welcome",

	updatedTodayText:     "Today",
	updatedThisWeekText:  "This Week",
	updatedThisMonthText: "This Month",
	updatedEarlierText:   "Earlier",
	updatedUnknownText:   "Unknown",
}
