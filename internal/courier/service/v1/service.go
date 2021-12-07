package servicev1

import (
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
)

type Service interface {
	Email() *emailSrvc
	Sms() *smsSrvc
	microservice.MicroServiceHelper
}

type service struct {
	microservice.MicroServiceObject
}

var _ Service = &service{}

func (s *service) Email() *emailSrvc { return newEmailSrvc(s) }
func (s *service) Sms() *smsSrvc     { return newSmsSrvc(s) }

func New() (Service, error) {
	log.Info("new service...")

	var servie Service = &service{}

	log.Info("service created")

	return servie, nil
}
