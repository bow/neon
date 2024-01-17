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
// TODO: How to maintain ordering of most-recently updated feeds first?
func (fp *feedsPane) startFeedsPoll(ch <-chan *entity.Feed) {
	root := fp.GetRoot()
	for feed := range ch {
		fnode, exists := fp.feedNodes[feed.FeedURL]
		newGroup := whenUpdated(feed)
		targetGNode := fp.groupNodes[newGroup]

		if !exists {
			fnode = feedNode(feed, fp.theme)
			fp.feedNodes[feed.FeedURL] = fnode
			targetGNode.AddChild(fnode)
		} else {
			oldFeed := fnode.GetReference().(*entity.Feed)
			oldGroup := whenUpdated(oldFeed)
			// TODO: Combine entries for existing feeds, set fnode reference, add to target.
			if oldGroup != newGroup {
				existingGNode := fp.groupNodes[oldGroup]
				existingGNode.RemoveChild(fnode)
				if len(existingGNode.GetChildren()) == 0 {
					root.RemoveChild(existingGNode)
				}
			}
			setFeedNodeDisplay(fnode, fp.theme)
		}
		if len(targetGNode.GetChildren()) == 1 && !exists {
			root.AddChild(targetGNode)
		}
	}
}

func (fp *feedsPane) initTree() {

	root := tview.NewTreeNode("")

	tree := tview.NewTreeView().
		SetRoot(root).
		SetGraphics(false).
		SetCurrentNode(root).
		SetPrefixes([]string{"  ", "· "}).
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

			return x, y + 1, width - 2, height - 1
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
	node := tview.NewTreeNode("").
		SetReference(feed).
		SetSelectable(true)
	setFeedNodeDisplay(node, theme)
	return node
}

func setFeedNodeDisplay(fnode *tview.TreeNode, theme *Theme) {
	feed := fnode.GetReference().(*entity.Feed)
	if c := feed.NumEntriesUnread(); c > 0 {
		fnode.SetText(fmt.Sprintf("%s (%d)", feed.Title, c)).
			SetColor(tcell.ColorYellow)
	} else {
		fnode.SetText(feed.Title).
			SetColor(theme.feedNode)
	}
}

func groupNode(ug feedUpdatePeriod, theme *Theme, lang *Lang) *tview.TreeNode {
	return tview.NewTreeNode(ug.Text(lang)).
		SetReference(ug).
		SetColor(theme.feedGroupNode).
		SetSelectable(false)
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
