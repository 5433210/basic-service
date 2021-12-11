package servicev1

import (
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
	"wailik.com/internal/sched/store"
)

type Service interface {
	Sched() *schedSrvc
	microservice.MicroServiceHelper
}

type service struct {
	microservice.MicroServiceObject
	store *store.Store
}

var _ Service = &service{}

func (s *service) Sched() *schedSrvc { return newSchedSrvc(s) }

func New(endpoint []string, size int) (Service, error) {
	log.Info("new service...")

	store, err := store.NewStore(endpoint, size)
	if err != nil {
		return nil, err
	}

	var servie Service = &service{
		store: store,
	}

	log.Info("service created")

	return servie, nil
}
