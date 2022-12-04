// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoArgs(t *testing.T) {
	stdout, stderr, err := execCommand(nil)
	require.NoError(t, err)

	assert.Empty(t, stderr)
	assert.Contains(t, stdout, "Feed reader suite")
	assert.Contains(t, stdout, `Use "iris [command] --help" for more information`)
}
