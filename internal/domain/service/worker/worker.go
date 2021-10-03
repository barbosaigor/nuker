package worker

import (
	"context"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

type Worker interface {
	ID() string
	Weight() int
	Do(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics)
}

type worker struct {
	id        string
	weight    int
	requester repository.Requester
}

func New(ID string, weight int, requester repository.Requester) Worker {
	return worker{
		id:        ID,
		weight:    weight,
		requester: requester,
	}
}

func (w worker) Do(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) {
	log.
		WithField("worker", w.id).
		Tracef("request count: %d", wl.RequestsCount)

	_ = w.requester.Assign(ctx, wl, metChan)
}

func (w worker) ID() string {
	return w.id
}

func (w worker) Weight() int {
	return w.weight
}
