package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	servicev1 "wailik.com/internal/authz/service/v1"
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
)

type permissionController struct {
	srvc servicev1.Service
}

func newPermissionController(c *controller) *permissionController {
	return &permissionController{srvc: c.srvc}
}

func (d *permissionController) Create(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, "78839721-a274-4a01-a2be-2725903bcf82")
	var r apiv1.Permission
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	if err := d.srvc.Permission().Create(domainId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *permissionController) Delete(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, "78839721-a274-4a01-a2be-2725903bcf82")
	permissionId := c.Params(constant.FldPermissionId)
	err := d.srvc.Permission().Delete(domainId, permissionId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *permissionController) Update(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, "78839721-a274-4a01-a2be-2725903bcf82")
	permissionId := c.Params(constant.FldPermissionId)
	var r apiv1.Permission
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r.Id = c.Params(constant.FldDomainId)

	if err := d.srvc.Permission().Update(domainId, permissionId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *permissionController) List(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, "78839721-a274-4a01-a2be-2725903bcf82")
	r, err := d.srvc.Permission().List(domainId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *permissionController) Get(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, "78839721-a274-4a01-a2be-2725903bcf82")
	permissionId := c.Params(constant.FldPermissionId)
	permission, err := d.srvc.Permission().Get(domainId, permissionId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, permission)
}
