package worker

import (
	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	"github.com/barbosaigor/nuker/internal/domain/service/worker/courier"
	fxpublisher "github.com/barbosaigor/nuker/internal/fx/module/domain/service/publisher"
	"go.uber.org/fx"
)

func Module(ID string, weight int) fx.Option {
	return fx.Options(
		fxpublisher.Module(),
		fx.Provide(
			func(pub publisher.Publisher) worker.Worker {
				return worker.New(ID, weight, courier.New(pub))
			},
		),
	)
}
