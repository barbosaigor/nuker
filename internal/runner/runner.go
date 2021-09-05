package runner

import (
	"context"
	"errors"
	"os"

	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	"github.com/barbosaigor/nuker/pkg/config"
	pkgrunner "github.com/barbosaigor/nuker/pkg/runner"
	log "github.com/sirupsen/logrus"
)

type runner struct {
	pipeline pipeline.Pipeline
	opts     Options
}

func New(pipeline pipeline.Pipeline, opts Options) pkgrunner.Runner {
	return &runner{
		pipeline: pipeline,
		opts:     opts,
	}
}

func (r *runner) Run(ctx context.Context) error {
	switch r.opts.Op {
	case exec:
		return r.exec(ctx)
	case run:
		return r.run(ctx)
	default:
		return nil
	}
}

func (r *runner) exec(ctx context.Context) error {
	cfg := cli.BuildExecCmdCfg()
	log.Tracef("config: %+v", cfg)
	if cfg == nil {
		return errors.New("nil exec config")
	}

	if cli.DryRunFlagExecCmd {
		log.Infof("plan: %+v", *cfg)
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
