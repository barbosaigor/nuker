package runner

import (
	"context"
	"errors"
	"os"

	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	"github.com/barbosaigor/nuker/pkg/config"
	pkgrunner "github.com/barbosaigor/nuker/pkg/runner"
	"github.com/rs/zerolog/log"
)

type runner struct {
	pipeline pipeline.Pipeline
}

func New(pipeline pipeline.Pipeline) pkgrunner.Runner {
	return &runner{
		pipeline: pipeline,
	}
}

func (r *runner) Run(ctx context.Context) error {
	if cli.IsExec {
		return r.exec(ctx)
	}

	if cli.IsRun {
		return r.run(ctx)
	}

	return nil
}

func (r *runner) exec(ctx context.Context) error {
	cfg := cli.BuildExecCmdCfg()
	log.Trace().Msgf("config: %+v", cfg)
	if cfg == nil {
		return errors.New("nil exec config")
	}

	if cli.DryRunFlagExecCmd {
		log.Info().Msgf("plan: %+v", *cfg)
		return nil
	}

	return r.pipeline.Run(ctx, *cfg)
}

func (r *runner) run(ctx context.Context) error {
	if len(cli.Args) == 0 {
		return nil
	}

	data, err := os.ReadFile(cli.Args[0])
	if err != nil {
		return err
	}

	cfg, err := config.YamlUnmarshal(data)
	if err != nil {
		return err
	}

	return r.pipeline.Run(ctx, *cfg)
}
