// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bow/iris/internal"
)

func TestNoArgs(t *testing.T) {
	stdout, stderr, err := execCommand(nil)
	require.NoError(t, err)

	assert.Empty(t, stderr)
	assert.Contains(t, stdout, "Feed reader suite")
	assert.Contains(t, stdout, `Use "iris [command] --help" for more information`)
}

func TestVersionOk(t *testing.T) {
	stdout, stderr, err := execCommand([]string{"version"})
	require.NoError(t, err)

	assert.Empty(t, stderr)
	assert.Contains(t, stdout, "App        : iris")
	assert.Contains(t, stdout, fmt.Sprintf("Version    : %s", internal.Version()))
}

// execCommand executes the command for testing.
func execCommand(args []string) (stdout string, stderr string, err error) {
	cmd, outb, errb := newCommand()
	cmd.SetArgs(args)

	err = cmd.Execute()

	return outb.String(), errb.String(), err
}

// newCommand creates new cobra command for testing. It also returns the buffers that
// capture the command's stdout and stderr, respectively.
func newCommand() (cmd *cobra.Command, outb *bytes.Buffer, errb *bytes.Buffer) {
	cmd = New()

	outb = bytes.NewBufferString("")
	cmd.SetOut(outb)

	errb = bytes.NewBufferString("")
	cmd.SetErr(errb)

	return cmd, outb, errb
}
