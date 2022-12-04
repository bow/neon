// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/bow/iris/internal"
)

func newViper(cmdName string) *viper.Viper {
	v := viper.New()
	envSuffix := ""
	if cmdName != "" {
		envSuffix = fmt.Sprintf("_%s", strings.ReplaceAll(cmdName, "-", "_"))
	}
	envPrefix := strings.ToUpper(fmt.Sprintf("%s%s", internal.AppName(), envSuffix))
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return v
}
