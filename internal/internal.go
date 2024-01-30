// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"fmt"
	"strings"
)

// AppName returns the application name.
func AppName() string {
	return "neon"
}

// AppHomepage returns the application homepage.
func AppHomepage() string {
	return "https://github.com/bow/neon"
}

// EnvKey returns the environment variable key for configuration.
func EnvKey(key string) string {
	if key == "" {
		return envPrefix()
	}
	return fmt.Sprintf("%s_%s", envPrefix(), strings.ToUpper(strings.ReplaceAll(key, "-", "_")))
}

// Banner shows the application name as ASCII art.
func Banner() string {
	return `    _   __
   / | / /___   ____   ____
  /  |/ // _ \ / __ \ / __ \
 / /|  //  __// /_/ // / / /
/_/ |_/ \___/ \____//_/ /_/`
}

// envPrefix returns the environment variable prefix for configuration.
func envPrefix() string {
	return strings.ToUpper(AppName())
}
