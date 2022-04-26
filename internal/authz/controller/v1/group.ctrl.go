package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	servicev1 "wailik.com/internal/authz/service/v1"
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
)

type groupController struct {
	srvc servicev1.Service
}

func newGroupController(c *controller) *groupController {
	return &groupController{srvc: c.srvc}
}

func (d *groupController) Create(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	var r apiv1.Group
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	if err := d.srvc.Group().Create(domainId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *groupController) Delete(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	err := d.srvc.Group().Delete(domainId, groupId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *groupController) Update(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	var r apiv1.Group
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r.Id = c.Params(constant.FldDomainId)

	if err := d.srvc.Group().Update(domainId, groupId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *groupController) List(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	r, err := d.srvc.Group().List(domainId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *groupController) Get(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	group, err := d.srvc.Group().Get(domainId, groupId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, group)
}

func (d *groupController) Permissions(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	permissions, err := d.srvc.Group().Permissions(domainId, groupId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, permissions)
}

func (d *groupController) Roles(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	roles, err := d.srvc.Group().Roles(domainId, groupId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, roles)
}

func (d *groupController) AddRole(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	var r apiv1.RoleOptions
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	if !d.srvc.Role().Has(domainId, r.Id) {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataNotExist, nil), nil)
	}
	err := d.srvc.Group().AddRole(domainId, groupId, r)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *groupController) DeleteRole(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	roleId := c.Params(constant.FldRoldId)
	err := d.srvc.Group().DeleteRole(domainId, groupId, roleId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *groupController) Denies(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	denies, err := d.srvc.Group().Get(domainId, groupId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, denies)
}

func (d *groupController) AddDeny(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	var r apiv1.DenyOptions
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	if !d.srvc.Deny().Has(domainId, r.Id) {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataNotExist, nil), nil)
	}
	err := d.srvc.Group().AddDeny(domainId, groupId, r)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *groupController) DeleteDeny(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	denyId := c.Params(constant.FldDenyId)
	err := d.srvc.Group().DeleteDeny(domainId, groupId, denyId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *groupController) Members(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	members, err := d.srvc.Group().Members(domainId, groupId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, members)
}

func (d *groupController) AddMember(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	var r apiv1.SubjectOptions
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	if !d.srvc.Subject().Has(domainId, r.Id) {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataNotExist, nil), nil)
	}
	err := d.srvc.Group().AddMember(domainId, groupId, r)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *groupController) DeleteMember(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	groupId := c.Params(constant.FldGroupId)
	subjectId := c.Params(constant.FldSubjectId)
	err := d.srvc.Group().DeleteMember(domainId, groupId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}
