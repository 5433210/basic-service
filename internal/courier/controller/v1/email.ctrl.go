package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	apiv1 "wailik.com/internal/courier/api/v1"
	servicev1 "wailik.com/internal/courier/service/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type emailController struct {
	srvc servicev1.Service
}

func newEmailController(c *controller) *emailController {
	return &emailController{srvc: c.srvc}
}

func (i *emailController) SendEmail(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.Email
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r, err := i.srvc.Email().Send(json)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldStatus] = r

	return core.WriteResponse(c, nil, data)
}
