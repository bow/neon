// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"google.golang.org/grpc"

	rp "github.com/bow/neon/internal/reader/repo"
	"github.com/bow/neon/internal/reader/ui"
)

//nolint:unused
type Reader struct {
	ctx      context.Context
	initPath string

	dsp  *ui.Display
	opr  ui.Operator
	repo rp.Repo

	stopped bool
}

func (r *Reader) Start() error {
	return r.dsp.Start()
}

type Builder struct {
	ctx       context.Context
	themeName string
	initPath  string
	scr       tcell.Screen

	// rpcRepo args.
	addr  string
	dopts []grpc.DialOption

	// For testing.
	rpo rp.Repo
	opr ui.Operator
}

func NewBuilder() *Builder {
	b := Builder{
		ctx:       context.Background(),
		themeName: "dark",
		dopts:     nil,
	}
	return &b
}

func (b *Builder) Address(addr string) *Builder {
	b.addr = addr
	return b
}

func (b *Builder) DialOpts(dialOpts ...grpc.DialOption) *Builder {
	b.dopts = dialOpts
	return b
}

func (b *Builder) Context(ctx context.Context) *Builder {
	b.ctx = ctx
	return b
}

func (b *Builder) InitPath(path string) *Builder {
	b.initPath = path
	return b
}

func (b *Builder) Theme(name string) *Builder {
	b.themeName = name
	return b
}

func (b *Builder) repo(rpo rp.Repo) *Builder {
	b.rpo = rpo
	return b
}

func (b *Builder) screen(scr tcell.Screen) *Builder {
	b.scr = scr
	return b
}

func (b *Builder) operator(opr ui.Operator) *Builder {
	b.opr = opr
	return b
}

func (b *Builder) Build() (*Reader, error) {

	if b.addr == "" && b.rpo == nil {
		return nil, fmt.Errorf("reader server address must be specified")
	}

	var (
		rpo rp.Repo
		err error
	)
	if b.rpo != nil {
		rpo = b.rpo
	} else {
		rpo, err = rp.NewRPCRepo(b.ctx, b.addr, b.dopts...)
		if err != nil {
			return nil, err
		}
	}

	var scr tcell.Screen
	if b.scr != nil {
		scr = b.scr
	} else {
		scr, err = tcell.NewScreen()
		if err != nil {
			return nil, err
		}
	}
	dsp, err := ui.NewDisplay(scr, b.themeName)
	if err != nil {
		return nil, err
	}

	var opr ui.Operator
	if b.opr != nil {
		opr = b.opr
	} else {
		opr = ui.NewDisplayOperator()
	}

	rdr := Reader{
		ctx:  b.ctx,
		dsp:  dsp,
		opr:  opr,
		repo: rpo,

		stopped: false,
	}
	rdr.dsp.Init(rdr.globalKeyHandler())

	return &rdr, nil
}

func (r *Reader) globalKeyHandler() ui.KeyHandler {

	return func(event *tcell.EventKey) *tcell.EventKey {
		var (
			key  = event.Key()
			keyr = event.Rune()
		)

		// nolint:gocritic,revive,exhaustive
		switch key {

		case tcell.KeyRune:
			switch keyr {
			case 'h', '?':
				r.opr.ToggleHelpPopup(r.dsp)
				return nil

			case 'q':
				r.dsp.Stop()
				r.stopped = true
				return nil
			}
		}

		return event
	}
}
