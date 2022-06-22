package metrics

import (
	"encoding/json"
	"time"
)

type MetricRate struct {
	Success int
	Failed  int
	Total   int

	AvgSuccess float64
	AvgFailed  float64

	MaxTime time.Duration
	AvgTime time.Duration
	MinTime time.Duration

	timeSum time.Duration

	SampleSize int
}

func NewMetricRate(sampleSize int) *MetricRate {
	return &MetricRate{
		SampleSize: sampleSize,
	}
}

func (mr *MetricRate) Append(m *NetworkMetrics) {
	if m == nil {
		return
	}

	if mr.SampleSize > 0 && mr.Success+mr.Failed%mr.SampleSize == 0 {
		mr.timeSum = m.ResponseTime
	} else {
		mr.timeSum += m.ResponseTime
	}

	if m.ResponseTime > mr.MaxTime {
		mr.MaxTime = m.ResponseTime
	}

	if mr.Failed == 0 && mr.Success == 0 {
		mr.MinTime = m.ResponseTime
	} else if m.ResponseTime < mr.MinTime {
		mr.MinTime = m.ResponseTime
	}

	if m.StatusCode >= 200 && m.StatusCode < 300 {
		mr.Success++
	} else {
		mr.Failed++
	}

	mr.calcAvg()

	mr.Total = mr.Success + mr.Failed
}

func (mr *MetricRate) calcAvg() {
	total := mr.Success + mr.Failed
	if total == 0 {
		return
	}

	mr.AvgSuccess = float64(mr.Success) / float64(total)
	mr.AvgFailed = float64(mr.Failed) / float64(total)

	mr.AvgTime = mr.timeSum / time.Duration(total)
}

func (mr MetricRate) String() string {
	data, err := json.Marshal(mr)
	if err != nil {
		return ""
	}

	return string(data)
}
