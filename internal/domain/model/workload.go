package model

import "github.com/barbosaigor/nuker/pkg/config"

type Workload struct {
	RequestsCount int            `json:"requestCount"`
	Cfg           config.Network `json:"cfg"`
}
