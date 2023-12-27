// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/bow/neon/cmd"
	"github.com/bow/neon/internal"
)

func main() {
	ctx := context.Background()
	command := cmd.New()
	internal.MustSetupLogging(command.ErrOrStderr())

	if err := command.ExecuteContext(ctx); err != nil {
		log.Logger.Error().Msgf("%s", err)
		os.Exit(1)
	}
}
