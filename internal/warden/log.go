package warden

import (
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

var Log *Logger

func NewLog(writer io.Writer, level zerolog.Level, timeFormat string) *Logger {
	zl := zerolog.New(zerolog.SyncWriter(writer)).Output(zerolog.ConsoleWriter{Out: writer, TimeFormat: timeFormat, FormatLevel: func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}}).With().Timestamp().Logger().Level(level)
	return &Logger{zl}
}

func SetLog(logger *Logger) {
	Log = logger
}
