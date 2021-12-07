package controllerv1

import (
	servicev1 "wailik.com/internal/sched/service/v1"
)

type schedController struct {
	srvc servicev1.Service
}

func newSchedController(c *controller) *schedController {
	return &schedController{srvc: c.srvc}
}
