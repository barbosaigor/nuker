package pipeline

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/probe"
	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"golang.org/x/sync/errgroup"
)

type Pipeline interface {
	Run(ctx context.Context, cfg config.Config) error
}

type pipeline struct {
	pubSvc   publisher.Publisher
	probeSvc probe.Probe
}

func New(pub publisher.Publisher, probeSvc probe.Probe) Pipeline {
	return &pipeline{
		pubSvc:   pub,
		probeSvc: probeSvc,
	}
}

func (p *pipeline) Run(ctx context.Context, cfg config.Config) (err error) {
	log.Ctx(ctx).Debug().Msg("starting pipeline")

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
	log.Ctx(ctx).Trace().Msgf("%+v", cfg)

	for _, stg := range cfg.Stages {
		log.Ctx(ctx).Info().Msg("running stage " + stg.Name)

		for _, step := range stg.Steps {
			log.Ctx(ctx).Info().Msg("running step " + step.Name)

			stepWg := &sync.WaitGroup{}
			stepWg.Add(len(step.Containers))

			for _, container := range step.Containers {
				container := container

				go func() {
					defer stepWg.Done()

					log.Ctx(ctx).Info().Msg("running container " + container.Name)

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
			log.Ctx(ctx).Trace().Msg("ctx timeout")
			return
		case <-ticker.C:
			log.Ctx(ctx).Trace().Msg("do request")
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
	endTime := time.Duration(container.Duration)
	reqCount := p.calcRequests(startTime, endTime, currTime, container.Min, container.Max)

	wg := &sync.WaitGroup{}
	wg.Add(reqCount)

	log.Ctx(ctx).Trace().Msgf("request count: %d", reqCount)

	for i := 0; i < reqCount; i++ {
		go func() {
			defer wg.Done()

			met, err := p.pubSvc.Publish(ctx, container.Network)
			if errors.Is(err, model.ErrProtNotSupported) {
				log.Ctx(ctx).Debug().Err(err)
				return
			}

			if met != nil {
				metChan <- met
			}
		}()
	}

	wg.Wait()
}

func (p *pipeline) calcRequests(start, end, curr time.Duration, min, max int) int {
	if max == 0 {
		max = min
	}

	dt := float64((curr - start) / time.Second)
	ratio := dt / float64(end)
	dr := ratio * float64(max)

	if int(dr) > max {
		return max
	}

	return int(dr)
}
