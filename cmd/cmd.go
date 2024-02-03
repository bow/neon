// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	"github.com/bow/neon/internal"
)

const (
	dbPathKey         = "db-path"
	defaultServerAddr = "127.0.0.1:5151"
)

func newViper(cmdName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(internal.EnvKey(cmdName))
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return v
}

type ctxKey string

func toCmdContext(cmd *cobra.Command, key string, value any) {
	ctx := context.WithValue(cmd.Context(), ctxKey(key), value)
	cmd.SetContext(ctx)
}

func fromCmdContext[T any](cmd *cobra.Command, key string) (T, error) {
	var zero T
	val, ok := cmd.Context().Value(ctxKey(key)).(T)
	if !ok {
		return zero, fmt.Errorf("can not retrieve %T value %[1]q from command context", key)
	}
	return val, nil
}

// showBanner prints the application banner to the given writer.
func showBanner(w io.Writer) {
	fmt.Fprintf(w, "%s\n\n", internal.Banner())
}

func makeAlias(name string) []string {
	return []string{name[:1]}
}

func derefOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func fmtOrEmpty(p *time.Time) string {
	if p == nil {
		return ""
	}
	return fmtTime(*p)
}

func fmtTime(value time.Time) string {
	return value.Local().Format("2 January 2006 • 15:04 MST")
}
