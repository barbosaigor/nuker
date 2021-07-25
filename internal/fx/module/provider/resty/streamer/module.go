package streamer

import (
	"github.com/barbosaigor/nuker/internal/provider/resty/streamer"
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			resty.New,
			streamer.New,
		),
	)
}
