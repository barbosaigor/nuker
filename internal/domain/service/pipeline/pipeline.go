package pipeline

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/barbosaigor/nuker/pkg/config"
)

type Pipeline interface {
	Run(ctx context.Context, cfg config.Config) error
}

type pipeline struct {
	orqSvc orchestrator.Orchestrator
}

func New(orqSvc orchestrator.Orchestrator) Pipeline {
	return &pipeline{
		orqSvc: orqSvc,
	}
}

func (p *pipeline) Run(ctx context.Context, cfg config.Config) (err error) {
	log.Trace("starting pipeline")
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

					p.startTicker(ctx, container)
				}()
			}

			stepWg.Wait()
		}
	}

	log.Trace("pipeline finished")
	return nil
}

func (p *pipeline) startTicker(ctx context.Context, container config.Container) {
	totalDuration := time.Duration(container.Duration+container.HoldFor) * time.Second
	tCtx, cancelFn := context.WithTimeout(ctx, totalDuration)
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
			p.runContainer(tCtx, startTime, totalDuration, container)
		}
	}
}

func (p *pipeline) runContainer(
	ctx context.Context,
	startTime, totalDuration time.Duration,
	container config.Container) {

	currTime := time.Duration(time.Now().UnixNano())
	endTime := startTime + totalDuration
	reqCount := p.calcRequests(startTime, endTime, currTime, container.Min, container.Max)

	log.
		WithField("container", container.Name).
		Infof("request count: %d", reqCount)

	wl := model.Workload{
		RequestsCount: reqCount,
		Cfg:           container.Network,
	}

	p.orqSvc.DistributeWorkload(ctx, wl)
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
