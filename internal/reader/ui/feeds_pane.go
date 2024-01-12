// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/neon/internal/entity"
)

type feedsPane struct {
	tview.TreeView

	theme *Theme
	lang  *Lang

	groupOrder []*tview.TreeNode
	groupNodes map[feedUpdatePeriod]*tview.TreeNode
	feedNodes  map[string]*tview.TreeNode
}

func newFeedsPane(theme *Theme, lang *Lang) *feedsPane {

	fp := feedsPane{
		theme:      theme,
		lang:       lang,
		groupOrder: make([]*tview.TreeNode, 0),
		groupNodes: make(map[feedUpdatePeriod]*tview.TreeNode),
		feedNodes:  make(map[string]*tview.TreeNode),
	}

	fp.initTree()

	focusf, unfocusf := fp.makeDrawFuncs()
	fp.SetDrawFunc(unfocusf)
	fp.SetFocusFunc(func() { fp.SetDrawFunc(focusf) })
	fp.SetBlurFunc(func() { fp.SetDrawFunc(unfocusf) })

	return &fp
}

// TODO: How to handle feeds being removed altogether?
func (fp *feedsPane) startFeedsPoll(ch <-chan *entity.Feed) {
	root := fp.GetRoot()
	for feed := range ch {
		fnode, exists := fp.feedNodes[feed.FeedURL]
		newGroup := whenUpdated(feed)
		if exists {
			oldGroup := whenUpdated(fnode.GetReference().(*entity.Feed))
			if oldGroup != newGroup {
				fp.groupNodes[oldGroup].RemoveChild(fnode)
			}
		} else {
			fnode = feedNode(feed, fp.theme)
			fp.feedNodes[feed.FeedURL] = fnode
		}
		fp.groupNodes[newGroup].AddChild(fnode)

		root.ClearChildren()
		for _, gnode := range fp.groupOrder {
			if len(gnode.GetChildren()) > 0 {
				root.AddChild(gnode)
			}
		}
	}
}

func (fp *feedsPane) initTree() {

	root := tview.NewTreeNode("")

	tree := tview.NewTreeView().
		SetRoot(root).
		SetGraphics(false).
		SetPrefixes([]string{"", "· "}).
		SetCurrentNode(root).
		SetTopLevel(1)

	fp.TreeView = *tree

	for i := uint8(0); i < uint8(updatedUnknown); i++ {
		ug := feedUpdatePeriod(i)
		gnode := groupNode(ug, fp.theme, fp.lang)
		fp.groupNodes[ug] = gnode
		fp.groupOrder = append(fp.groupOrder, gnode)
	}
}

func (fp *feedsPane) makeDrawFuncs() (focusf, unfocusf drawFunc) {

	var titleUF, titleF string
	if fp.lang.feedsPaneTitle != "" {
		titleUF = fmt.Sprintf(" %s ", fp.lang.feedsPaneTitle)
		titleF = fmt.Sprintf("[::b]▶ %s[::-] ", fp.lang.feedsPaneTitle)
	} else {
		titleF = "[::b]▶[::-] "
	}

	drawf := func(
		focused bool,
	) func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {

		var (
			title   string
			leftPad int
		)
		if focused {
			title = titleF
			leftPad = 0
		} else {
			title = titleUF
			leftPad = 1
		}

		return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			lineStyle := fp.theme.lineStyle()
			// Draw top and optionally bottom borders.
			for cx := x; cx < x+width; cx++ {
				screen.SetContent(cx, y, tview.BoxDrawingsLightHorizontal, nil, lineStyle)
			}

			// Write the title text.
			tview.Print(
				screen,
				title,
				x+leftPad,
				y,
				width-2,
				tview.AlignLeft,
				fp.theme.titleFG,
			)

			return x + 2, y + 1, width - 2, height - 1
		}
	}

	focusf = drawf(true)
	unfocusf = drawf(false)

	return focusf, unfocusf
}

func (fp *feedsPane) refreshColors() {
	for _, gnode := range fp.TreeView.GetRoot().GetChildren() {
		gnode.SetColor(fp.theme.feedGroupNode)
		for _, fnode := range gnode.GetChildren() {
			fnode.SetColor(fp.theme.feedNode)
		}
	}
}

type feedUpdatePeriod uint8

const (
	updatedToday feedUpdatePeriod = iota
	updatedThisWeek
	updatedThisMonth
	updatedEarlier
	updatedUnknown
)

func (ug feedUpdatePeriod) Text(lang *Lang) string {
	switch ug {
	case updatedToday:
		return lang.updatedTodayText
	case updatedThisWeek:
		return lang.updatedThisWeekText
	case updatedThisMonth:
		return lang.updatedThisMonthText
	case updatedEarlier:
		return lang.updatedEarlierText
	case updatedUnknown:
		return lang.updatedUnknownText
	default:
		return lang.updatedUnknownText
	}
}

func feedNode(feed *entity.Feed, theme *Theme) *tview.TreeNode {
	return tview.NewTreeNode(feed.Title).
		SetReference(feed.FeedURL).
		SetColor(theme.feedNode).
		SetSelectable(true)
}

func groupNode(ug feedUpdatePeriod, theme *Theme, lang *Lang) *tview.TreeNode {
	return tview.NewTreeNode(ug.Text(lang)).
		SetReference(ug).
		SetColor(theme.feedGroupNode).
		SetSelectable(true)
}

func whenUpdated(feed *entity.Feed) feedUpdatePeriod {
	if feed.Updated == nil {
		return updatedUnknown
	}

	var (
		now       = time.Now()
		yesterday = now.AddDate(0, 0, -1)
		lastWeek  = now.AddDate(0, 0, -7)
		lastMonth = now.AddDate(0, -1, 0)
	)

	ft := *feed.Updated
	switch {
	case ft.Before(lastMonth):
		return updatedEarlier
	case ft.Before(lastWeek):
		return updatedThisMonth
	case ft.Before(yesterday):
		return updatedThisWeek
	default:
		return updatedToday
	}
}
