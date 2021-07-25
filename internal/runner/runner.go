package runner

import (
	"context"
	"os"

	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	"github.com/barbosaigor/nuker/pkg/config"
	pkgrunner "github.com/barbosaigor/nuker/pkg/runner"
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
	data, err := os.ReadFile(cfgFileName)
	if err != nil {
		return err
	}

	cfg, err := config.YamlUnmarshal(data)
	if err != nil {
		return err
	}

	return r.pipeline.Run(ctx, *cfg)
}
