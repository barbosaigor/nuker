package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/model"
	m "github.com/barbosaigor/nuker/internal/domain/service/master"
	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	"github.com/barbosaigor/nuker/internal/domain/service/probe"
	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/barbosaigor/nuker/pkg/net"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type master struct {
	probeSvc    probe.Probe
	orchSvc     orchestrator.Orchestrator
	pipelineSvc pipeline.Pipeline
	server      *fiber.App
	done        bool
	opts        Options
}

func NewMaster(probe probe.Probe, orch orchestrator.Orchestrator, pipeline pipeline.Pipeline, opts Options) m.Master {
	return &master{
		probeSvc:    probe,
		orchSvc:     orch,
		pipelineSvc: pipeline,
		server: fiber.New(fiber.Config{
			DisableKeepalive:      true,
			DisableStartupMessage: true,
		}),
		done: false,
		opts: opts,
	}
}

func (m *master) Run(ctx context.Context, cfg config.Config) (err error) {
	metChan := make(chan *metrics.NetworkMetrics)
	defer close(metChan)

	cCtx, cancelCtx := context.WithCancel(ctx)

	errG := &errgroup.Group{}

	errG.Go(func() error {
		defer cancelCtx()
		defer m.server.Shutdown()
		return m.probeSvc.Listen(cCtx, metChan)
	})

	errG.Go(func() error {
		defer cancelCtx()
		return m.listen(cCtx, metChan)
	})

	errG.Go(func() error {
		defer cancelCtx()
		defer func() {
			log.Trace("waiting for graceful worker shutdowns")
			m.done = true
			select {
			case <-time.After(1 * time.Minute):
				return
			case <-m.isDrained():
				log.Trace("awaiting for remaining workers metric")
				<-time.After(5 * time.Second)
				return
			}
		}()

		// start pipeline when there are enough workers assigned
		for {
			select {
			case <-cCtx.Done():
				return nil
			default:
				if m.opts.MinWorkers <= m.orchSvc.TotalWorkers() {
					return m.pipelineSvc.Run(cCtx, cfg)
				}
				<-time.After(time.Second)
			}
		}
	})

	return errG.Wait()
}

func (m *master) isDrained() <-chan struct{} {
	drainCh := make(chan struct{})
	go func() {
		for {
			if !m.orchSvc.HasAnyWorkload() {
				drainCh <- struct{}{}
			}
			<-time.After(100 * time.Millisecond)
		}
	}()
	return drainCh
}

func (m *master) listen(ctx context.Context, metChan chan<- *metrics.NetworkMetrics) error {
	log.Infof("master URL: http://%s:%s", net.IP(), m.opts.Port)
	// TODO: refactor, extract handler logic to separate function
	m.server.Post("/worker/:id", func(c *fiber.Ctx) error {
		wb, err := m.parseWorkerBody(c)
		if err != nil {
			return err
		}

		wb.ID = utils.ImmutableString(c.Params("id"))

		m.orchSvc.AddWorker(wb.ID, wb.Weight)

		resp, err := json.Marshal(wb)
		if err != nil {
			return err
		}

		c.Send(resp)
		return nil
	})

	m.server.Post("/worker", func(c *fiber.Ctx) error {
		wb, err := m.parseWorkerBody(c)
		if err != nil {
			return err
		}

		wb.ID = uuid.NewString()

		m.orchSvc.AddWorker(wb.ID, wb.Weight)

		resp, err := json.Marshal(wb)
		if err != nil {
			return err
		}

		c.Set("Content-Type", "application/json")
		c.Send(resp)

		return nil
	})

	m.server.Delete("/worker/:id", func(c *fiber.Ctx) error {
		workerID := utils.ImmutableString(c.Params("id"))

		m.orchSvc.DelWorker(workerID)

		return nil
	})

	m.server.Get("/worker/:id", func(c *fiber.Ctx) error {
		workerID := utils.ImmutableString(c.Params("id"))

		if !m.orchSvc.HasWorker(workerID) {
			log.Infof("worker %q not found", workerID)
			return c.SendStatus(http.StatusNotFound)
		}

		m.orchSvc.FlushWorker(workerID)

		wls := m.orchSvc.TakeWorkloads(workerID)
		if !m.done && len(wls) == 0 {
			return c.SendStatus(http.StatusNoContent)
		}

		var op model.WorkerOp
		if m.done && len(wls) == 0 {
			op = model.Detach
		} else {
			op = model.Assignment
		}

		lc := model.LaborContract{
			Operation: op,
			Workloads: wls,
		}

		resp, err := json.Marshal(lc)
		if err != nil {
			return err
		}

		c.Set("Content-Type", "application/json")
		c.Send(resp)

		return nil
	})

	m.server.Post("/worker/:id/metrics", func(c *fiber.Ctx) error {
		workerID := utils.ImmutableString(c.Params("id"))

		if !m.orchSvc.HasWorker(workerID) {
			return c.SendStatus(http.StatusNotFound)
		}

		var met metrics.NetworkMetrics

		err := c.BodyParser(&met)
		if err != nil {
			return err
		}

		metChan <- &met

		return nil
	})

	return m.server.Listen(":" + m.opts.Port)
}

func (m master) parseWorkerBody(c *fiber.Ctx) (model.WorkerBody, error) {
	var wb model.WorkerBody
	err := c.BodyParser(&wb)
	if err != nil {
		return wb, err
	}

	if wb.Weight < 1 {
		wb.Weight = 1
	}

	return wb, nil
}
