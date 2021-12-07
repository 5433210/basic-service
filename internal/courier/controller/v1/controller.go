package controllerv1

import (
	servicev1 "wailik.com/internal/courier/service/v1"
)

type Controller interface {
	Email() *emailController
	Sms() *smsController
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

func (c *controller) Email() *emailController { return newEmailController(c) }
func (c *controller) Sms() *smsController     { return newSmsController(c) }
