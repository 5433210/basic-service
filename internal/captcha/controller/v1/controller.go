package controllerv1

import (
	servicev1 "wailik.com/internal/captcha/service/v1"
)

type Controller interface {
	Captcha() *captchaController
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

func (c *controller) Captcha() *captchaController { return newCaptchaController(c) }
