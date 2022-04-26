package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	apiv1 "wailik.com/internal/courier/api/v1"
	servicev1 "wailik.com/internal/courier/service/v1"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type smsController struct {
	srvc servicev1.Service
}

func newSmsController(c *controller) *smsController {
	return &smsController{srvc: c.srvc}
}

func (i *smsController) SendSms(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.Sms
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	if err := i.srvc.Sms().Send(json); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}
