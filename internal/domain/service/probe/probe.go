package probe

import (
	"context"
	"io"
	"sync"

	"github.com/barbosaigor/nuker/internal/domain/service/bufwriter"
	"github.com/barbosaigor/nuker/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

// Probe listens to metrics incoming and operates metrics info
type Probe interface {
	Listen(ctx context.Context, met <-chan *metrics.NetworkMetrics) error
}

type probe struct {
	logger  bufwriter.BufWriter
	metRate *MetricRate
	mut     *sync.Mutex
}

func New(logger bufwriter.BufWriter) Probe {
	return &probe{
		logger:  logger,
		metRate: &MetricRate{},
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
			return nil
		case m, ok := <-met:
			if !ok {
				log.Trace("probe: metrics channel close")
				return nil
			}
			_ = p.writeMetr(ctx, m)
			p.metRate.Append(m)
		}
	}
}

func (p *probe) writeMetr(ctx context.Context, m *metrics.NetworkMetrics) error {
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
		Info(p.metRate.String())
}
