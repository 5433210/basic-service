package authn

import (
	fiber "github.com/gofiber/fiber/v2"
	servicev1 "wailik.com/internal/authn/service/v1"
	"wailik.com/internal/pkg/microservice"
	"wailik.com/internal/pkg/server"
)

type AuthzServer struct {
	server.Server
	service servicev1.Service
	Dsn     string
}

func CreateService(s *AuthzServer) (*AuthzServer, error) {
	service, err := servicev1.New(s.Dsn)
	if err != nil {
		return nil, err
	}
	s.service = service
	return s, nil
}

func (s *AuthzServer) SetMicroService(ms microservice.MicroService) {
	s.service.SetMicroService(ms)
}

func (s *AuthzServer) Bind(app *fiber.App) {
	route(app, s.service)
}
