package master

import (
	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/master/http"
	"github.com/barbosaigor/nuker/internal/domain/service/master/local"
	pipelinefx "github.com/barbosaigor/nuker/internal/fx/module/domain/service/pipeline"
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/publisher"
	"go.uber.org/fx"
)

func Module() fx.Option {
	var opt fx.Option

	if cli.Master {
		opt = fx.Options(
			fx.Provide(
				func() http.Options {
					return http.Options{
						Port:       cli.Port,
						MinWorkers: cli.MinWorkers,
					}
				},
				http.NewMaster,
			),
		)
	} else {
		opt = fx.Options(
			fx.Provide(
				local.NewMaster,
			),
		)
	}

	return fx.Options(
		publisher.Module(),
		pipelinefx.Module(),
		opt,
	)
}
