package tview

import (
	"github.com/barbosaigor/nuker/internal/domain/service/view/nop"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			nop.New,
		),
	)
}
