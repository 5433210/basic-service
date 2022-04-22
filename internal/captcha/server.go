package captcha

import (
	fiber "github.com/gofiber/fiber/v2"
	servicev1 "wailik.com/internal/captcha/service/v1"
	"wailik.com/internal/pkg/microservice"
	"wailik.com/internal/pkg/server"
)

type SchedServer struct {
	server.Server
	service       servicev1.Service
	StoreEndpoint []string
	StorePoolSize int
}

func CreateService(s *SchedServer) (*SchedServer, error) {
	service, err := servicev1.New()
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
