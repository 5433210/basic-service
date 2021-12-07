package controllerv1

import (
	servicev1 "wailik.com/internal/authn/service/v1"
)

type Controller interface {
	Identities() *identityController
	Identifiers() *identifierController
}

type controller struct {
	srvc servicev1.Service
}

func NewController(s servicev1.Service) Controller {
	var c Controller = &controller{
		srvc: s,
	}

	return c
}

func (c *controller) Identities() *identityController    { return newIdentityController(c) }
func (c *controller) Identifiers() *identifierController { return newIdentifierController(c) }
