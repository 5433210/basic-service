package sched

import (
	"github.com/gofiber/fiber/v2"

	controllerv1 "wailik.com/internal/sched/controller/v1"
	servicev1 "wailik.com/internal/sched/service/v1"
)

func route(app *fiber.App, srvc servicev1.Service) {
	c := controllerv1.NewController(srvc)

	v1 := app.Group("v1")
	{
		schedGroup(v1, c)
	}
}

func schedGroup(v1 fiber.Router, c controllerv1.Controller) {
	v1Sched := v1.Group("sched")
	{
		v1Sched.Post("", c.Sched().VerifySched)
		v1Sched.Get("", c.Sched().GenerateSched)
	}
}
