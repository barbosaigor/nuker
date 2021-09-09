package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

type Worker interface {
	Do(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics)
	ID() string
	Weight() int
}

type worker struct {
	id     string
	pub    publisher.Publisher
	weight int
}

func New(ID string, pub publisher.Publisher, weight int) Worker {
	return &worker{
		id:     ID,
		pub:    pub,
		weight: weight,
	}
}

func (w worker) Do(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) {
	wg := &sync.WaitGroup{}
	wg.Add(wl.RequestsCount)

	log.
		WithField("worker", w.id).
		Tracef("request count: %d", wl.RequestsCount)

	for i := 0; i < wl.RequestsCount; i++ {
		go func() {
			defer wg.Done()

			met, err := w.pub.Publish(ctx, wl.Cfg)
			if errors.Is(err, model.ErrProtNotSupported) {
				log.Debug(err)
				return
			}

			if met != nil {
				metChan <- met
			}
		}()
	}

	wg.Wait()
}

func (w worker) ID() string {
	return w.id
}

func (w worker) Weight() int {
	return w.weight
}
