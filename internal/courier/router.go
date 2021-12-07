package courier

import (
	"github.com/gofiber/fiber/v2"

	controllerv1 "wailik.com/internal/courier/controller/v1"
	servicev1 "wailik.com/internal/courier/service/v1"
)

func route(app *fiber.App, srvc servicev1.Service) {
	c := controllerv1.NewController(srvc)

	v1 := app.Group("v1")
	{
		emailGroup(v1, c)
		smsGroup(v1, c)
	}
}

func emailGroup(v1 fiber.Router, c controllerv1.Controller) {
	v1Email := v1.Group("email")
	{
		v1Email.Post("", c.Email().SendEmail)
	}
}

func smsGroup(v1 fiber.Router, c controllerv1.Controller) {
	v1Sms := v1.Group("sms")
	{
		v1Sms.Post("", c.Sms().SendSms)
	}
}
