package local

import (
	"context"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/master"
	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	"github.com/barbosaigor/nuker/internal/domain/service/probe"
	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type localMaster struct {
	orchSvc     orchestrator.Orchestrator
	probeSvc    probe.Probe
	pubSvc      publisher.Publisher
	pipelineSvc pipeline.Pipeline
}

func NewMaster(
	orchSvc orchestrator.Orchestrator,
	probeSvc probe.Probe,
	pubSvc publisher.Publisher,
	pipelineSvc pipeline.Pipeline,
) master.Master {
	return &localMaster{
		orchSvc:     orchSvc,
		probeSvc:    probeSvc,
		pubSvc:      pubSvc,
		pipelineSvc: pipelineSvc,
	}
}

const masterWorkerID = "master-worker-1"

func (lm *localMaster) Run(ctx context.Context, cfg config.Config) error {
	lm.orchSvc.AddWorker(masterWorkerID, 1)

	metChan := make(chan *metrics.NetworkMetrics)
	defer close(metChan)

	cCtx, cancelCtx := context.WithCancel(ctx)
	wCtx, wCancelCtx := context.WithCancel(cCtx)

	errG := &errgroup.Group{}

	errG.Go(func() error {
		return lm.probeSvc.Listen(cCtx, metChan)
	})

	errG.Go(func() error {
		defer cancelCtx()
		lm.pollWorkloads(wCtx, metChan)
		return nil
	})

	errG.Go(func() error {
		defer wCancelCtx()
		return lm.pipelineSvc.Run(cCtx, cfg)
	})

	return errG.Wait()
}

func (lm *localMaster) pollWorkloads(ctx context.Context, metChan chan<- *metrics.NetworkMetrics) {
	for {
		select {
		case <-ctx.Done():
			log.Tracef("context canceled: %v", ctx.Err())
			return
		default:
			lm.orchSvc.FlushWorker(masterWorkerID)

			wls := lm.orchSvc.TakeWorkloads(masterWorkerID)
			if len(wls) == 0 {
				continue
			}

			lm.assignWls(ctx, wls, metChan)
		}
	}
}

func (lm *localMaster) assignWls(ctx context.Context, wls []model.Workload, metChan chan<- *metrics.NetworkMetrics) {
	wg := &sync.WaitGroup{}
	wg.Add(len(wls))

	for _, wl := range wls {
		wl := wl

		go func() {
			defer wg.Done()
			lm.doWl(ctx, wl, metChan)
		}()
	}

	wg.Wait()
}

func (lm *localMaster) doWl(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) {
	log.Debugf("running %d requests", wl.RequestsCount)

	wg := &sync.WaitGroup{}
	wg.Add(wl.RequestsCount)

	for i := 0; i < wl.RequestsCount; i++ {
		go func() {
			defer wg.Done()

			met, err := lm.pubSvc.Publish(ctx, wl.Cfg)
			if err != nil {
				log.Trace(err)
			}

			if met != nil {
				metChan <- met
			}
		}()
	}

	wg.Wait()
}
