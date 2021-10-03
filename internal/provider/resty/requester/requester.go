package requester

import (
	"context"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

type requester struct{}

func New() repository.Requester {
	return &requester{}
}

func (r *requester) Assign(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error {
	log.Tracef("assigning workload to worker")
	return nil
}
