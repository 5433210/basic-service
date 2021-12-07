package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	apiv1 "wailik.com/internal/captcha/api/v1"
	servicev1 "wailik.com/internal/captcha/service/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type captchaController struct {
	srvc servicev1.Service
}

func newCaptchaController(c *controller) *captchaController {
	return &captchaController{srvc: c.srvc}
}

func (i *captchaController) GenerateCaptcha(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	challengeMode := c.Query(constant.FldChallengeMode, "")
	r, err := i.srvc.Captcha().GenerateCaptcha(apiv1.ChallengeMode(challengeMode))
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldChallenge] = r

	return core.WriteResponse(c, nil, data)
}

func (i *captchaController) VerifyCaptcha(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.Captcha
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r := i.srvc.Captcha().VerifyCaptcha(json)
	data := make(map[string]interface{})
	data[constant.FldOk] = r

	return core.WriteResponse(c, nil, data)
}
