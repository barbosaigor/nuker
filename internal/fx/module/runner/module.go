package runner

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/barbosaigor/nuker/internal/cli"
	masterfx "github.com/barbosaigor/nuker/internal/fx/module/domain/service/master"
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/worker"
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
	} else if strings.EqualFold(cli.LogLevel, "trace") {
		log.SetLevel(log.TraceLevel)
	} else if strings.EqualFold(cli.LogLevel, "debug") {
		log.SetLevel(log.DebugLevel)
	} else if strings.EqualFold(cli.LogLevel, "warn") {
		log.SetLevel(log.WarnLevel)
	} else if strings.EqualFold(cli.LogLevel, "error") {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var consoleWriter io.Writer = os.Stdout
	if cli.Quiet {
		consoleWriter = io.Discard
	}

	logger := log.New()
	logger.SetOutput(consoleWriter)

	return fx.Options(
		fx.Logger(newFxLooger(logger)),
		fx.Provide(context.Background),
		masterfx.Module(),
		worker.Module(cli.WorkerID, cli.MasterURI, cli.WorkerWeight),
		fx.Provide(
			runner.LoadCfg,
			runner.New,
		),
	)
}
