package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	servicev1 "wailik.com/internal/authz/service/v1"
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
)

type denyController struct {
	srvc servicev1.Service
}

func newDenyController(c *controller) *denyController {
	return &denyController{srvc: c.srvc}
}

func (d *denyController) Create(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	var r apiv1.Deny
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	if err := d.srvc.Deny().Create(domainId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *denyController) Delete(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	denyId := c.Params(constant.FldDenyId)
	err := d.srvc.Deny().Delete(domainId, denyId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *denyController) Update(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	denyId := c.Params(constant.FldDenyId)
	var r apiv1.Deny
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r.Id = c.Params(constant.FldDomainId)

	if err := d.srvc.Deny().Update(domainId, denyId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *denyController) List(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	r, err := d.srvc.Deny().List(domainId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *denyController) Get(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	denyId := c.Params(constant.FldDenyId)
	deny, err := d.srvc.Deny().Get(domainId, denyId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, deny)
}
