package view

import (
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/view/tview"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		tview.Module(),
	)
}
