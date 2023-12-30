// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Display struct {
	inner  *tview.Application
	screen tcell.Screen
	theme  *Theme
}

func NewDisplay(screen tcell.Screen, theme string) (*Display, error) {
	th, err := LoadTheme(theme)
	if err != nil {
		return nil, err
	}

	root := makeRoot(th)
	inner := tview.NewApplication().
		EnableMouse(true).
		SetRoot(root, true).
		SetScreen(screen)

	return &Display{inner: inner, screen: screen, theme: th}, nil
}

func (l *Display) Start() error {
	return l.inner.Run()
}

func makeRoot(theme *Theme) tview.Primitive {
	// TODO Replace placeholder box with actual page.
	return tview.NewBox().SetBorder(true)
}
