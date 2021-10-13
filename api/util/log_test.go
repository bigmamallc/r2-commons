package util

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	log1 := DefaultRootLog()
	log1.Info().Msg("hi")

	os.Setenv("LOG_JSON", "true")
	log2 := DefaultRootLog()
	log2.Info().Msg("hi")
}
