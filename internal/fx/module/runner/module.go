package runner

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/pipeline"
	"github.com/barbosaigor/nuker/internal/runner"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func Module() fx.Option {
	cli.ExecCli()

	log.SetFormatter(&log.TextFormatter{
		ForceQuote:    true,
		FullTimestamp: true,
	})

	if cli.Verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var consoleWriter io.Writer = os.Stdout
	if cli.Quiet {
		consoleWriter = ioutil.Discard
	}

	logger := log.New()
	logger.SetOutput(consoleWriter)

	return fx.Options(
		fx.Logger(newFxLooger(logger)),
		fx.Provide(context.Background),
		pipeline.Module(),
		fx.Provide(
			runner.LoadCfg,
			runner.New,
		),
	)
}
