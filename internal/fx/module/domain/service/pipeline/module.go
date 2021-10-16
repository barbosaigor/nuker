package pipeline

import (
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	orchestratorfx "github.com/barbosaigor/nuker/internal/fx/module/domain/service/orchestrator"
	probefx "github.com/barbosaigor/nuker/internal/fx/module/domain/service/probe"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		probefx.Module(),
		orchestratorfx.Module(),
		fx.Provide(
			pipeline.New,
		),
	)
}
