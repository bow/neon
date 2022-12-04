// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/bow/iris/cmd"
)

func main() {
	ctx := context.Background()
	command := cmd.New()

	if err := command.ExecuteContext(ctx); err != nil {
		log.Logger.Error().Err(err).Msg("command failed")
		os.Exit(1)
	}
}
