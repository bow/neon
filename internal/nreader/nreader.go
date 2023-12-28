// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package nreader

import (
	"context"
	"fmt"

	"github.com/bow/neon/internal/nreader/view"
)

//nolint:unused
type Reader struct {
	ctx      context.Context
	initPath string

	view  view.Viewer
	model *model
}

func (r *Reader) Start() error {
	return fmt.Errorf("Start is unimplemented")
}
