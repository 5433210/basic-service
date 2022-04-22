package servicev1

import (
	storev1 "wailik.com/internal/authn/store/v1"
	"wailik.com/internal/pkg/cache"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/microservice"
)

type Service interface {
	Identity() *identitySrvc
	Identifier() *identifierSrvc
	Credential() *credentialSrvc
	microservice.MicroServiceHelper
}

type service struct {
	store storev1.Store
	cache cache.Cache
	microservice.MicroServiceObject
}

var _ Service = &service{}

func (s *service) Identity() *identitySrvc     { return newIdentitySrvc(s) }
func (s *service) Identifier() *identifierSrvc { return newIdentifierSrvc(s) }
func (s *service) Credential() *credentialSrvc { return newCredentialSrvc(s) }

func New(dsn string) (Service, error) {
	log.Infof("new service...%s", dsn)

	store, err := storev1.New(&storev1.Options{DSN: dsn})
	if err != nil {
		return nil, err
	}

	cache, err := cache.New(cache.Options{Type: "local"})
	if err != nil {
		return nil, err
	}

	var servie Service = &service{
		store: store,
		cache: cache,
	}

	log.Info("service created")

	return servie, nil
}
