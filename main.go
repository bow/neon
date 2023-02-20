// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bow/iris/cmd"
)

func main() {
	ctx := context.Background()
	command := cmd.New()

	if err := command.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(command.ErrOrStderr(), "%s\n", err)
		os.Exit(1)
	}
}
