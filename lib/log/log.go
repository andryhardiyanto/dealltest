package log

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func RegisterLogger(debug bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log = zerolog.
		New(os.Stderr).
		With().
		Timestamp().
		Logger()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Panic() *zerolog.Event {
	return log.Panic()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Hook(hook zerolog.Hook) zerolog.Logger {
	return log.Hook(hook)
}

func GetLogger() zerolog.Logger {
	return log
}
