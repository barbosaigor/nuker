package bufwriter

import (
	"github.com/barbosaigor/nuker/internal/provider/file/bufwriter"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			bufwriter.New,
		),
	)
}
