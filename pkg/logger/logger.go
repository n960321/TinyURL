package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetLogger(local bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.With().Caller().Logger()
	if local {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05.999Z07:00"}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
	}
}
