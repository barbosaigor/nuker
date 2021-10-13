package orchestrator

import (
	"context"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/requester"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

type Orchestrator interface {
	WorkloadAssigner
	AddWorker(ID string, weight int)
	DelWorker(ID string)
	TotalWorkers() int
}

type orchestrator struct {
	workers     map[string]worker.Worker
	totalWeight int
	mut         *sync.RWMutex
	reqFactory  requester.Factory
}

func New(workers map[string]worker.Worker, reqFactory requester.Factory) Orchestrator {
	if workers == nil {
		workers = map[string]worker.Worker{}
	}

	return &orchestrator{
		workers:     workers,
		totalWeight: sumWeights(workers),
		mut:         &sync.RWMutex{},
		reqFactory:  reqFactory,
	}
}

func sumWeights(workers map[string]worker.Worker) int {
	total := 0
	for _, w := range workers {
		total += w.Weight()
	}
	return total
}

func (o *orchestrator) AddWorker(ID string, weight int) {
	o.mut.Lock()
	defer o.mut.Unlock()

	_, ok := o.workers[ID]
	if ok {
		log.
			WithField("worker-id", ID).
			Infof("worker already registered")
		return
	}

	o.totalWeight += weight
	o.workers[ID] = worker.New(ID, weight, o.reqFactory.Create())
}

func (o *orchestrator) DelWorker(ID string) {
	o.mut.Lock()
	defer o.mut.Unlock()

	_, ok := o.workers[ID]
	if !ok {
		return
	}

	o.totalWeight -= o.workers[ID].Weight()
	delete(o.workers, ID)

	log.
		WithField("worker-id", ID).
		Info("worker deleted")
}

// AssignWorkload distribute workload among workers
func (o orchestrator) AssignWorkload(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error {
	wg := &sync.WaitGroup{}
	o.mut.RLock()
	wg.Add(len(o.workers))
	o.mut.RUnlock()

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
	o.mut.RLock()
	defer o.mut.RUnlock()
	if len(o.workers) == 1 {
		return total
	}
	ratio := float64(o.workers[wID].Weight()) / float64(o.totalWeight)
	return int(float64(total) * ratio)
}

func (o orchestrator) TotalWorkers() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return len(o.workers)
}
