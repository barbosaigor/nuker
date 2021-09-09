package orchestrator

import (
	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	fxworker "github.com/barbosaigor/nuker/internal/fx/module/domain/service/worker"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fxworker.Module("w-master", 1),
		fx.Provide(
			func(w worker.Worker) orchestrator.Orchestrator {
				return orchestrator.New(map[string]worker.Worker{w.ID(): w})
			},
		),
	)
}
