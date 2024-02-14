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

	theme      *Theme
	lang       *Lang
	store      *feedStore
	focusStack *entity.ID
}

func newFeedsPane(theme *Theme, lang *Lang) *feedsPane {

	fp := feedsPane{
		theme: theme,
		lang:  lang,
		store: newFeedStore(),
	}

	fp.initTree()

	focusf, unfocusf := fp.makeDrawFuncs()
	fp.SetDrawFunc(unfocusf)
	fp.SetFocusFunc(
		func() {
			fp.SetDrawFunc(focusf)

			if previous := fp.findFeedNode(fp.focusStack); previous != nil {
				fp.SetCurrentNode(previous)
			} else if fallback := fp.getFirstFeedNode(); fallback != nil {
				fp.SetCurrentNode(fallback)
			}
		},
	)
	fp.SetBlurFunc(
		func() {
			fp.SetDrawFunc(unfocusf)
			if feed := fp.getCurrentFeed(); feed != nil {
				fp.focusStack = &feed.ID
			}
			fp.SetCurrentNode(nil)
		},
	)

	return &fp
}

func (fp *feedsPane) updateFeed(feed *entity.Feed) {
	root := fp.GetRoot()

	fp.store.upsert(feed)

	var currentFeedID *entity.ID
	if currentFeed := fp.getCurrentFeed(); currentFeed != nil {
		currentFeedID = &currentFeed.ID
	}

	root.ClearChildren()

	for _, group := range fp.store.feedsByPeriod() {
		gnode := groupNode(group.label, fp.theme, fp.lang)
		root.AddChild(gnode)

		for _, feed := range group.feedsSlice() {
			fnode := feedNode(feed, fp.theme)
			setFeedNodeDisplay(fnode, fp.theme)
			gnode.AddChild(fnode)
			if currentFeedID != nil && feed.ID == *currentFeedID {
				fp.SetCurrentNode(fnode)
			}
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
}

func (fp *feedsPane) findFeedNode(id *entity.ID) *tview.TreeNode {
	if id == nil {
		return nil
	}
	root := fp.GetRoot()
	if root == nil {
		return nil
	}
	for _, gnode := range root.GetChildren() {
		for _, fnode := range gnode.GetChildren() {
			feed := fnode.GetReference().(*entity.Feed)
			if feed.ID == *id {
				return fnode
			}
		}
	}
	return nil
}

func (fp *feedsPane) getFirstFeedNode() *tview.TreeNode {
	root := fp.GetRoot()
	if root == nil {
		return nil
	}
	for i, gnode := range root.GetChildren() {
		for j, fnode := range gnode.GetChildren() {
			if i == 0 && j == 0 {
				return fnode
			}
		}
	}
	return nil
}

func (fp *feedsPane) getCurrentGroupNode() *tview.TreeNode {
	root := fp.GetRoot()
	if root == nil {
		return nil
	}
	current := fp.GetCurrentNode()
	if current == nil {
		return nil
	}
	switch t := current.GetReference().(type) {
	case feedUpdatePeriod:
		return current
	case *entity.Feed:
		targetGroup := whenUpdated(t)
		for _, gnode := range root.GetChildren() {
			if group, ok := gnode.GetReference().(feedUpdatePeriod); ok && group == targetGroup {
				return gnode
			}
		}
	}
	return nil
}

func (fp *feedsPane) getCurrentFeed() *entity.Feed {
	if fnode := fp.GetCurrentNode(); fnode != nil {
		feed, ok := fnode.GetReference().(*entity.Feed)
		if ok {
			return feed
		}
	}
	return nil
}

func (fp *feedsPane) toggleAllFeedsFold() {
	root := fp.GetRoot()
	if root == nil {
		return
	}
	var allExpanded, allCollapsed bool
	for i, gnode := range root.GetChildren() {
		expanded := gnode.IsExpanded()
		if i == 0 {
			allExpanded = expanded
			allCollapsed = !expanded
			continue
		}
		allExpanded = allExpanded && expanded
		allCollapsed = allCollapsed && !expanded
	}
	// treat mixed state as if everything is collapsed. that is, we then expand everything.
	if (!allExpanded && !allCollapsed) || (allCollapsed && !allExpanded) {
		for _, gnode := range root.GetChildren() {
			gnode.Expand()
		}
		return
	} else if allExpanded && !allCollapsed {
		current := fp.getCurrentGroupNode()
		for _, gnode := range root.GetChildren() {
			gnode.Collapse()
		}
		// Set selection to nearest group prior to collapsing.
		fp.SetCurrentNode(current)
		return
	}
	panic("impossible fold state")
}

func (fp *feedsPane) toggleCurrentFeedFold() {
	root := fp.GetRoot()
	if root == nil {
		return
	}
	current := fp.GetCurrentNode()
	if current == nil {
		return
	}
	if gnode := fp.getCurrentGroupNode(); gnode != nil {
		if gnode.IsExpanded() {
			gnode.Collapse()
		} else {
			gnode.Expand()
		}
		fp.SetCurrentNode(gnode)
		return
	}
}

func (fp *feedsPane) makeDrawFuncs() (focusf, unfocusf drawFunc) {

	titleUF, titleF := fmtPaneTitle(fp.lang.feedsPaneTitle)

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
	for _, gnode := range fp.GetRoot().GetChildren() {
		gnode.SetColor(fp.theme.feedGroupNode)
		for _, fnode := range gnode.GetChildren() {
			setFeedNodeDisplay(fnode, fp.theme)
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
			SetColor(theme.feedNodeUnread)
	} else {
		fnode.SetText(feed.Title).
			SetColor(theme.feedNode)
	}
}

func groupNode(ug feedUpdatePeriod, theme *Theme, lang *Lang) *tview.TreeNode {
	node := tview.NewTreeNode(ug.Text(lang)).
		SetReference(ug).
		SetColor(theme.feedGroupNode).
		SetSelectable(true)

	node.SetSelectedFunc(func() {
		if node.IsExpanded() {
			node.Collapse()
		} else {
			node.Expand()
		}
	})

	return node
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
