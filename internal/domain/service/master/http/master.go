package http

import (
	"context"
	"time"

	m "github.com/barbosaigor/nuker/internal/domain/service/master"
	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/barbosaigor/nuker/pkg/net"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type master struct {
	orchSvc     orchestrator.Orchestrator
	pipelineSvc pipeline.Pipeline
	server      *fiber.App
	opts        Options
}

func NewMaster(orch orchestrator.Orchestrator, pipeline pipeline.Pipeline, opts Options) m.Master {
	return &master{
		orchSvc:     orch,
		pipelineSvc: pipeline,
		server:      fiber.New(),
		opts:        opts,
	}
}

func (m *master) Run(ctx context.Context, cfg config.Config) (err error) {
	go func() {
		err = m.listen()
	}()

	// wait and start pipeline when there are enough workers assigned
	for {
		if m.opts.MinWorkers <= m.orchSvc.TotalWorkers() {
			return m.pipelineSvc.Run(ctx, cfg)
		}
		<-time.After(time.Second)
	}
}

func (m *master) listen() error {
	log.Infof("master %s:%s", net.IP(), m.opts.Port)

	m.server.Post("/worker/:id", func(c *fiber.Ctx) error {
		workerID := string(append([]byte{}, c.Params("id")[:]...))
		wb, _ := m.parseWorkerBody(c)

		m.orchSvc.AddWorker(workerID, wb.Weight)

		return nil
	})

	m.server.Delete("/worker/:id", func(c *fiber.Ctx) error {
		workerID := string(append([]byte{}, c.Params("id")[:]...))

		m.orchSvc.DelWorker(workerID)

		return nil
	})

	return m.server.Listen(":" + m.opts.Port)
}

type workerBody struct {
	Weight int
}

func (m master) parseWorkerBody(c *fiber.Ctx) (workerBody, error) {
	var wb workerBody
	err := c.BodyParser(&wb)

	if wb.Weight < 1 {
		wb.Weight = 1
	}

	return wb, err
}
