package master

import (
	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/master/http"
	pipelinefx "github.com/barbosaigor/nuker/internal/fx/module/domain/service/pipeline"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		pipelinefx.Module(),
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
}
