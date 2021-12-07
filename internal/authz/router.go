package server

import (
	"github.com/gofiber/fiber/v2"

	controllerv1 "wailik.com/internal/authz/controller/v1"
	servicev1 "wailik.com/internal/authz/service/v1"
)

func route(app *fiber.App, srvc servicev1.Service) {
	c := controllerv1.NewController(srvc)

	v1 := app.Group("v1")
	{
		groupRoute(v1, c)
		roleRoute(v1, c)
		domainRoute(v1, c)
		denyRoute(v1, c)
		permissionRoute(v1, c)
		subjectRoute(v1, c)
	}
}

func groupRoute(v1 fiber.Router, c controllerv1.Controller) {
	v1Group := v1.Group("groups")
	{
		v1Group.Post("", c.Groups().Create)
		v1Group.Get("", c.Groups().List)
		v1Group.Delete(":groupId", c.Groups().Delete)
		v1Group.Patch(":groupId", c.Groups().Update)
		v1Group.Get(":groupId", c.Groups().Get)
		v1Group.Get(":groupId/permissions", c.Groups().Permissions)
		v1Group.Get(":groupId/roles", c.Groups().Roles)
		v1Group.Post(":groupId/roles", c.Groups().AddRole)
		v1Group.Delete(":groupId/roles/:roleId", c.Groups().DeleteRole)
		v1Group.Get(":groupId/denies", c.Groups().Denies)
		v1Group.Post(":groupId/denies", c.Groups().AddDeny)
		v1Group.Delete(":groupId/denies/:denyId", c.Groups().DeleteDeny)
		v1Group.Get(":groupId/members", c.Groups().Members)
		v1Group.Post(":groupId/members", c.Groups().AddMember)
		v1Group.Delete(":groupId/members/:memberId", c.Groups().DeleteMember)
	}
}

func roleRoute(v1 fiber.Router, c controllerv1.Controller) {
	v1Role := v1.Group("roles")
	{
		v1Role.Post("", c.Roles().Create)
		v1Role.Get("", c.Roles().List)
		v1Role.Delete(":roleId", c.Roles().Delete)
		v1Role.Patch(":roleId", c.Roles().Update)
		v1Role.Get(":roleId", c.Roles().Get)
	}
}

func domainRoute(v1 fiber.Router, c controllerv1.Controller) {
	v1Domain := v1.Group("domains")
	{
		v1Domain.Post("", c.Domains().Create)
		v1Domain.Get("", c.Domains().List)
		v1Domain.Delete(":domainId", c.Domains().Delete)
		v1Domain.Patch(":domainId", c.Domains().Update)
		v1Domain.Get(":domainId", c.Domains().Get)
	}
}

func denyRoute(v1 fiber.Router, c controllerv1.Controller) {
	v1Deny := v1.Group("denies")
	{
		v1Deny.Post("", c.Denies().Create)
		v1Deny.Get("", c.Denies().List)
		v1Deny.Delete(":denyId", c.Denies().Delete)
		v1Deny.Patch(":denyId", c.Denies().Update)
		v1Deny.Get(":denyId", c.Denies().Get)
	}
}

func permissionRoute(v1 fiber.Router, c controllerv1.Controller) {
	v1Permission := v1.Group("permissions")
	{
		v1Permission.Post("", c.Permissions().Create)
		v1Permission.Get("", c.Permissions().List)
		v1Permission.Delete(":permissionId", c.Permissions().Delete)
		v1Permission.Patch(":permissionId", c.Permissions().Update)
		v1Permission.Get(":permissionId", c.Permissions().Get)
	}
}

func subjectRoute(v1 fiber.Router, c controllerv1.Controller) {
	v1Subject := v1.Group("subjects")
	{
		v1Subject.Post("", c.Subjects().Create)
		v1Subject.Get("", c.Subjects().List)
		v1Subject.Delete(":subjectId", c.Subjects().Delete)
		v1Subject.Patch(":subjectId", c.Subjects().Update)
		v1Subject.Get(":subjectId", c.Subjects().Get)
		v1Subject.Get(":subjectId/permissions", c.Subjects().Permissions)
		v1Subject.Get(":subjectId/roles", c.Subjects().Roles)
		v1Subject.Post(":subjectId/roles", c.Subjects().AddRole)
		v1Subject.Get(":subjectId/roles/can_be_granted", c.Subjects().RolesCanBeGrantedTo)
		v1Subject.Get(":subjectId/roles/can_be_accessed", c.Subjects().RolesCanBeAccessedBy)
		v1Subject.Delete(":subjectId/roles/:roleId", c.Subjects().DeleteRole)
		v1Subject.Get(":subjectId/denies", c.Subjects().Denies)
		v1Subject.Post(":subjectId/denies", c.Subjects().AddDeny)
		v1Subject.Get(":subjectId/denies/can_be_granted", c.Subjects().DeniesCanBeGrantedTo)
		v1Subject.Get(":subjectId/denies/can_be_accessed", c.Subjects().DeniesCanBeAccessedBy)
		v1Subject.Delete(":subjectId/denies/:denyId", c.Subjects().DeleteDeny)
		v1Subject.Get(":subjectId/groups", c.Subjects().Groups)
		v1Subject.Post(":subjectId/groups", c.Subjects().AddGroup)
		v1Subject.Delete(":subjectId/groups/:groupId", c.Subjects().DeleteGroup)
	}
}
