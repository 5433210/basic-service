package captcha

import (
	"github.com/gofiber/fiber/v2"

	controllerv1 "wailik.com/internal/captcha/controller/v1"
	servicev1 "wailik.com/internal/captcha/service/v1"
)

func route(app *fiber.App, srvc servicev1.Service) {
	c := controllerv1.NewController(srvc)

	v1 := app.Group("v1")
	{
		captchaGroup(v1, c)
	}
}

func captchaGroup(v1 fiber.Router, c controllerv1.Controller) {
	v1Captcha := v1.Group("captcha")
	{
		v1Captcha.Post("", c.Captcha().VerifyCaptcha)
		v1Captcha.Get("", c.Captcha().GenerateCaptcha)
	}
}
