package controllerv1

import (
	servicev1 "wailik.com/internal/authz/service/v1"
)

type Controller interface {
	Domains() *domainController
	Permissions() *permissionController
	Groups() *groupController
	Roles() *roleController
	Denies() *denyController
	Subjects() *subjectController
}

type controller struct {
	srvc servicev1.Service
}

func NewController(s servicev1.Service) Controller {
	var c Controller = &controller{
		srvc: s,
	}

	return c
}

func (c *controller) Domains() *domainController         { return newDomainController(c) }
func (c *controller) Permissions() *permissionController { return newPermissionController(c) }
func (c *controller) Groups() *groupController           { return newGroupController(c) }
func (c *controller) Roles() *roleController             { return newRoleController(c) }
func (c *controller) Denies() *denyController            { return newDenyController(c) }
func (c *controller) Subjects() *subjectController       { return newSubjectController(c) }
