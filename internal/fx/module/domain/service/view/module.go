package view

import (
	"github.com/barbosaigor/nuker/internal/domain/service/view"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			view.New,
		),
	)
}
