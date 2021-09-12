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

	Worker worker.Worker `optional:"true"`
}

func Module() fx.Option {
	opt := fx.Options()

	if !cli.Master || cli.Worker {
		opt = fxworker.Module(uuid.New().String(), 1)
	}

	return fx.Options(
		opt,
		fx.Provide(
			func(params orchestratorParams) orchestrator.Orchestrator {
				var workers map[string]worker.Worker

				if params.Worker != nil {
					workers = map[string]worker.Worker{params.Worker.ID(): params.Worker}
				}

				return orchestrator.New(workers)
			},
		),
	)
}
