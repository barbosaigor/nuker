package requester

import (
	"github.com/barbosaigor/nuker/internal/domain/service/requester"
	publisherfx "github.com/barbosaigor/nuker/internal/fx/module/domain/service/publisher"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		publisherfx.Module(),
		fx.Provide(
			requester.New,
		),
	)
}
