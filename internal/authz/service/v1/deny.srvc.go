package servicev1

import (
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type denySrvc struct {
	service *service
}

func newDenySrvc(s *service) *denySrvc {
	return &denySrvc{
		service: s,
	}
}

func (s *denySrvc) Create(domainId string, o apiv1.Deny) error {
	log.Infof("create domain(%+v) deny(%+v):%+v", domainId, o.Id, o)

	if s.Has(domainId, o.Id) {
		return errors.NewErrorCf(errors.ErrCdDataExist, nil, "role(%v:%v)", domainId, o.Id)
	}

	if err := s.checkPermissions(domainId, o.Permissions); err != nil {
		return err
	}
	path := BuildPath(constant.FldDomains, domainId, constant.FldDenies, o.Id)
	err := s.service.create(path, o)
	if err != nil {
		return err
	}
	if err := s.addPermissions(domainId, o.Id, o.Permissions); err != nil {
		return err
	}
	log.Infof("domain(%+v) deny(%+v) created:%+v", domainId, o.Id, o)

	return nil
}

func (s *denySrvc) Delete(domainId string, denyId string) error {
	log.Infof("delete domain(%+v) deny(%+v)", domainId, denyId)

	deny, err := s.Get(domainId, denyId)
	if deny == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "role(%v:%v)", domainId, denyId)
	}
	path := BuildPath(constant.FldDomains, domainId, constant.FldDenies, denyId, constant.FldGroups)
	err = s.service.ldb.TraverseWithPrefix([]byte(path), func(key []byte, value interface{}) error {
		groupId := string(value.([]byte))

		return s.service.Group().DeleteRole(domainId, groupId, denyId)
	})
	if err != nil {
		return err
	}
	path = BuildPath(constant.FldDomains, domainId, constant.FldDenies, denyId, constant.FldSubjects)
	err = s.service.ldb.TraverseWithPrefix([]byte(path), func(key []byte, value interface{}) error {
		subjectId := string(value.([]byte))

		return s.service.Subject().DeleteRole(domainId, subjectId, denyId)
	})
	if err != nil {
		return err
	}

	path = BuildPath(constant.FldDomains, domainId, constant.FldDenies, denyId)

	if err = s.service.deleteWithPrefix(path); err != nil {
		return err
	}
	if err = s.deletePermissions(domainId, denyId, deny.Permissions); err != nil {
		return err
	}

	log.Infof("domain(%+v) deny(%+v) deleted", domainId, denyId)

	return nil
}

func (s *denySrvc) Update(domainId string, denyId string, o apiv1.Deny) error {
	log.Infof("update domain(%+v) deny(%+v):%+v", domainId, denyId, o)

	deny, err := s.Get(domainId, o.Id)
	if deny == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "role(%v:%+v)", domainId, o)
	}

	path := BuildPath(constant.FldDomains, domainId, constant.FldDenies, o.Id)
	err = s.service.update(path, o)
	if err != nil {
		return err
	}
	if err = s.deletePermissions(domainId, denyId, deny.Permissions); err != nil {
		return err
	}
	if err := s.addPermissions(domainId, denyId, deny.Permissions); err != nil {
		return err
	}
	log.Infof("domain(%+v) deny(%+v) updated:%+v", domainId, denyId, o)

	return nil
}

func (s *denySrvc) Get(domainId string, denyId string) (*apiv1.Deny, error) {
	log.Infof("get domain(%+v) deny(%+v)", domainId, denyId)

	path := BuildPath(constant.FldDomains, domainId, constant.FldDenies, denyId)
	deny := apiv1.Deny{}
	err := s.service.get(path, &deny)
	if err != nil {
		return nil, err
	}

	deny.Id = denyId
	log.Infof("domain(%+v) deny(%+v) got:%+v", domainId, denyId, deny)

	return &deny, nil
}

func (s *denySrvc) List(domainId string) (*[]apiv1.Deny, error) {
	log.Infof("list domain(%+v) denies", domainId)

	path := BuildPath(constant.FldDomains, domainId, constant.FldDenies)
	denies := make([]apiv1.Deny, 0)
	if err := s.service.list(path, func(id string, j []byte) error {
		deny := apiv1.Deny{}
		if err := json.Unmarshal(j, &deny); err != nil {
			return err
		}
		deny.Id = id
		denies = append(denies, deny)

		return nil
	}); err != nil {
		return nil, err
	}
	log.Infof("domain(%+v) denies listed:%+v", domainId, denies)

	return &denies, nil
}

func (s *denySrvc) Has(domainId string, denyId string) bool {
	path := BuildPath(constant.FldDomains, domainId, constant.FldDenies, denyId)

	return s.service.exist(path)
}

func (s *denySrvc) checkPermissions(domainId string, pid apiv1.PermissionsInDeny) error {
	inclPerms := pid.Included
	for _, p := range inclPerms {
		if !s.service.Permission().Has(domainId, p) {
			return errors.NewErrorC(errors.ErrCdInvalidPermission, nil)
		}
	}

	exclPerms := pid.Excluded
	for _, p := range exclPerms {
		if !s.service.Permission().Has(domainId, p) {
			return errors.NewErrorC(errors.ErrCdInvalidPermission, nil)
		}
	}

	return nil
}

func (s *denySrvc) addPermissions(domainId string, denyId string, pid apiv1.PermissionsInDeny) error {
	var path string
	inclPerms := pid.Included
	exclPerms := pid.Excluded

	// add relationship between permission and role for subsequent "Permission" delete operation
	for _, p := range inclPerms {
		path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldPermissions, p, constant.FldDenies, constant.FldIncluded, denyId)
		if err := s.service.ldb.WriteObject([]byte(path), []byte(denyId)); err != nil {
			return err
		}
	}
	for _, p := range exclPerms {
		path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldPermissions, p, constant.FldDenies, constant.FldExcluded, denyId)
		if err := s.service.ldb.WriteObject([]byte(path), []byte(denyId)); err != nil {
			return err
		}
	}

	return nil
}

func (s *denySrvc) deletePermissions(domainId string, denyId string, pid apiv1.PermissionsInDeny) error {
	var path string
	inclPerms := pid.Included
	exclPerms := pid.Excluded

	// add relationship between permission and role for subsequent "Permission" delete operation
	for _, p := range inclPerms {
		path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldPermissions, p, constant.FldDenies, constant.FldIncluded, denyId)
		if err := s.service.ldb.Delete([]byte(path)); err != nil {
			return err
		}
	}
	for _, p := range exclPerms {
		path = BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldPermissions, p, constant.FldDenies, constant.FldExcluded, denyId)
		if err := s.service.ldb.Delete([]byte(path)); err != nil {
			return err
		}
	}

	return nil
}
