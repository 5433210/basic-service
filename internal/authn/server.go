package authn

import (
	fiber "github.com/gofiber/fiber/v2"
	servicev1 "wailik.com/internal/authn/service/v1"
	"wailik.com/internal/pkg/microservice"
	"wailik.com/internal/pkg/server"
)

type AuthnServer struct {
	server.Server
	service    servicev1.Service
	DataSource string
}

func CreateService(s *AuthnServer) (*AuthnServer, error) {
	service, err := servicev1.New(s.DataSource)
	if err != nil {
		return nil, err
	}
	s.service = service
	return s, nil
}

func (s *AuthnServer) SetMicroService(ms microservice.MicroService) {
	s.service.SetMicroService(ms)
}

func (s *AuthnServer) Bind(app *fiber.App) {
	route(app, s.service)
}
