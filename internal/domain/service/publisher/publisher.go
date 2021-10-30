package publisher

import (
	"context"
	"strings"

	"github.com/barbosaigor/nuker/internal/domain/model"
	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/barbosaigor/nuker/pkg/metrics"
)

type Publisher interface {
	Publish(ctx context.Context, cfg config.Network) (*metrics.NetworkMetrics, error)
}

type publisher struct {
	httpSteamer repository.Streamer
}

func New(httpSteamer repository.Streamer) Publisher {
	return &publisher{
		httpSteamer: httpSteamer,
	}
}

func (p *publisher) Publish(ctx context.Context, cfg config.Network) (*metrics.NetworkMetrics, error) {
	if cfg.Protocol == "" || strings.EqualFold(cfg.Protocol, "http") || strings.EqualFold(cfg.Protocol, "https") {
		return p.httpSteamer.Stream(ctx, cfg)
	}

	return nil, model.ErrProtNotSupported
}
