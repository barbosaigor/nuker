package probe

import (
	"github.com/barbosaigor/nuker/internal/domain/service/probe"
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/bufwriter"
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/view"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		bufwriter.Module(),
		view.Module(),
		fx.Provide(
			probe.New,
		),
	)
}
