package util

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

func DefaultRootLog() zerolog.Logger {
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		panic(err)
	}
	if logLevel == zerolog.NoLevel {
		logLevel = zerolog.InfoLevel
	}

	logJSON := os.Getenv("LOG_JSON") == "true"

	var writer io.Writer = os.Stdout
	if !logJSON {
		writer = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.RFC3339
		})
	}

	log := zerolog.New(writer).With().Timestamp().Logger().Level(logLevel)
	return log
}
