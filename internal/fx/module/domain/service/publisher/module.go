package publisher

import (
	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	fxstreamer "github.com/barbosaigor/nuker/internal/fx/module/provider/resty/streamer"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fxstreamer.Module(),
		fx.Provide(
			publisher.New,
		),
	)
}
