package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	servicev1 "wailik.com/internal/authz/service/v1"
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
)

type roleController struct {
	srvc servicev1.Service
}

func newRoleController(c *controller) *roleController {
	return &roleController{srvc: c.srvc}
}

func (d *roleController) Create(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	var r apiv1.Role
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	if err := d.srvc.Role().Create(domainId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *roleController) Delete(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	roleId := c.Params(constant.FldRoldId)
	err := d.srvc.Role().Delete(domainId, roleId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *roleController) Update(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	roleId := c.Params(constant.FldRoldId)
	var r apiv1.Role
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r.Id = c.Params(constant.FldDomainId)

	if err := d.srvc.Role().Update(domainId, roleId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *roleController) List(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	r, err := d.srvc.Role().List(domainId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *roleController) Get(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	roleId := c.Params(constant.FldRoldId)
	role, err := d.srvc.Role().Get(domainId, roleId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, role)
}
