package repository

import (
	"context"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/pkg/metrics"
)

type Requester interface {
	Assign(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error
}
