// Copyright (c) 2023 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package api

import (
	// Embed needs to be imported so we can use the go:embed directive.
	_ "embed"
)

//go:embed iris.proto
var proto []byte

// Proto() returns the proto file that describes the server interface.
func Proto() []byte {
	return proto
}
