package pipeline

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/barbosaigor/nuker/internal/domain/service/probe"
	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"golang.org/x/sync/errgroup"
)

type Pipeline interface {
	Run(ctx context.Context, cfg config.Config) error
}

type pipeline struct {
	opts     Options
	probeSvc probe.Probe
	orqSvc   orchestrator.Orchestrator
}

func New(probeSvc probe.Probe, orqSvc orchestrator.Orchestrator, opts Options) Pipeline {
	return &pipeline{
		probeSvc: probeSvc,
		orqSvc:   orqSvc,
		opts:     opts,
	}
}

func (p *pipeline) Run(ctx context.Context, cfg config.Config) (err error) {
	log.Debug("starting pipeline")

	go p.orqSvc.Listen()

	metChan := make(chan *metrics.NetworkMetrics)
	defer close(metChan)

	cCtx, cancelCtx := context.WithCancel(ctx)

	errG := &errgroup.Group{}

	errG.Go(func() error {
		return p.probeSvc.Listen(cCtx, metChan)
	})

	errG.Go(func() error {
		defer cancelCtx()

		p.run(cCtx, cfg, metChan)

		return nil
	})

	return errG.Wait()
}

func (p *pipeline) run(ctx context.Context, cfg config.Config, metChan chan<- *metrics.NetworkMetrics) {
	log.Tracef("%+v", cfg)

	for _, stg := range cfg.Stages {
		log.Info("running stage " + stg.Name)

		for _, step := range stg.Steps {
			log.Info("running step " + step.Name)

			stepWg := &sync.WaitGroup{}
			stepWg.Add(len(step.Containers))

			for _, container := range step.Containers {
				container := container

				go func() {
					defer stepWg.Done()

					log.Info("running container " + container.Name)

					p.startTicker(ctx, container, metChan)
				}()
			}

			stepWg.Wait()
		}
	}
}

func (p *pipeline) startTicker(ctx context.Context, container config.Container, metChan chan<- *metrics.NetworkMetrics) {
	tCtx, cancelFn := context.WithTimeout(ctx, time.Second*time.Duration(container.Duration+container.HoldFor))
	defer cancelFn()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	startTime := time.Duration(time.Now().UnixNano())
	for {
		select {
		case <-tCtx.Done():
			log.Trace("ctx timeout")
			return
		case <-ticker.C:
			log.Trace("do request")
			p.runContainer(tCtx, startTime, container, metChan)
		}
	}
}

func (p *pipeline) runContainer(
	ctx context.Context,
	startTime time.Duration,
	container config.Container,
	metChan chan<- *metrics.NetworkMetrics) {

	currTime := time.Duration(time.Now().UnixNano())
	endTime := time.Duration(container.Duration * int(time.Second))
	reqCount := p.calcRequests(startTime, endTime, currTime, container.Min, container.Max)

	log.Tracef("request count: %d", reqCount)

	wl := model.Workload{
		RequestsCount: reqCount,
		Cfg:           container.Network,
	}

	_ = p.orqSvc.AssignWorkload(ctx, wl, metChan)
}

func (p *pipeline) calcRequests(start, end, curr time.Duration, min, max int) int {
	if end <= 0 || min < 0 {
		return 0
	}

	if max == 0 || max < min {
		max = min
	}

	a := (float64(curr-start) / float64(time.Second)) * float64(max-min)
	b := float64(end-start) / float64(time.Second)
	if b == 0 {
		return 0
	}

	requests := min + int(a/b)

	if int(requests) > max {
		return max
	}

	return requests
}
