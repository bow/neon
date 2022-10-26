// Copyright (c) 2022 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type LogStyle uint8

const (
	PrettyLogStyle LogStyle = iota
	JSONLogStyle
)

func SetLogPID() {
	zlog.Logger = zlog.Logger.With().Int("pid", os.Getpid()).Logger()
}

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
		cw = zerolog.ConsoleWriter{Out: writer, TimeFormat: tf}
	case JSONLogStyle:
		cw = writer
	}

	if logLevel != "" {
		ll, err = zerolog.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			return fmt.Errorf("invalid log level '%s'", logLevel)
		}
	}

	lcs := logConfigState{
		logLevel:          ll,
		timestampFunc:     func() time.Time { return time.Now().UTC() },
		timeFieldFormat:   tf,
		durationFieldUnit: time.Millisecond,
		logger: zerolog.New(cw).
			With().
			Timestamp().
			Str("app", AppName()).
			Str("version", Version()).
			Logger(),
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
