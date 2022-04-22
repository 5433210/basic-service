package sched

import (
	fiber "github.com/gofiber/fiber/v2"
	"wailik.com/internal/pkg/microservice"
	"wailik.com/internal/pkg/server"
	servicev1 "wailik.com/internal/sched/service/v1"
)

type SchedServer struct {
	server.Server
	service       servicev1.Service
	StoreEndpoint []string
	StorePoolSize int
}

func CreateService(s *SchedServer) (*SchedServer, error) {
	service, err := servicev1.New(s.StoreEndpoint, s.StorePoolSize)
	if err != nil {
		return nil, err
	}
	s.service = service
	return s, nil
}

func (s *SchedServer) SetMicroService(ms microservice.MicroService) {
	s.service.SetMicroService(ms)
}

func (s *SchedServer) Bind(app *fiber.App) {
	route(app, s.service)
}
