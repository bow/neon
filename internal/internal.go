// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"fmt"
	"strings"
)

// AppName returns the application name.
func AppName() string {
	return "iris"
}

// EnvKey returns the environment variable key for configuration.
func EnvKey(key string) string {
	if key == "" {
		return envPrefix()
	}
	return fmt.Sprintf("%s_%s", envPrefix(), strings.ToUpper(strings.ReplaceAll(key, "-", "_")))
}

// envPrefix returns the environment variable prefix for configuration.
func envPrefix() string {
	return strings.ToUpper(AppName())
}
