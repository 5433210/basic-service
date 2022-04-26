package servicev1

import (
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type permissionSrvc struct {
	service *service
}

func newPermissionSrvc(s *service) *permissionSrvc {
	return &permissionSrvc{
		service: s,
	}
}

func (s *permissionSrvc) Create(domainId string, o apiv1.Permission) error {
	log.Infof("create domain(%+v) permission(%+v):%+v", domainId, o.Id, o)
	if permission, err := s.Get(domainId, o.Id); permission != nil {
		return errors.NewErrorCf(errors.ErrCdDataExist, err, "permission(%v:%v)", domainId, o.Id)
	}
	path := BuildPath(constant.FldDomains, domainId, constant.FldPermissions, o.Id)
	err := s.service.create(path, o)
	if err != nil {
		return err
	}

	log.Infof("domain(%+v) permission(%+v) created:%+v", domainId, o.Id, o)

	return nil
}

func (s *permissionSrvc) Delete(domainId string, permissionId string) error {
	log.Infof("delete domain(%+v) permission(%+v)", domainId, permissionId)
	if permission, err := s.Get(domainId, permissionId); permission == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "permission(%v:%v)", domainId, permissionId)
	}
	path := BuildPath(constant.FldDomains, domainId, constant.FldPermissions, permissionId)

	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	// todo: Delete permission refereced by roles and denies
	// exluded roles
	path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldPermissions, permissionId, constant.FldRoles, constant.FldExcluded)
	err := s.service.ldb.TraverseWithPrefix([]byte(path), func(key []byte, value interface{}) error {
		rp := BuildPath(constant.FldDomains, domainId, constant.FldRoles,
			constant.FldExcluded, string(value.([]byte)))

		return s.service.odb.DeleteOne(rp)
	})
	if err != nil {
		return err
	}
	// included roles
	path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldPermissions, permissionId, constant.FldRoles, constant.FldIncluded)
	err = s.service.ldb.TraverseWithPrefix([]byte(path), func(key []byte, value interface{}) error {
		rp := BuildPath(constant.FldDomains, domainId, constant.FldRoles,
			constant.FldIncluded, string(value.([]byte)))

		return s.service.odb.DeleteOne(rp)
	})
	if err != nil {
		return err
	}

	// excluded denies
	path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldPermissions, permissionId, constant.FldDenies, constant.FldExcluded)
	err = s.service.ldb.TraverseWithPrefix([]byte(path), func(key []byte, value interface{}) error {
		rp := BuildPath(constant.FldDomains, domainId, constant.FldDenies,
			constant.FldExcluded, string(value.([]byte)))

		return s.service.odb.DeleteOne(rp)
	})
	if err != nil {
		return err
	}

	// included denies
	path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldPermissions, permissionId, constant.FldDenies, constant.FldIncluded)
	err = s.service.ldb.TraverseWithPrefix([]byte(path), func(key []byte, value interface{}) error {
		rp := BuildPath(constant.FldDomains, domainId, constant.FldDenies,
			constant.FldIncluded, string(value.([]byte)))

		return s.service.odb.DeleteOne(rp)
	})
	if err != nil {
		return err
	}

	// releationship
	path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldPermissions, permissionId)
	err = s.service.ldb.DeleteWithPrefix([]byte(path))
	if err != nil {
		return err
	}

	log.Infof("domain(%+v) permission(%+v) deleted", domainId, permissionId)

	return nil
}

func (s *permissionSrvc) Update(domainId string, permissionId string, o apiv1.Permission) error {
	log.Infof("update domain(%+v) permission(%+v):%+v", domainId, permissionId, o)
	if permission, err := s.Get(domainId, o.Id); permission == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "permission(%v:%+v)", domainId, o)
	}

	path := BuildPath(constant.FldDomains, domainId, constant.FldPermissions, o.Id)
	err := s.service.update(path, o)
	if err != nil {
		return err
	}

	log.Infof("domain(%+v) permission(%+v) updated:%+v", domainId, permissionId, o)

	return nil
}

func (s *permissionSrvc) Get(domainId string, permissionId string) (*apiv1.Permission, error) {
	log.Infof("get domain(%+v) permission(%+v)", domainId, permissionId)
	path := BuildPath(constant.FldDomains, domainId, constant.FldPermissions, permissionId)
	permission := apiv1.Permission{}
	err := s.service.get(path, &permission)
	if err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) permission(%+v) go:%+v", domainId, permissionId, permission)

	return &permission, nil
}

func (s *permissionSrvc) List(domainId string) (*[]apiv1.Permission, error) {
	path := BuildPath(constant.FldDomains, domainId, constant.FldPermissions)
	permissions := make([]apiv1.Permission, 0)
	if err := s.service.list(path, func(id string, j []byte) error {
		permission := apiv1.Permission{}
		if err := json.Unmarshal(j, &permission); err != nil {
			return err
		}
		permissions = append(permissions, permission)

		return nil
	}); err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) permission listed:%+v", domainId, permissions)

	return &permissions, nil
}

func (s *permissionSrvc) Has(domainId string, permissionId string) bool {
	path := BuildPath(constant.FldDomains, domainId, constant.FldPermissions, permissionId)

	return s.service.exist(path)
}
