package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/bow/iris/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Logger.Error().Err(err).Msg("command failed")
		os.Exit(1)
	}
}
