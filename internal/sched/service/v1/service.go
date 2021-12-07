package servicev1

import (
	"wailik.com/internal/pkg/cache"
	"wailik.com/internal/pkg/log"
)

type Service interface {
	Sched() *schedSrvc
}

type service struct {
	cache cache.Cache
}

var _ Service = &service{}

func (s *service) Sched() *schedSrvc { return newSchedSrvc(s) }

func New() (Service, error) {
	log.Info("new service...")
	cache, err := cache.New(cache.Options{Type: "local"})
	if err != nil {
		return nil, err
	}

	var servie Service = &service{cache: cache}

	log.Info("service created")

	return servie, nil
}
