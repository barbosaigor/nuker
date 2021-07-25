package fxrunner

import (
	"context"
	"os"

	"github.com/barbosaigor/nuker/pkg/runner"
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

func Run(opt fx.Option) error {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger()

	ctx := logger.WithContext(context.Background())

	app := fx.New(
		fx.Logger(newFxLooger(ctx)),
		fx.Provide(
			func() context.Context {
				return ctx
			},
		),
		opt,
		fx.Invoke(runner.StartRunner),
	)

	if err := app.Err(); err != nil {
		return err
	}

	return nil
}
