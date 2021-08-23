package probe

import (
	"github.com/barbosaigor/nuker/internal/domain/service/probe"
	"github.com/barbosaigor/nuker/internal/fx/module/domain/service/bufwriter"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		bufwriter.Module(),
		fx.Provide(
			probe.New,
		),
	)
}
