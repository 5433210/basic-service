package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	servicev1 "wailik.com/internal/authz/service/v1"
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
)

type subjectController struct {
	srvc servicev1.Service
}

func newSubjectController(c *controller) *subjectController {
	return &subjectController{srvc: c.srvc}
}

func (d *subjectController) Create(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	var r apiv1.Subject
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	if err := d.srvc.Subject().Create(domainId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *subjectController) Delete(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	err := d.srvc.Subject().Delete(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *subjectController) Update(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	var r apiv1.Subject
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r.Id = c.Params(constant.FldDomainId)

	if err := d.srvc.Subject().Update(domainId, subjectId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *subjectController) List(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	r, err := d.srvc.Subject().List(domainId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, r)
}

func (d *subjectController) Get(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	subject, err := d.srvc.Subject().Get(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, subject)
}

func (d *subjectController) Permissions(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	permissions, err := d.srvc.Subject().Permissions(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, permissions)
}

func (d *subjectController) Roles(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	roles, err := d.srvc.Subject().Roles(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, roles)
}

func (d *subjectController) AddRole(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	var r apiv1.RoleOptions
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	if !d.srvc.Role().Has(domainId, r.Id) {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataNotExist, nil), nil)
	}

	if err := d.srvc.Subject().AddRole(domainId, subjectId, r); err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *subjectController) DeleteRole(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	roleId := c.Params(constant.FldRoldId)
	err := d.srvc.Subject().DeleteRole(domainId, subjectId, roleId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *subjectController) Denies(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	denies, err := d.srvc.Subject().Get(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, denies)
}

func (d *subjectController) AddDeny(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	var r apiv1.DenyOptions
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	if !d.srvc.Deny().Has(domainId, r.Id) {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataNotExist, nil), nil)
	}

	err := d.srvc.Subject().AddDeny(domainId, subjectId, r)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *subjectController) DeleteDeny(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	denyId := c.Params(constant.FldDenyId)
	err := d.srvc.Subject().DeleteDeny(domainId, subjectId, denyId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *subjectController) Groups(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	members, err := d.srvc.Subject().Groups(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, members)
}

func (d *subjectController) AddGroup(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	var r apiv1.GroupOptions
	if err := c.BodyParser(&r); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	if !d.srvc.Group().Has(domainId, r.Id) {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataNotExist, nil), nil)
	}
	err := d.srvc.Subject().AddGroup(domainId, subjectId, r)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *subjectController) DeleteGroup(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	groupId := c.Params(constant.FldGroupId)
	err := d.srvc.Subject().DeleteGroup(domainId, subjectId, groupId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

func (d *subjectController) RolesCanBeGrantedTo(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	role, err := d.srvc.Subject().RolesCanBeGrantedTo(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, role)
}

func (d *subjectController) RolesCanBeAccessedBy(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	role, err := d.srvc.Subject().RolesCanBeAccessedBy(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, role)
}

func (d *subjectController) DeniesCanBeGrantedTo(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	deny, err := d.srvc.Subject().DeniesCanBeGrantedTo(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, deny)
}

func (d *subjectController) DeniesCanBeAccessedBy(c *fiber.Ctx) error {
	domainId := c.Query(constant.FldDomainId, constant.DefaultDomainId)
	subjectId := c.Params(constant.FldSubjectId)
	deny, err := d.srvc.Subject().DeniesCanBeAccessedBy(domainId, subjectId)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, deny)
}
