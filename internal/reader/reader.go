// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	rp "github.com/bow/neon/internal/reader/repo"
	"github.com/bow/neon/internal/reader/ui"
)

//nolint:unused
type Reader struct {
	ctx      context.Context
	initPath string

	view ui.Viewer
	repo rp.Repo
}

func (r *Reader) Show() error {
	return r.view.Show()
}

type Builder struct {
	ctx       context.Context
	initPath  string
	themeName string

	// rpcRepo args.
	addr  string
	dopts []grpc.DialOption

	// For testing.
	rpo rp.Repo
	vwr ui.Viewer
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

func (b *Builder) repo(rpo rp.Repo) *Builder {
	b.rpo = rpo
	return b
}

func (b *Builder) viewer(v ui.Viewer) *Builder {
	b.vwr = v
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

	var viewer ui.Viewer
	if b.vwr != nil {
		viewer = b.vwr
	} else {
		viewer, err = ui.NewView(b.themeName)
		if err != nil {
			return nil, err
		}
	}

	rdr := Reader{
		ctx:  b.ctx,
		view: viewer,
		repo: rpo,
	}
	rdr.setKeyHandlers()

	return &rdr, nil
}

func (r *Reader) setKeyHandlers() {
}
