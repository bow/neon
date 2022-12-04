// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"testing"

	"github.com/bow/iris/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionOk(t *testing.T) {
	stdout, stderr, err := execCommand([]string{"version"})
	require.NoError(t, err)

	assert.Empty(t, stderr)
	assert.Contains(t, stdout, "App        : iris")
	assert.Contains(t, stdout, fmt.Sprintf("Version    : %s", internal.Version()))
}
