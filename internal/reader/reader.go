// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"google.golang.org/grpc"

	bknd "github.com/bow/neon/internal/reader/backend"
	"github.com/bow/neon/internal/reader/ui"
)

//nolint:unused
type Reader struct {
	ctx      context.Context
	initPath string

	dsp *ui.Display
	opr ui.Operator
	be  bknd.Backend
}

func (r *Reader) Start() error {
	return r.dsp.Start()
}

func (r *Reader) globalKeyHandler() ui.KeyHandler {
	r.mustDefinedFields()

	return func(event *tcell.EventKey) *tcell.EventKey {
		var (
			key  = event.Key()
			keyr = event.Rune()
		)

		// nolint:exhaustive
		switch key {

		case tcell.KeyRune:
			switch keyr {
			case 'A':
				r.opr.ToggleAboutPopup(r.dsp, r.be)
				return nil

			case 'h', '?':
				r.opr.ToggleHelpPopup(r.dsp)
				return nil

			case 'q':
				r.dsp.Stop()
				return nil
			}

		case tcell.KeyEscape:
			r.opr.UnfocusFront(r.dsp)
			return nil
		}

		return event
	}
}

func (r *Reader) mustDefinedFields() {
	if r.dsp == nil {
		panic("can not set handler with nil display")
	}

	if r.opr == nil {
		panic("can not set handler with nil operator")
	}

	if r.be == nil {
		panic("can not set handler with nil backend")
	}
}

type Builder struct {
	ctx       context.Context
	themeName string
	initPath  string
	scr       tcell.Screen

	// rpcBackend args.
	addr  string
	dopts []grpc.DialOption

	// For testing.
	be  bknd.Backend
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

func (b *Builder) backend(be bknd.Backend) *Builder {
	b.be = be
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

	if b.addr == "" && b.be == nil {
		return nil, fmt.Errorf("reader server address must be specified")
	}

	var (
		be  bknd.Backend
		err error
	)
	if b.be != nil {
		be = b.be
	} else {
		be, err = bknd.NewRPC(b.ctx, b.addr, b.dopts...)
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
		ctx: b.ctx,
		dsp: dsp,
		opr: opr,
		be:  be,
	}
	rdr.dsp.SetHandlers(rdr.globalKeyHandler())

	return &rdr, nil
}
