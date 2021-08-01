package bufwriter

import (
	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/barbosaigor/nuker/internal/provider/file/bufwriter"
	nopBufWriter "github.com/barbosaigor/nuker/internal/provider/nop/bufwriter"
	"go.uber.org/fx"
)

func Module() fx.Option {
	if cli.NoLogFile {
		return fx.Options(
			fx.Provide(
				nopBufWriter.New,
			),
		)
	}

	return fx.Options(
		fx.Provide(
			func() (repository.BufWriter, error) {
				return bufwriter.New(cli.LogFile)
			},
		),
	)
}
