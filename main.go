// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/bow/iris/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		log.Logger.Error().Err(err).Msg("command failed")
		os.Exit(1)
	}
}
