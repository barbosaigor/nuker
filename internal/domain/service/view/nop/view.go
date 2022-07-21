package nop

import (
	"github.com/barbosaigor/nuker/internal/domain/service/view"
	"github.com/barbosaigor/nuker/pkg/metrics"
)

type nopView struct{}

func New() (view.View, error) {
	return nopView{}, nil
}

func (nopView) SetMetric(mr *metrics.MetricRate) {}

func (nopView) ShutDown() {}
