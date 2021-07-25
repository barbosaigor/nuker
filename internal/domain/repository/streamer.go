package repository

import (
	"context"

	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/barbosaigor/nuker/pkg/metrics"
)

type Streamer interface {
	Stream(ctx context.Context, cfg config.Network) (*metrics.NetworkMetrics, error)
}
