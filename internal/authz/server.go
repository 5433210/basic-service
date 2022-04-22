package authz

import (
	fiber "github.com/gofiber/fiber/v2"
	servicev1 "wailik.com/internal/authz/service/v1"
	"wailik.com/internal/pkg/microservice"
	"wailik.com/internal/pkg/server"
)

type AuthzServer struct {
	server.Server
	service  servicev1.Service
	RegoPath string
	DataPath string
	DBPath   string
}

func CreateService(s *AuthzServer) (*AuthzServer, error) {
	service, err := servicev1.New(s.DBPath, s.RegoPath, s.DataPath)
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
