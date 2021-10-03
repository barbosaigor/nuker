package orchestrator

import (
	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	fxworker "github.com/barbosaigor/nuker/internal/fx/module/domain/service/worker"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type orchestratorParams struct {
	fx.In

	Worker  worker.Worker `optional:"true"`
	Options orchestrator.Options
}

func Module() fx.Option {
	opt := fx.Options()

	if !cli.Master || cli.Worker {
		opt = fxworker.Module(uuid.New().String(), 1)
	}

	return fx.Options(
		opt,
		fx.Provide(
			func() orchestrator.Options {
				return orchestrator.Options{
					Port: cli.Port,
				}
			},
			func(params orchestratorParams) orchestrator.Orchestrator {
				workers := map[string]worker.Worker{}

				if params.Worker != nil {
					workers[params.Worker.ID()] = params.Worker
				}

				return orchestrator.New(workers, params.Options)
			},
		),
	)
}
