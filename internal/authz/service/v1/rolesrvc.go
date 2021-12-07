package servicev1

import (
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type roleSrvc struct {
	service *service
}

func newRoleSrvc(s *service) *roleSrvc {
	return &roleSrvc{
		service: s,
	}
}

func (s *roleSrvc) Create(domainId string, o apiv1.Role) error {
	log.Infof("create domain(%+v) role(%+v):%+v", domainId, o.Id, o)
	if s.Has(domainId, o.Id) {
		return errors.NewErrorCf(errors.ErrCdDataExist, nil, "role(%v:%v)", domainId, o.Id)
	}

	if err := s.checkPermissions(domainId, o.Permissions); err != nil {
		return err
	}

	path := BuildPath(constant.FldDomains, domainId, constant.FldRoles, o.Id)
	err := s.service.create(path, o)
	if err != nil {
		return err
	}

	if err := s.addPermissions(domainId, o.Id, o.Permissions); err != nil {
		return err
	}

	log.Infof("domain(%+v) role(%+v) created:%+v", domainId, o.Id, o)

	return nil
}

func (s *roleSrvc) Delete(domainId string, roleId string) error {
	log.Infof("delete domain(%+v) role(%+v)", domainId, roleId)
	role, err := s.Get(domainId, roleId)
	if role == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "role(%v:%v)", domainId, roleId)
	}
	path := BuildPath(constant.FldDomains, domainId, constant.FldRoles, roleId, constant.FldGroups)
	err = s.service.ldb.TraverseWithPrefix([]byte(path), func(key []byte, value interface{}) error {
		groupId := string(value.([]byte))

		return s.service.Group().DeleteRole(domainId, groupId, roleId)
	})
	if err != nil {
		return err
	}

	err = s.service.ldb.DeleteWithPrefix([]byte(path))
	if err != nil {
		return err
	}

	path = BuildPath(constant.FldDomains, domainId, constant.FldRoles, roleId, constant.FldSubjects)
	err = s.service.ldb.TraverseWithPrefix([]byte(path), func(key []byte, value interface{}) error {
		subjectId := string(value.([]byte))

		return s.service.Subject().DeleteRole(domainId, subjectId, roleId)
	})
	if err != nil {
		return err
	}

	err = s.service.ldb.DeleteWithPrefix([]byte(path))
	if err != nil {
		return err
	}

	path = BuildPath(constant.FldDomains, domainId, constant.FldRoles, roleId)

	if err = s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	if err = s.deletePermissions(domainId, roleId, role.Permissions); err != nil {
		return err
	}

	log.Infof("domain(%+v) role(%+v) deleted", domainId, roleId)

	return nil
}

func (s *roleSrvc) Update(domainId string, roleId string, o apiv1.Role) error {
	log.Infof("update domain(%+v) role(%+v):%+v", domainId, o.Id, o)
	role, err := s.Get(domainId, o.Id)
	if role == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "role(%v:%+v)", domainId, o)
	}

	path := BuildPath(constant.FldDomains, domainId, constant.FldRoles, o.Id)
	if err = s.service.update(path, o); err != nil {
		return err
	}

	if err = s.deletePermissions(domainId, roleId, role.Permissions); err != nil {
		return err
	}

	if err := s.addPermissions(domainId, roleId, o.Permissions); err != nil {
		return err
	}
	log.Infof("domain(%+v) role(%+v) updated:%+v", domainId, roleId, o)

	return nil
}

func (s *roleSrvc) Get(domainId string, roleId string) (*apiv1.Role, error) {
	log.Infof("get domain(%+v) role(%+v)", domainId, roleId)
	path := BuildPath(constant.FldDomains, domainId, constant.FldRoles, roleId)
	role := apiv1.Role{}
	err := s.service.get(path, &role)
	if err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) role(%+v) got:%+v", domainId, roleId, role)

	return &role, nil
}

func (s *roleSrvc) List(domainId string) (*[]apiv1.Role, error) {
	log.Infof("list domain(%+v) roles", domainId)
	path := BuildPath(constant.FldDomains, domainId, constant.FldRoles)
	roles := make([]apiv1.Role, 0)

	if err := s.service.list(path, func(id string, j []byte) error {
		role := apiv1.Role{}
		if err := json.Unmarshal(j, &role); err != nil {
			return err
		}
		role.Id = id
		roles = append(roles, role)

		return nil
	}); err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) roles listed:%+v", domainId, roles)

	return &roles, nil
}

func (s *roleSrvc) Has(domainId string, roleId string) bool {
	path := BuildPath(constant.FldDomains, domainId, constant.FldRoles, roleId)

	return s.service.exist(path)
}

func (s *roleSrvc) checkPermissions(domainId string, o apiv1.PermissionsInRole) error {
	inclPerms := o.Included
	for p := range inclPerms {
		if !s.service.Permission().Has(domainId, p) {
			return errors.NewErrorC(errors.ErrCdInvalidPermission, nil)
		}
	}

	exclPerms := o.Excluded
	for _, p := range exclPerms {
		if !s.service.Permission().Has(domainId, p) {
			return errors.NewErrorC(errors.ErrCdInvalidPermission, nil)
		}
	}

	return nil
}

func (s *roleSrvc) addPermissions(domainId string, roleId string, o apiv1.PermissionsInRole) error {
	var path string
	inclPerms := o.Included
	exclPerms := o.Excluded

	// add relationship between permission and role for subsequent "Permission" delete operation
	for p := range inclPerms {
		path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
			constant.FldPermissions, p, constant.FldRoles, constant.FldIncluded, roleId)
		if err := s.service.ldb.WriteObject([]byte(path), []byte(roleId)); err != nil {
			return err
		}
	}
	for _, p := range exclPerms {
		path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldPermissions, p,
			constant.FldRoles, constant.FldExcluded, roleId)
		if err := s.service.ldb.WriteObject([]byte(path), []byte(roleId)); err != nil {
			return err
		}
	}

	return nil
}

func (s *roleSrvc) deletePermissions(domainId string, roleId string, o apiv1.PermissionsInRole) error {
	var path string
	inclPerms := o.Included
	exclPerms := o.Excluded

	// add relationship between permission and role for subsequent "Permission" delete operation
	for p := range inclPerms {
		path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldPermissions, p,
			constant.FldRoles, constant.FldIncluded, roleId)
		if err := s.service.ldb.Delete([]byte(path)); err != nil {
			return err
		}
	}
	for _, p := range exclPerms {
		path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldPermissions, p,
			constant.FldRoles, constant.FldExcluded, roleId)
		if err := s.service.ldb.Delete([]byte(path)); err != nil {
			return err
		}
	}

	return nil
}
