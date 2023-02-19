// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type LogStyle uint8

const (
	PrettyLogStyle LogStyle = iota
	JSONLogStyle
)

func InitGlobalLog(
	logLevel string,
	style LogStyle,
	writer io.Writer,
) error {

	var (
		err error
		cw  io.Writer
		ll  = zerolog.InfoLevel
		tf  = time.RFC3339
	)

	switch style {
	case PrettyLogStyle:
		cw = zerolog.ConsoleWriter{
			Out:                 writer,
			TimeFormat:          tf,
			PartsExclude:        []string{"time"},
			FormatLevel:         prettyLevelFormatter(),
			FormatMessage:       prettyMessageFormatter(),
			FormatFieldName:     prettyFieldNameFormatter(),
			FormatFieldValue:    prettyFieldValueFormatter(),
			FormatErrFieldName:  prettyErrFieldNameFormatter(),
			FormatErrFieldValue: prettyErrFieldValueFormatter(),
		}
	case JSONLogStyle:
		cw = writer
	}

	if logLevel != "" {
		ll, err = zerolog.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			return fmt.Errorf("invalid log level '%s'", logLevel)
		}
	}

	logb := zerolog.New(cw).
		With().
		Timestamp().
		Str("app", AppName()).
		Str("version", Version())
	if inDocker() {
		logb = logb.Int("pid", os.Getpid())
	}

	lcs := logConfigState{
		logLevel:          ll,
		timestampFunc:     func() time.Time { return time.Now().UTC() },
		timeFieldFormat:   tf,
		durationFieldUnit: time.Millisecond,
		logger:            logb.Logger(),
	}
	defer func() {
		if err != nil {
			defaultLogConfigState.apply()
		}
	}()
	lcs.apply()

	return nil
}

type logConfigState struct {
	logger            zerolog.Logger
	logLevel          zerolog.Level
	timestampFunc     func() time.Time
	timeFieldFormat   string
	durationFieldUnit time.Duration
}

var defaultLogConfigState = &logConfigState{
	logger:            zlog.Logger,
	logLevel:          zerolog.GlobalLevel(),
	timestampFunc:     zerolog.TimestampFunc,
	timeFieldFormat:   zerolog.TimeFieldFormat,
	durationFieldUnit: zerolog.DurationFieldUnit,
}

func (s *logConfigState) apply() {
	zerolog.TimestampFunc = s.timestampFunc
	zerolog.TimeFieldFormat = s.timeFieldFormat
	zerolog.DurationFieldUnit = s.durationFieldUnit
	zerolog.SetGlobalLevel(s.logLevel)
	zlog.Logger = s.logger
}

// Adapted from:
// https://github.com/rs/zerolog/blob/762546b5c64e03f3d23f867213e80aa45906aaf7/console.go
func prettyLevelFormatter() zerolog.Formatter {

	return func(i any) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case zerolog.LevelTraceValue:
				l = colorizeLevel("TRC", colorWhite)
			case zerolog.LevelDebugValue:
				l = colorizeLevel("DBG", colorMagenta)
			case zerolog.LevelInfoValue:
				l = colorizeLevel("INF", colorCyan)
			case zerolog.LevelWarnValue:
				l = colorizeLevel("WRN", colorYellow)
			case zerolog.LevelErrorValue:
				l = colorizeLevel("ERR", colorRed)
			case zerolog.LevelFatalValue:
				l = colorizeLevel("FTL", colorRed)
			case zerolog.LevelPanicValue:
				l = colorizeLevel("PNC", colorRed)
			default:
				l = colorize(ll, colorBold)
			}
		} else {
			if i == nil {
				l = colorize("???", colorBold)
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}
		return l
	}
}

func prettyMessageFormatter() zerolog.Formatter {
	return func(i any) string {
		if i != nil {
			msg := fmt.Sprintf("%s", i)
			// Use 'range' to get to the first character, which may span
			// more than one byte.
			for idx, c := range msg {
				return colorize(
					string(unicode.ToUpper(c))+msg[idx+1:],
					colorBold,
				)
			}
		}
		return ""
	}
}

func prettyFieldNameFormatter() zerolog.Formatter {
	return func(i any) string {
		return fmt.Sprintf(
			"%s=",
			colorize(fmt.Sprintf("%s", i), colorDarkGray),
		)
	}
}

func prettyFieldValueFormatter() zerolog.Formatter {
	return func(i any) string {
		return colorize(fmt.Sprintf("%s", i), colorYellow)
	}
}

func prettyErrFieldNameFormatter() zerolog.Formatter {
	return func(i any) string {
		return fmt.Sprintf(
			"%s=",
			colorize(fmt.Sprintf("%s", i), colorRed),
		)
	}
}

func prettyErrFieldValueFormatter() zerolog.Formatter {
	return func(i any) string {
		return fmt.Sprintf("%s", i)
	}
}

// Adapted from:
// https://github.com/rs/zerolog/blob/762546b5c64e03f3d23f867213e80aa45906aaf7/console.go
func colorize(s any, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func invert(s any) string {
	return fmt.Sprintf("\x1b[7m%s\x1b[27m", s)
}

func colorizeLevel(s any, c int) string {
	return invert(colorize(colorize(fmt.Sprintf(" %s ", s), c), colorBold))
}

// Copied from:
// https://github.com/rs/zerolog/blob/762546b5c64e03f3d23f867213e80aa45906aaf7/console.go
const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

func inDocker() bool {
	_, errStat := os.Stat("/.dockerenv")
	return errStat != nil
}
