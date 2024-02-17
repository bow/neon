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

	theme    *Theme
	lang     *Lang
	incoming <-chan *entity.Feed

	store *feedStore
}

func newFeedsPane(
	theme *Theme,
	lang *Lang,
	incoming <-chan *entity.Feed,
) *feedsPane {

	fp := feedsPane{
		theme:    theme,
		lang:     lang,
		incoming: incoming,
		store:    newFeedStore(),
	}

	fp.initTree()

	focusf, unfocusf := fp.makeDrawFuncs()
	fp.SetDrawFunc(unfocusf)
	fp.SetFocusFunc(
		func() {
			fp.SetDrawFunc(focusf)
			current := fp.GetCurrentNode()
			if current == nil || current == fp.GetRoot() {
				if target := fp.getFirstFeedNode(); target != nil {
					fp.SetCurrentNode(target)
				}
			}
		},
	)
	fp.SetBlurFunc(func() { fp.SetDrawFunc(unfocusf) })

	return &fp
}

func (fp *feedsPane) startPoll() (stop func()) {
	done := make(chan struct{})
	stop = func() {
		defer close(done)
		done <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case feed := <-fp.incoming:
				fp.updateFeed(feed)
			}
		}
	}()

	return stop
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
		SetPrefixes([]string{"  ", "· "}).
		SetTopLevel(1)

	fp.TreeView = *tree
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
			if group := periodOf(gnode); group != nil && targetGroup == *group {
				return gnode
			}
		}
	}
	return nil
}

func (fp *feedsPane) getCurrentFeed() *entity.Feed {
	return feedOf(fp.GetCurrentNode())
}

func (fp *feedsPane) getFoldState() foldState {
	root := fp.GetRoot()
	if root == nil {
		return foldUnknown
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
	if !allExpanded && !allCollapsed {
		return foldMixed
	}
	if allCollapsed && !allExpanded {
		return foldAllCollapsed
	}
	if allExpanded && !allCollapsed {
		return foldAllExpanded
	}
	panic("impossible fold state")
}

func (fp *feedsPane) toggleAllFeedsFold() {
	root := fp.GetRoot()
	if root == nil {
		return
	}

	switch state := fp.getFoldState(); state {

	case foldMixed, foldAllCollapsed:
		for _, gnode := range root.GetChildren() {
			if period := periodOf(gnode); period != nil {
				gnode.SetText(period.Text(fp.lang))
			}
			gnode.Expand()
		}
		return

	case foldAllExpanded:
		current := fp.getCurrentGroupNode()
		for _, gnode := range root.GetChildren() {
			if unread := countGroupUnread(gnode); unread > 0 {
				if period := periodOf(gnode); period != nil {
					gnode.SetText(fmt.Sprintf("%s (%d)", period.Text(fp.lang), unread))
				}
			}
			gnode.Collapse()
		}
		// Set selection to nearest group prior to collapsing.
		fp.SetCurrentNode(current)
		return

	case foldUnknown:
		return
	}
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
			if unread := countGroupUnread(gnode); unread > 0 {
				if period := periodOf(gnode); period != nil {
					gnode.SetText(fmt.Sprintf("%s (%d)", period.Text(fp.lang), unread))
				}
			}
			gnode.Collapse()
		} else {
			if period := periodOf(gnode); period != nil {
				gnode.SetText(period.Text(fp.lang))
			}
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

type foldState uint8

const (
	foldUnknown foldState = iota
	foldMixed
	foldAllCollapsed
	foldAllExpanded
)

type feedUpdatePeriod uint8

const (
	updatedToday feedUpdatePeriod = iota
	updatedThisWeek
	updatedThisMonth
	updatedEarlier
	updatedUnknown
)

func (period feedUpdatePeriod) Text(lang *Lang) string {
	switch period {
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
	feed := feedOf(fnode)
	if feed == nil {
		return
	}
	if c := feed.NumEntriesUnread(); c > 0 {
		fnode.SetText(fmt.Sprintf("%s (%d)", feed.Title, c)).
			SetColor(theme.feedNodeUnread)
	} else {
		fnode.SetText(feed.Title).
			SetColor(theme.feedNode)
	}
}

func groupNode(period feedUpdatePeriod, theme *Theme, lang *Lang) *tview.TreeNode {
	return tview.NewTreeNode(period.Text(lang)).
		SetReference(period).
		SetColor(theme.feedGroupNode).
		SetSelectable(true)
}

func countGroupUnread(gnode *tview.TreeNode) int {
	var unread int
	period := periodOf(gnode)
	if period == nil {
		return 0
	}
	for _, fnode := range gnode.GetChildren() {
		feed := feedOf(fnode)
		if feed == nil {
			continue
		}
		unread += feed.NumEntriesUnread()
	}
	return unread
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

func feedOf(node *tview.TreeNode) *entity.Feed {
	if node == nil {
		return nil
	}
	feed, ok := node.GetReference().(*entity.Feed)
	if !ok {
		return nil
	}
	return feed
}

func periodOf(node *tview.TreeNode) *feedUpdatePeriod {
	if node == nil {
		return nil
	}
	period, ok := node.GetReference().(feedUpdatePeriod)
	if !ok {
		return nil
	}
	return &period
}
