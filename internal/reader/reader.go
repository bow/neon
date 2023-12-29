// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package reader

import (
	"context"
	"fmt"

	"github.com/bow/neon/internal/reader/model"
	"github.com/bow/neon/internal/reader/view"
)

//nolint:unused
type Reader struct {
	ctx      context.Context
	initPath string

	view  view.Viewer
	model model.Model
}

func (r *Reader) Start() error {
	return fmt.Errorf("Start is unimplemented")
}
