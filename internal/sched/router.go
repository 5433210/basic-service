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
	v1Sched := v1.Group("jobs")
	{
		v1Sched.Get("", c.Sched().GetAllJobs)
		v1Sched.Post("", c.Sched().CreateJob)
		v1Sched.Get(":job_id", c.Sched().GetJobById)
		v1Sched.Patch(":job_id", c.Sched().UpdateJob)
		v1Sched.Delete(":job_id", c.Sched().DeleteJob)
		v1Sched.Get(":job_id/executions", c.Sched().GetJobAllExecutions)
	}
}
