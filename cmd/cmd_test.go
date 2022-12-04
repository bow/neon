// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"bytes"

	"github.com/spf13/cobra"
)

// execCommand executes the command for testing.
func execCommand(args []string) (err error, stdout string, stderr string) {
	cmd, outb, errb := newCommand()
	cmd.SetArgs(args)

	err = cmd.Execute()

	return err, outb.String(), errb.String()
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
