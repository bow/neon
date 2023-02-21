// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/viper"

	"github.com/bow/iris/internal"
)

func newViper(cmdName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(internal.EnvKey(cmdName))
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return v
}

// showBanner prints the application banner to the given writer.
func showBanner(w io.Writer) {
	fmt.Fprintf(
		w,
		`    ____       _
   /  _/_____ (_)_____
   / / / ___// // ___/
 _/ / / /   / /(__  )
/___//_/   /_//____/

`)
}
