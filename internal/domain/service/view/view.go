package view

import "github.com/barbosaigor/nuker/pkg/metrics"

type View interface {
	SetMetric(mr *metrics.MetricRate)
	Start()
	ShutDown()
}
