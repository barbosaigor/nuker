package pipeline

import (
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	fxprobe "github.com/barbosaigor/nuker/internal/fx/module/domain/service/probe"
	fxpublisher "github.com/barbosaigor/nuker/internal/fx/module/domain/service/publisher"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fxpublisher.Module(),
		fxprobe.Module(),
		fx.Provide(
			pipeline.New,
		),
	)
}
