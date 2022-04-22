package servicev1

import (
	"wailik.com/internal/pkg/cache"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
)

type Service interface {
	microservice.MicroServiceHelper
	Captcha() *captchaSrvc
}

type service struct {
	microservice.MicroServiceObject
	cache cache.Cache
}

var _ Service = &service{}

func (s *service) Captcha() *captchaSrvc { return newCaptchaSrvc(s) }

func New() (Service, error) {
	log.Info("new service...")
	cache, err := cache.New(cache.Options{Type: "local"})
	if err != nil {
		return nil, err
	}

	var servie Service = &service{cache: cache}

	log.Info("service created")

	return servie, nil
}
