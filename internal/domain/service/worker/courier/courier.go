package courier

import (
	"context"
	"errors"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

type courier struct {
	pub publisher.Publisher
}

func New(pub publisher.Publisher) repository.Courier {
	return &courier{
		pub: pub,
	}
}

func (c courier) Do(ctx context.Context, wl model.Workload, metChan chan<- *metrics.NetworkMetrics) error {
	wg := &sync.WaitGroup{}
	wg.Add(wl.RequestsCount)

	for i := 0; i < wl.RequestsCount; i++ {
		go func() {
			defer wg.Done()

			met, err := c.pub.Publish(ctx, wl.Cfg)
			if errors.Is(err, model.ErrProtNotSupported) {
				log.Debug(err)
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
