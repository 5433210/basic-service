package controllerv1

import (
	servicev1 "wailik.com/internal/sched/service/v1"
)

type Controller interface {
	Sched() *schedController
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

func (c *controller) Sched() *schedController { return newSchedController(c) }
