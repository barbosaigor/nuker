package requester

import (
	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
)

type Factory interface {
	Create() repository.Requester
}

type factory struct {
	pub publisher.Publisher
}

func NewFactory(pub publisher.Publisher) Factory {
	return &factory{
		pub: pub,
	}
}

func (f factory) Create() repository.Requester {
	return New(f.pub)
}
