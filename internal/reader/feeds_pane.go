// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/bow/lens/internal"
)

type feedsPane struct {
	*tview.TreeView

	theme *Theme

	groupOrder []*tview.TreeNode
	groupNodes map[feedUpdatedGroup]*tview.TreeNode
	feedNodes  map[string]*tview.TreeNode

	feeds <-chan *internal.Feed
}

func newFeedsPane(theme *Theme, feeds <-chan *internal.Feed) *feedsPane {

	fp := feedsPane{
		theme:      theme,
		feeds:      feeds,
		groupOrder: make([]*tview.TreeNode, 0),
		groupNodes: make(map[feedUpdatedGroup]*tview.TreeNode),
		feedNodes:  make(map[string]*tview.TreeNode),
	}

	fp.initTree()

	focusf, unfocusf := fp.makeDrawFuncs()
	fp.SetDrawFunc(unfocusf)
	fp.SetFocusFunc(func() { fp.SetDrawFunc(focusf) })
	fp.SetBlurFunc(func() { fp.SetDrawFunc(unfocusf) })

	go fp.listenForUpdates()

	return &fp
}

// TODO: How to handle feeds being removed altogether?
func (fp *feedsPane) listenForUpdates() {
	root := fp.GetRoot()
	for feed := range fp.feeds {
		fnode, exists := fp.feedNodes[feed.FeedURL]
		newGroup := whenUpdated(feed)
		if exists {
			oldGroup := whenUpdated(fnode.GetReference().(*internal.Feed))
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
		SetCurrentNode(root).
		SetTopLevel(1)

	fp.TreeView = tree

	for i := uint8(0); i < uint8(updatedUnknown); i++ {
		ug := feedUpdatedGroup(i)
		gnode := groupNode(ug, fp.theme)
		fp.groupNodes[ug] = gnode
		fp.groupOrder = append(fp.groupOrder, gnode)
	}
}

func (fp *feedsPane) makeDrawFuncs() (focusf, unfocusf drawFunc) {

	var titleUF, titleF string
	if fp.theme.FeedsPaneTitle != "" {
		titleUF = fmt.Sprintf(" %s ", fp.theme.FeedsPaneTitle)
		titleF = fmt.Sprintf("[::b]▶ %s[::-] ", fp.theme.FeedsPaneTitle)
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
				fp.theme.TitleForeground,
			)

			return x + 2, y + 1, width - 2, height - 1
		}
	}

	focusf = drawf(true)
	unfocusf = drawf(false)

	return focusf, unfocusf
}

func (fp *feedsPane) refreshColors() {
	for _, node := range fp.TreeView.GetRoot().GetChildren() {
		node.SetColor(fp.theme.FeedsGroup)
	}
}

type feedUpdatedGroup uint8

const (
	updatedToday feedUpdatedGroup = iota
	updatedThisWeek
	updatedThisMonth
	updatedEarlier
	updatedUnknown
)

func (ug feedUpdatedGroup) Text(theme *Theme) string {
	switch ug {
	case updatedToday:
		return theme.UpdatedTodayText
	case updatedThisWeek:
		return theme.UpdatedThisWeekText
	case updatedThisMonth:
		return theme.UpdatedThisMonthText
	case updatedEarlier:
		return theme.UpdatedEarlier
	case updatedUnknown:
		return theme.UpdatedUnknownText
	default:
		return theme.UpdatedUnknownText
	}
}

func feedNode(feed *internal.Feed, _ *Theme) *tview.TreeNode {
	return tview.NewTreeNode(feed.Title).
		SetReference(feed.FeedURL).
		SetColor(tcell.ColorWhite).
		SetSelectable(true)
}

func groupNode(ug feedUpdatedGroup, theme *Theme) *tview.TreeNode {
	return tview.NewTreeNode(ug.Text(theme)).
		SetReference(ug).
		SetColor(theme.FeedsGroup).
		SetSelectable(false)
}

func whenUpdated(feed *internal.Feed) feedUpdatedGroup {
	if feed.Updated == nil {
		return updatedUnknown
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -7)
	lastMonth := now.AddDate(0, -1, 0)

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
