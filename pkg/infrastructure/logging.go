// Package infrastructure is implements an adapter to talks low-level modules.
// # This manifest was generated by ymir. DO NOT EDIT.
package infrastructure

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/shared/tracer"
)

// InitializeLogger will set logging format.
func InitializeLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	writer := &tracer.ZeroWriter{MinLevel: zerolog.DebugLevel}
	log.Logger = zerolog.New(
		zerolog.MultiLevelWriter(os.Stdout, writer)).
		With().Timestamp().Caller().Logger()
	if strings.EqualFold(Envs.App.Environment, Development) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
