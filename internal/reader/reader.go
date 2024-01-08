// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"google.golang.org/grpc"

	bknd "github.com/bow/neon/internal/reader/backend"
	st "github.com/bow/neon/internal/reader/state"
	"github.com/bow/neon/internal/reader/ui"
)

//nolint:unused
type Reader struct {
	ctx context.Context

	display *ui.Display
	opr     ui.Operator
	backend bknd.Backend
	state   st.State
}

func (r *Reader) Start() error {
	if !r.state.IntroSeen() {
		r.opr.ShowIntroPopup(r.display)
		defer r.state.MarkIntroSeen()
	}
	return r.display.Start()
}

func (r *Reader) globalKeyHandler() ui.KeyHandler {
	r.mustDefinedFields()

	statsPopupLock := make(chan struct{}, 1)

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
				r.opr.ToggleAboutPopup(r.display, r.backend)
				return nil

			case 'S':
				go func() {
					select {
					case statsPopupLock <- struct{}{}:
						defer func() { <-statsPopupLock }()
					default:
						return
					}
					r.opr.ToggleStatsPopup(r.display, r.backend)
					r.opr.Draw(r.display)
				}()
				return nil

			case 'h', '?':
				r.opr.ToggleHelpPopup(r.display)
				return nil

			case 'q':
				r.display.Stop()
				return nil
			}

		case tcell.KeyEscape:
			r.opr.UnfocusFront(r.display)
			return nil
		}

		return event
	}
}

func (r *Reader) mustDefinedFields() {
	if r.display == nil {
		panic("can not set handler with nil display")
	}

	if r.opr == nil {
		panic("can not set handler with nil operator")
	}

	if r.backend == nil {
		panic("can not set handler with nil backend")
	}
}

type Builder struct {
	ctx       context.Context
	themeName string
	scr       tcell.Screen

	// rpcBackend args.
	addr  string
	dopts []grpc.DialOption

	// For testing.
	be  bknd.Backend
	opr ui.Operator
	stt st.State
}

func NewBuilder(ctx context.Context) *Builder {
	b := Builder{
		ctx:       ctx,
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

func (b *Builder) state(stt st.State) *Builder {
	b.stt = stt
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
		opr = ui.NewDisplayOperator(b.ctx)
	}

	var stt st.State
	if b.stt != nil {
		stt = b.stt
	} else {
		stt = st.NewState()
	}

	rdr := Reader{
		ctx:     b.ctx,
		display: dsp,
		opr:     opr,
		backend: be,
		state:   stt,
	}
	rdr.display.SetHandlers(rdr.globalKeyHandler())

	return &rdr, nil
}
