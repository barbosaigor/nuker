package pipeline

import (
	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	fxorchestrator "github.com/barbosaigor/nuker/internal/fx/module/domain/service/orchestrator"
	fxprobe "github.com/barbosaigor/nuker/internal/fx/module/domain/service/probe"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fxprobe.Module(),
		fxorchestrator.Module(),
		fx.Provide(
			func() pipeline.Options {
				return pipeline.Options{
					IsMaster: cli.Master,
				}
			},
			pipeline.New,
		),
	)
}
