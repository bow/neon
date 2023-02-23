// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	"github.com/bow/iris/internal"
)

func newViper(cmdName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(internal.EnvKey(cmdName))
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return v
}

type ctxKey string

func inCmdContext(cmd *cobra.Command, key string, value any) {
	ctx := context.WithValue(cmd.Context(), ctxKey(key), value)
	cmd.SetContext(ctx)
}

func fromCmdContext[T any](cmd *cobra.Command, key string) (T, error) {
	var zero T
	val, ok := cmd.Context().Value(ctxKey(key)).(T)
	if !ok {
		return zero, fmt.Errorf("error retrieving %q from command context", key)
	}
	return val, nil
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
