package probe

import (
	"github.com/barbosaigor/nuker/internal/domain/service/probe"
	"github.com/barbosaigor/nuker/internal/fx/module/provider/file/bufwriter"
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
