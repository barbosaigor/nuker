package metrics

import (
	"encoding/json"
	"time"
)

type NetworkMetrics struct {
	Host         string
	Path         string
	StatusCode   int
	ResponseTime time.Duration
	Body         string
	Err          string
}

func (n NetworkMetrics) String() string {
	metrics, _ := json.Marshal(n)
	return string(metrics)
}
