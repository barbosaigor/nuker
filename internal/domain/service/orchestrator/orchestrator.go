package orchestrator

import (
	"context"
	"sync"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/model"
	log "github.com/sirupsen/logrus"
)

type Orchestrator interface {
	DistributeWorkload(ctx context.Context, wl model.Workload)
	AddWorker(ID string, weight int)
	DelWorker(ID string)
	HasWorker(ID string) bool
	FlushWorker(ID string)
	TakeWorkload(ID string) (model.Workload, bool)
	TotalWorkers() int
}

type orchestrator struct {
	workers     map[string]wlWorker
	totalWeight int
	mut         *sync.RWMutex
}

type wlWorker struct {
	ID              string
	weight          int
	workload        model.Workload
	shouldProcessWL bool
	lastFlush       time.Time
}

func (ws wlWorker) clone() wlWorker {
	return ws
}

func New() Orchestrator {
	return &orchestrator{
		workers:     map[string]wlWorker{},
		totalWeight: 0,
		mut:         &sync.RWMutex{},
	}
}

func (o *orchestrator) HasWorker(ID string) bool {
	o.mut.RLock()
	defer o.mut.RUnlock()

	_, ok := o.workers[ID]
	return ok
}

func (o *orchestrator) FlushWorker(ID string) {
	if !o.HasWorker(ID) {
		return
	}

	o.mut.Lock()
	defer o.mut.Unlock()

	newWorker := o.workers[ID].clone()
	newWorker.lastFlush = time.Now()
	o.workers[ID] = newWorker

	log.
		WithField("id", newWorker.ID).
		WithField("weight", newWorker.weight).
		WithField("last-flush", newWorker.lastFlush.String()).
		Tracef("worker flushed")
}

func (o *orchestrator) GarbageCollectWorkers() {
	// TODO: tweak
	const workerTTL = 15 * time.Second
	// TODO: data race, should create workers copy
	for _, w := range o.workers {
		if time.Now().After(w.lastFlush.Add(workerTTL)) {
			log.
				WithField("elapsed-time", time.Until(w.lastFlush.Add(workerTTL)).String()).
				WithField("worker-id", w.ID).
				WithField("worker-last-flush", w.lastFlush.Add(workerTTL).String()).
				Tracef("garbage collecting worker")
			o.DelWorker(w.ID)
		}
	}
}

func (o orchestrator) TakeWorkload(ID string) (model.Workload, bool) {
	o.mut.Lock()
	defer o.mut.Unlock()

	w, ok := o.workers[ID]
	if !ok || !o.workers[ID].shouldProcessWL {
		return w.workload, false
	}

	wl := w.workload

	newWorker := o.workers[ID].clone()
	newWorker.shouldProcessWL = false
	o.workers[ID] = newWorker

	return wl, true
}

func (o *orchestrator) AddWorker(ID string, weight int) {
	log.
		WithField("id", ID).
		WithField("weight", weight).
		Tracef("adding worker")
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
	o.workers[ID] = wlWorker{
		ID:              ID,
		weight:          weight,
		shouldProcessWL: false,
		lastFlush:       time.Now(),
	}

	log.
		WithField("id", ID).
		WithField("weight", weight).
		WithField("last-flush", o.workers[ID].lastFlush.String()).
		Debug("worker added")
}

func (o *orchestrator) DelWorker(ID string) {
	o.mut.Lock()
	defer o.mut.Unlock()

	_, ok := o.workers[ID]
	if !ok {
		return
	}

	o.totalWeight -= o.workers[ID].weight
	delete(o.workers, ID)

	log.
		WithField("worker-id", ID).
		Debug("worker deleted")
}

// DistributeWorkload calculates workload among workers
func (o orchestrator) DistributeWorkload(ctx context.Context, wl model.Workload) {
	o.GarbageCollectWorkers()

	// TODO: data race, create workers copy
	for _, w := range o.workers {
		workerWL := model.Workload{
			RequestsCount: o.calcRequests(w.ID, wl.RequestsCount),
			Cfg:           wl.Cfg,
		}

		// log.
		// 	WithField("worker-id", w.ID).
		// 	Tracef("worker workload: %v", workerWL)

		o.mut.Lock()
		newWorker := w.clone()
		newWorker.workload = workerWL
		newWorker.shouldProcessWL = true
		o.workers[w.ID] = newWorker
		o.mut.Unlock()
	}
}

// calcRequests calculate request amount for a specific worker
func (o orchestrator) calcRequests(wID string, total int) int {
	o.mut.RLock()
	defer o.mut.RUnlock()

	if len(o.workers) <= 1 {
		return total
	}

	ratio := float64(o.workers[wID].weight) / float64(o.totalWeight)
	return int(float64(total) * ratio)
}

func (o orchestrator) TotalWorkers() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return len(o.workers)
}
