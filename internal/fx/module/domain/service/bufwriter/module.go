package bufwriter

import (
	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/bufwriter"
	"github.com/barbosaigor/nuker/internal/domain/service/bufwriter/filebufwriter"
	nopbufwriter "github.com/barbosaigor/nuker/internal/domain/service/bufwriter/nopbufwriter"
	"go.uber.org/fx"
)

func Module() fx.Option {
	if cli.NoLogFile {
		return fx.Options(
			fx.Provide(
				nopbufwriter.New,
			),
		)
	}

	return fx.Options(
		fx.Provide(
			func() (bufwriter.BufWriter, error) {
				return filebufwriter.New(cli.LogFile)
			},
		),
	)
}
