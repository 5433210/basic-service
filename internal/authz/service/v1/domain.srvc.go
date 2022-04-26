package servicev1

import (
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type domainSrvc struct {
	service *service
}

func newDomainSrvc(s *service) *domainSrvc {
	return &domainSrvc{
		service: s,
	}
}

func (s *domainSrvc) Create(d apiv1.Domain) error {
	log.Infof("create domain:%+v", d)
	type null struct{}

	if domain, err := s.Get(d.Id); domain != nil {
		return errors.NewErrorCf(errors.ErrCdDataExist, err, "domain(%v)", d.Id)
	}

	path := BuildPath(constant.FldDomains, d.Id)
	err := s.service.create(path, d)
	if err != nil {
		return err
	}
	// create subpath
	path = BuildPath(constant.FldDomains, d.Id, constant.FldPermissions)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}
	path = BuildPath(constant.FldDomains, d.Id, constant.FldGroups)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}
	path = BuildPath(constant.FldDomains, d.Id, constant.FldSubjects)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}
	path = BuildPath(constant.FldDomains, d.Id, constant.FldRoles)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}
	path = BuildPath(constant.FldDomains, d.Id, constant.FldDenies)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}

	log.Infof("domain created:%+v", d)

	return nil
}

func (s *domainSrvc) Delete(domainId string) error {
	log.Infof("delete domain:%+v", domainId)
	if domain, err := s.Get(domainId); domain == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "domain(%v)", domainId)
	}
	path := BuildPath(constant.FldDomains, domainId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	log.Infof("domain deleted:%v", domainId)

	return nil
}

func (s *domainSrvc) Update(domainId string, d apiv1.Domain) error {
	log.Infof("update domain:%+v", d)
	if domain, err := s.Get(domainId); domain == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "domain(%v)", domainId)
	}

	path := BuildPath(constant.FldDomains, domainId)
	err := s.service.update(path, d)
	if err != nil {
		return err
	}

	log.Infof("domain updated:%v:%+v", domainId, d)

	return nil
}

func (s *domainSrvc) Get(domainId string) (*apiv1.Domain, error) {
	log.Infof("get domain:%+v", domainId)
	input := &DomainInput{DomainId: domainId}
	domain := apiv1.Domain{}
	err := s.service.query(constant.RuleDomain, input, &domain)
	if err != nil {
		return nil, err
	}
	log.Infof("domain got:%+v", domain)

	return &domain, nil
}

func (s *domainSrvc) List() (*[]apiv1.Domain, error) {
	log.Infof("list domain")
	domains := []apiv1.Domain{}
	err := s.service.query(constant.RuleDomains, nil, &domains)
	if err != nil {
		return nil, err
	}

	log.Infof("domains listed:%+v", domains)

	return &domains, nil
}

func (s *domainSrvc) Has(domainId string) bool {
	path := BuildPath(constant.FldDomains, domainId)

	return s.service.exist(path)
}
