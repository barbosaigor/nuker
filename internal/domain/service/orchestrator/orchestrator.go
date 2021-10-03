package orchestrator

import (
	"context"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	"github.com/barbosaigor/nuker/internal/provider/resty/requester"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/barbosaigor/nuker/pkg/net"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
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
	opts        Options
}

type workerBody struct {
	Weight int
}

func New(workers map[string]worker.Worker, opts Options) Orchestrator {
	if workers == nil {
		workers = map[string]worker.Worker{}
	}

	return &orchestrator{
		workers:     workers,
		totalWeight: sumWeights(workers),
		server:      fiber.New(),
		mut:         &sync.RWMutex{},
		opts:        opts,
	}
}

func sumWeights(workers map[string]worker.Worker) int {
	total := 0
	for _, w := range workers {
		total += w.Weight()
	}
	return total
}

// TODO: Move to another layer, provider logic
// Listen to workers creation and deletion events
func (o *orchestrator) Listen() error {
	log.Infof("master %s:%s", net.IP(), o.opts.Port)

	o.server.Post("/worker/:id", func(c *fiber.Ctx) error {
		workerID := string(append([]byte{}, c.Params("id")[:]...))
		wb, _ := o.parseWorkerBody(c)

		o.addWorker(workerID, wb.Weight)

		return nil
	})

	o.server.Delete("/worker/:id", func(c *fiber.Ctx) error {
		workerID := string(append([]byte{}, c.Params("id")[:]...))

		o.delWorker(workerID)

		return nil
	})

	return o.server.Listen(":" + o.opts.Port)
}

func (o *orchestrator) addWorker(ID string, weight int) {
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
	o.workers[ID] = worker.New(ID, weight, requester.New())
}

func (o *orchestrator) delWorker(ID string) {
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

func (o orchestrator) parseWorkerBody(c *fiber.Ctx) (workerBody, error) {
	var wb workerBody
	err := c.BodyParser(&wb)

	if wb.Weight < 1 {
		wb.Weight = 1
	}

	return wb, err
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
