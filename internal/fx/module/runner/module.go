package runner

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/pipeline"
	"github.com/barbosaigor/nuker/internal/runner"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func Module() fx.Option {
	cli.ExecCli()

	if cli.Verbose {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	var consoleWriter io.Writer = os.Stdout
	if cli.Quiet {
		consoleWriter = ioutil.Discard
	}

	logger := zerolog.
		New(consoleWriter).
		With().
		Timestamp().
		Logger()

	ctx := logger.WithContext(context.Background())

	return fx.Options(
		fx.Logger(newFxLooger(ctx)),
		fx.Provide(
			func() context.Context {
				return ctx
			},
		),
		pipeline.Module(),
		fx.Provide(runner.New),
	)
}
