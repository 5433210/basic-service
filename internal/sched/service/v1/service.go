package servicev1

import (
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
)

type Service interface {
	Sched() *schedSrvc
	microservice.MicroServiceHelper
}

type service struct {
	microservice.MicroServiceObject
}

var _ Service = &service{}

func (s *service) Sched() *schedSrvc { return newSchedSrvc(s) }

func New() (Service, error) {
	log.Info("new service...")

	var servie Service = &service{}

	log.Info("service created")

	return servie, nil
}
