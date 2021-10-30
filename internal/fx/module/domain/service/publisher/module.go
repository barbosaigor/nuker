package publisher

import (
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	fxstreamer "github.com/barbosaigor/nuker/internal/fx/module/provider/resty/streamer"
	"go.uber.org/fx"
)

var once = &sync.Once{}

func Module() fx.Option {
	opts := fx.Options()

	once.Do(func() {
		opts = fx.Options(
			fxstreamer.Module(),
			fx.Provide(
				publisher.New,
			),
		)
	})

	return opts
}
