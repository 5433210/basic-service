package servicev1

import (
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
	"wailik.com/internal/sched/store"
)

type Service interface {
	Sched() *schedSrvc
	microservice.MicroServiceHelper
	Run()
	Stop()
	LoadSchedules() error
}

type service struct {
	microservice.MicroServiceObject
	store     *store.Store
	scheduler *Scheduler
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
		store:     store,
		scheduler: NewScheduler(store),
	}

	log.Info("service created")

	return servie, nil
}

func (s *service) Run() {
	s.GetMicroService().Start()
	s.scheduler.Start()
}

func (s *service) Stop() {
	s.scheduler.Stop()
}

func (s *service) LoadSchedules() error {
	return s.scheduler.Load()
}
