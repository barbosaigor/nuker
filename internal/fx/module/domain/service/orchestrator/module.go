package orchestrator

import (
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"go.uber.org/fx"
)

var once = &sync.Once{}

func Module() fx.Option {
	opts := fx.Options()
	once.Do(func() {
		opts = fx.Options(
			fx.Provide(
				orchestrator.New,
			),
		)
	})
	return opts
}
