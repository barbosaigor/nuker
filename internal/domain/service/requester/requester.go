package requester

import (
	"context"
	"errors"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

type Requester interface {
	Assign(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error
}

type requester struct {
	pub publisher.Publisher
}

func New(pub publisher.Publisher) Requester {
	return &requester{
		pub: pub,
	}
}

func (c requester) Assign(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error {
	wg := &sync.WaitGroup{}
	wg.Add(wl.RequestsCount)

	for i := 0; i < wl.RequestsCount; i++ {
		go func() {
			defer wg.Done()

			met, err := c.pub.Publish(ctx, wl.Cfg)
			if errors.Is(err, model.ErrProtNotSupported) {
				log.Trace(err)
				return
			}

			if met != nil {
				metChan <- met
			}
		}()
	}

	wg.Wait()

	return nil
}
