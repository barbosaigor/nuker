package orchestrator

import (
	"context"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/pkg/metrics"
)

type WorkloadAssigner interface {
	AssignWorkload(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error
}
