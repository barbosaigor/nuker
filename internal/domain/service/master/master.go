package master

import (
	"context"

	"github.com/barbosaigor/nuker/pkg/config"
)

type Master interface {
	Run(ctx context.Context, cfg config.Config) error
}
