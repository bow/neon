// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	m "github.com/bow/neon/internal/reader/model"
	v "github.com/bow/neon/internal/reader/view"
)

//nolint:unused
type Reader struct {
	ctx      context.Context
	initPath string

	view  v.Viewer
	model m.Model
}

func (r *Reader) Show() error {
	return r.view.Show()
}

type Builder struct {
	ctx       context.Context
	initPath  string
	themeName string

	// rpcModel args.
	addr  string
	dopts []grpc.DialOption

	// For testing.
	mod m.Model
	vwr v.Viewer
}

func NewBuilder() *Builder {
	b := Builder{
		ctx:       context.Background(),
		dopts:     nil,
		themeName: "dark",
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

func (b *Builder) model(mod m.Model) *Builder {
	b.mod = mod
	return b
}

func (b *Builder) viewer(vwr v.Viewer) *Builder {
	b.vwr = vwr
	return b
}

func (b *Builder) Build() (*Reader, error) {

	if b.addr == "" && b.mod == nil {
		return nil, fmt.Errorf("reader server address must be specified")
	}

	var (
		mod m.Model
		err error
	)
	if b.mod != nil {
		mod = b.mod
	} else {
		mod, err = m.NewRPCModel(b.ctx, b.addr, b.dopts...)
		if err != nil {
			return nil, err
		}
	}

	var viewer v.Viewer
	if b.vwr != nil {
		viewer = b.vwr
	} else {
		viewer, err = v.New(b.themeName)
		if err != nil {
			return nil, err
		}
	}

	rdr := Reader{
		ctx:   b.ctx,
		view:  viewer,
		model: mod,
	}
	rdr.setKeyHandlers()

	return &rdr, nil
}

func (r *Reader) setKeyHandlers() {
}
