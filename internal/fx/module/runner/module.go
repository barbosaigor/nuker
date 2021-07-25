package runner

import (
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/pipeline"
	"github.com/barbosaigor/nuker/internal/runner"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		pipeline.Module(),
		fx.Provide(runner.New),
	)
}
