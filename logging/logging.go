package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/bow/courier/version"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Style uint8

const (
	PrettyConsoleStyle Style = iota
	JSONStyle
)

func Init(
	logLevel string,
	style Style,
	writer io.Writer,
) error {

	var (
		err error
		cw  io.Writer
		ll  = zerolog.InfoLevel
		tf  = time.RFC3339
	)

	switch style {
	case PrettyConsoleStyle:
		cw = zerolog.ConsoleWriter{Out: writer, TimeFormat: tf}
	case JSONStyle:
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
			Str("app", version.AppName()).
			Str("app_version", version.Version()).
			Str("grpc.version", grpc.Version).
			Int("pid", os.Getpid()).
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
