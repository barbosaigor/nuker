package runner

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type fxLogger struct {
	logger *zerolog.Logger
}

func newFxLooger(ctx context.Context) fx.Printer {
	return &fxLogger{
		logger: log.Ctx(ctx),
	}
}

func (l fxLogger) Printf(format string, v ...interface{}) {
	l.logger.
		Trace().
		Msgf(format, v...)
}
