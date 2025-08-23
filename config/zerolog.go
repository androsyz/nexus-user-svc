package config

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return logger
}
