// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"google.golang.org/grpc"

	bknd "github.com/bow/neon/internal/reader/backend"
	st "github.com/bow/neon/internal/reader/state"
	"github.com/bow/neon/internal/reader/ui"
)

type Reader struct {
	ctx context.Context

	display *ui.Display
	opr     ui.Operator
	backend bknd.Backend
	state   st.State

	callTimeout time.Duration
}

func (r *Reader) Start() error {
	if !r.state.IntroSeen() {
		r.opr.ShowIntroPopup(r.display)
		defer r.state.MarkIntroSeen()
	}
	r.opr.ShowAllFeeds(r.display, r.backend.GetAllFeedsF(r.ctx))
	r.opr.RefreshStats(r.display, r.backend.GetStatsF(r.ctx))
	return r.display.Start()
}

// nolint:revive
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
				r.opr.ToggleAboutPopup(r.display, r.backend.String())
				return nil

			case 'E':
				r.opr.FocusEntriesPane(r.display)
				return nil

			case 'F':
				r.opr.FocusFeedsPane(r.display)
				return nil

			case 'R':
				r.opr.FocusReadingPane(r.display)
				return nil

			case 'S':
				go func() {
					select {
					case statsPopupLock <- struct{}{}:
						defer func() { <-statsPopupLock }()
					default:
						return
					}
					ctx, cancel := r.callCtx()
					defer cancel()
					r.opr.ToggleStatsPopup(r.display, r.backend.GetStatsF(ctx))
					r.display.Draw()
				}()
				return nil

			case 'H', '?':
				r.opr.ToggleHelpPopup(r.display)
				return nil

			case 'b':
				r.opr.ToggleStatusBar(r.display)
				return nil

			case 'c':
				r.opr.ClearStatusBar(r.display)
				return nil

			case 'q':
				r.display.Stop()
				return nil
			}

		case tcell.KeyTab:
			if event.Modifiers()&tcell.ModAlt == 0 {
				r.opr.FocusNextPane(r.display)
			} else {
				r.opr.FocusPreviousPane(r.display)
			}
			return nil

		case tcell.KeyEscape:
			r.opr.UnfocusFront(r.display)
			return nil
		}

		return event
	}
}

func (r *Reader) feedsPaneKeyHandler() ui.KeyHandler {
	pullFeedsLock := make(chan struct{}, 1)

	return func(event *tcell.EventKey) *tcell.EventKey {
		keyr := event.Rune()

		if keyr == 'P' {
			go func() {
				select {
				case pullFeedsLock <- struct{}{}:
					defer func() { <-pullFeedsLock }()
				default:
					return
				}
				ctxf, cancelf := r.callCtx()
				defer cancelf()
				r.opr.RefreshFeeds(r.display, r.backend.PullFeedsF(ctxf, nil, false))

				ctxs, cancels := r.callCtx()
				defer cancels()
				r.opr.RefreshStats(r.display, r.backend.GetStatsF(ctxs))

				r.display.Draw()
			}()
			return nil
		}
		return event
	}
}

func (r *Reader) callCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.ctx, r.callTimeout)
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
	addr        string
	dopts       []grpc.DialOption
	callTimeout time.Duration

	// For testing.
	be  bknd.Backend
	opr ui.Operator
	stt st.State
}

func NewBuilder(ctx context.Context) *Builder {
	b := Builder{
		ctx:         ctx,
		themeName:   "dark",
		dopts:       nil,
		callTimeout: 3 * time.Second,
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

func (b *Builder) CallTimeout(timeout time.Duration) *Builder {
	b.callTimeout = timeout
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
		opr = ui.NewDisplayOperator()
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

		callTimeout: b.callTimeout,
	}
	rdr.display.SetHandlers(
		rdr.globalKeyHandler(),
		rdr.feedsPaneKeyHandler(),
	)

	return &rdr, nil
}
