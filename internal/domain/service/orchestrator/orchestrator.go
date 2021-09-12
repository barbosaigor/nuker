package orchestrator

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Orchestrator interface {
	Listen() error
	AssignWorkload(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error
}

type orchestrator struct {
	workers     map[string]worker.Worker
	totalWeight int
	server      *fiber.App
	mut         *sync.RWMutex
}

func New(workers map[string]worker.Worker) Orchestrator {
	return &orchestrator{
		workers:     workers,
		totalWeight: sumWeights(workers),
		server:      fiber.New(),
		mut:         &sync.RWMutex{},
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
	o.server.Post("/worker/:id", func(c *fiber.Ctx) error {
		workerID := c.Params("id")
		o.mut.RLock()
		_, ok := o.workers[workerID]
		o.mut.RUnlock()
		if !ok {
			logrus.
				WithField("worker-id", workerID).
				Infof("worker %s:%s registered", c.IP(), c.Port())
			o.mut.Lock()
			o.workers[workerID] = worker.New(workerID, 1, nil)
			o.mut.Unlock()
		}

		c.Status(http.StatusCreated)
		return c.SendString(fmt.Sprintf("Hello %s, World 👋!", workerID))
	})

	return o.server.Listen(":9050")
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
