package model

import "github.com/barbosaigor/nuker/pkg/config"

type Workload struct {
	RequestsCount int
	Cfg           config.Network
}
