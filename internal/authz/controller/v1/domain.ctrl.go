package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	servicev1 "wailik.com/internal/authz/service/v1"
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
)

type domainController struct {
	srv servicev1.Service
}

func newDomainController(c *controller) *domainController {
	return &domainController{srv: c.srvc}
}

func (d *domainController) Create(c *fiber.Ctx) error {
	var r apiv1.Domain
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	if err := d.srv.Domain().Create(r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *domainController) Delete(c *fiber.Ctx) error {
	domainId := c.Params(constant.FldDomainId)
	err := d.srv.Domain().Delete(domainId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *domainController) Update(c *fiber.Ctx) error {
	var r apiv1.Domain
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	domainId := c.Params(constant.FldDomainId)
	if err := d.srv.Domain().Update(domainId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *domainController) List(c *fiber.Ctx) error {
	r, err := d.srv.Domain().List()
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *domainController) Get(c *fiber.Ctx) error {
	domainId := c.Params(constant.FldDomainId)
	domain, err := d.srv.Domain().Get(domainId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, domain)
}
