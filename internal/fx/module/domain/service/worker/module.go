package worker

import (
	"github.com/barbosaigor/nuker/internal/domain/service/requester"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	requesterfx "github.com/barbosaigor/nuker/internal/fx/module/domain/service/requester"
	"go.uber.org/fx"
)

func Module(ID, masterURI string, weight int) fx.Option {
	return fx.Options(
		requesterfx.Module(),
		fx.Provide(
			func(req requester.Requester) worker.Worker {
				return worker.New(ID, masterURI, weight, req)
			},
		),
	)
}
