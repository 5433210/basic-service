package courier

import (
	fiber "github.com/gofiber/fiber/v2"
	servicev1 "wailik.com/internal/courier/service/v1"
	"wailik.com/internal/pkg/microservice"
	"wailik.com/internal/pkg/server"
)

type CourierServer struct {
	server.Server
	service servicev1.Service
}

func CreateService(s *CourierServer) (*CourierServer, error) {
	service, err := servicev1.New()
	if err != nil {
		return nil, err
	}
	s.service = service
	return s, nil
}

func (s *CourierServer) SetMicroService(ms microservice.MicroService) {
	s.service.SetMicroService(ms)
}

func (s *CourierServer) Bind(app *fiber.App) {
	route(app, s.service)
}
