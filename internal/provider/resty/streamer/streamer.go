package streamer

import (
	"context"
	"strings"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/barbosaigor/nuker/pkg/config"
	"github.com/barbosaigor/nuker/pkg/metrics"
	"github.com/go-resty/resty/v2"
)

type streamer struct {
	client *resty.Client
}

func New(client *resty.Client) repository.Streamer {
	return &streamer{
		client: client,
	}
}

func (s streamer) Stream(ctx context.Context, cfg config.Network) (*metrics.NetworkMetrics, error) {

	timeout := 3 * time.Second
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Second
	}

	res, err := s.client.
		SetTimeout(timeout).
		R().
		SetBody(cfg.Body).
		SetHeaders(cfg.Headers).
		Execute(strings.ToUpper(cfg.Method), cfg.Host+cfg.Path)

	met := metrics.NetworkMetrics{
		Host:         cfg.Host,
		Path:         cfg.Path,
		StatusCode:   res.StatusCode(),
		Body:         string(res.Body()),
		ResponseTime: res.Time(),
	}

	if err != nil {
		met.Err = err.Error()
		return &met, err
	}

	return &met, nil
}
