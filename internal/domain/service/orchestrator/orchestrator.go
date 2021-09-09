package orchestrator

import (
	"context"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	"github.com/barbosaigor/nuker/pkg/metrics"
)

type Orchestrator interface {
	Listen() error
	AssignWorkload(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error
}

type orchestrator struct {
	workers     map[string]worker.Worker
	totalWeight int
}

func New(workers map[string]worker.Worker) Orchestrator {
	return &orchestrator{
		workers:     workers,
		totalWeight: sumWeights(workers),
	}
}

func sumWeights(workers map[string]worker.Worker) int {
	total := 0
	for _, w := range workers {
		total += w.Weight()
	}
	return total
}

// TODO
// Listen to new/delete workers
func (o *orchestrator) Listen() error {
	// start http server
	return nil
}

// AssignWorkload distribute workload among workers
func (o orchestrator) AssignWorkload(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(o.workers))

	for _, w := range o.workers {
		w := w
		wl := model.Workload{
			RequestsCount: o.calcRequests(w.ID(), wl.RequestsCount),
			Cfg:           wl.Cfg,
		}
		go func() {
			defer wg.Done()
			w.Do(ctx, wl, metChan)
		}()
	}

	wg.Wait()

	return nil
}

// calcRequests calculate request amount for a specific worker
func (o orchestrator) calcRequests(wID string, total int) int {
	if len(o.workers) == 1 {
		return total
	}

	ratio := float64(o.workers[wID].Weight()) / float64(o.totalWeight)
	return int(float64(total) * ratio)
}
