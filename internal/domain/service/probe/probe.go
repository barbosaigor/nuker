package probe

import (
	"context"
	"io"

	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/rs/zerolog/log"
)

// Probe listens to metrics incoming and operate metrics info
type Probe interface {
	Listen(ctx context.Context, met <-chan *metrics.NetworkMetrics) error
}

type probe struct {
	// logReader write metrics info into either a file, stdout, networking or etc.
	logger  repository.BufWriter
	metRate *MetricRate
}

func New(logger repository.BufWriter) Probe {
	return &probe{
		logger:  logger,
		metRate: &MetricRate{},
	}
}

func (p *probe) Listen(ctx context.Context, met <-chan *metrics.NetworkMetrics) error {
	log.Ctx(ctx).Trace().Msg("probe: listening")

	for {
		select {
		case <-ctx.Done():
			log.Ctx(ctx).Trace().Msg("probe: context close")
			p.metrReport(ctx)
			return nil
		case m, ok := <-met:
			if !ok {
				log.Ctx(ctx).Trace().Msg("probe: metrics channel close")
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

	log.Ctx(ctx).Debug().Msg("metric: " + m.String())
	n, err := io.WriteString(p.logger, m.String())
	if err != nil {
		log.Ctx(ctx).Error().Msg("Error to write metric in probe writer")
	}

	log.Ctx(ctx).Debug().Msgf("metrics with %d bytes wrote to writer", n)

	return err
}

func (p *probe) metrReport(ctx context.Context) {
	log.
		Ctx(ctx).
		Info().
		Str("log-location", p.logger.Location()).
		Msg(p.metRate.String())
}
