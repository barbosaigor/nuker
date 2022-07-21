package probe

import (
	"context"
	"io"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/service/bufwriter"
	"github.com/barbosaigor/nuker/internal/domain/service/view"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

// Probe listens to metrics incoming and operates metrics info
type Probe interface {
	Listen(ctx context.Context, met <-chan *metrics.NetworkMetrics) error
}

type probe struct {
	logger  bufwriter.BufWriter
	vw      view.View
	metRate *metrics.MetricRate
	mut     *sync.Mutex
}

func New(logger bufwriter.BufWriter, vw view.View) Probe {
	return &probe{
		logger:  logger,
		vw:      vw,
		metRate: &metrics.MetricRate{},
		mut:     &sync.Mutex{},
	}
}

func (p *probe) Listen(ctx context.Context, met <-chan *metrics.NetworkMetrics) error {
	log.Trace("probe: listening")

	for {
		select {
		case <-ctx.Done():
			log.Trace("probe: context close")
			p.metrReport(ctx)
			p.vw.ShutDown()
			return nil
		case m, ok := <-met:
			if !ok {
				log.Trace("probe: metrics channel close")
				return nil
			}
			_ = p.writeMetrToLogger(ctx, m)
			p.metRate.Append(m)
			p.vw.SetMetric(p.metRate)
		}
	}
}

func (p *probe) writeMetrToLogger(ctx context.Context, m *metrics.NetworkMetrics) error {
	if m == nil {
		return nil
	}

	log.Trace("metric: " + m.String())
	p.mut.Lock()
	defer p.mut.Unlock()
	n, err := io.WriteString(p.logger, m.String())
	if err != nil {
		log.Error("Error to write metric in probe writer")
	}

	log.Tracef("metrics with %d bytes wrote to writer", n)

	return err
}

func (p *probe) metrReport(ctx context.Context) {
	log.
		WithField("logPath", p.logger.Location()).
		Debug(p.metRate.String())
}
