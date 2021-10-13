package orchestrator

import (
	"sync"

	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/barbosaigor/nuker/internal/domain/service/requester"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	publisherfx "github.com/barbosaigor/nuker/internal/fx/module/domain/service/publisher"
	fxworker "github.com/barbosaigor/nuker/internal/fx/module/domain/service/worker"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type orchestratorParams struct {
	fx.In

	Worker           worker.Worker `optional:"true"`
	RequesterFactory requester.Factory
}

var once = &sync.Once{}

func Module() fx.Option {
	opts := fx.Options()

	once.Do(func() {
		opt := fx.Options()
		if !cli.Master || cli.Worker {
			opt = fxworker.Module(uuid.New().String(), 1)
		}

		opts = fx.Options(
			opt,
			publisherfx.Module(),
			fx.Provide(
				requester.NewFactory,
				func(params orchestratorParams) orchestrator.Orchestrator {
					workers := map[string]worker.Worker{}

					if params.Worker != nil {
						workers[params.Worker.ID()] = params.Worker
					}

					return orchestrator.New(workers, params.RequesterFactory)
				},
			),
		)
	})

	return opts
}
