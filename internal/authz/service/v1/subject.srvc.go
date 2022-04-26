package servicev1

import (
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type subjectSrvc struct {
	service *service
}

func newSubjectSrvc(s *service) *subjectSrvc {
	return &subjectSrvc{
		service: s,
	}
}

func (s *subjectSrvc) Create(domainId string, o apiv1.Subject) error {
	log.Infof("create domain(%+v) subject(%+v):%+v", domainId, o.Id, o)

	type null struct{}
	if subject, err := s.Get(domainId, o.Id); subject != nil {
		return errors.NewErrorCf(errors.ErrCdDataExist, err, "subject(%v:%v)", domainId, o.Id)
	}
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, o.Id)
	err := s.service.create(path, o)
	if err != nil {
		return err
	}
	path = BuildPath(constant.FldDomains, domainId, constant.FldSubjects, o.Id, constant.FldDenies)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}
	path = BuildPath(constant.FldDomains, domainId, constant.FldSubjects, o.Id, constant.FldRoles)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}
	path = BuildPath(constant.FldDomains, domainId, constant.FldSubjects, o.Id, constant.FldGroups)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}

	log.Infof("domain(%+v) group created:%+v", domainId, o)

	return nil
}

func (s *subjectSrvc) Delete(domainId string, subjectId string) error {
	log.Infof("delete domain(%+v) subject(%+v)", domainId, subjectId)
	if subject, err := s.Get(domainId, subjectId); subject == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "subject(%v:%v)", domainId, subjectId)
	}
	roles, err := s.Roles(domainId, subjectId)
	if err != nil {
		return err
	}
	for _, role := range *roles {
		err = s.DeleteRole(domainId, subjectId, role.Id)
		if err != nil {
			return err
		}
	}

	denies, err := s.Denies(domainId, subjectId)
	if err != nil {
		return err
	}
	for _, deny := range *denies {
		err = s.DeleteDeny(domainId, subjectId, deny.Id)
		if err != nil {
			return err
		}
	}

	rs, err := s.Groups(domainId, subjectId)
	if err != nil {
		return err
	}

	for _, group := range *rs {
		err = s.service.Group().DeleteMember(domainId, group.Id, subjectId)
		if err != nil {
			return err
		}
	}

	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId)
	if err = s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	log.Infof("domain(%+v) group deleted:%+v", domainId, subjectId)

	return nil
}

func (s *subjectSrvc) Update(domainId string, subjectId string, o apiv1.Subject) error {
	log.Infof("update domain(%+v) subject(%+v):%+v", domainId, subjectId, o)

	if subject, err := s.Get(domainId, o.Id); subject == nil {
		return errors.NewErrorCf(errors.ErrCdDataNotExist, err, "subject(%v:%+v)", domainId, o)
	}

	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, o.Id)
	err := s.service.update(path, o)
	if err != nil {
		return err
	}

	log.Infof("domain(%+v) group updated:%+v", domainId, subjectId)

	return nil
}

func (s *subjectSrvc) Get(domainId string, subjectId string) (*apiv1.Subject, error) {
	log.Infof("get domain(%+v) subject(%+v)", domainId, subjectId)
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId)
	subject := apiv1.Subject{}
	err := s.service.get(path, &subject)
	if err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) subject(%+v) got:%+v", domainId, subject, subject)

	return &subject, nil
}

func (s *subjectSrvc) List(domainId string) (*[]apiv1.Subject, error) {
	log.Infof("list domain(%+v) subjects", domainId)
	input := &DomainInput{
		DomainId: domainId,
	}
	subjects := []apiv1.Subject{}
	err := s.service.query(constant.RuleSubjects, input, &subjects)
	if err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) subjects listed:%+v", domainId, subjects)

	return &subjects, nil
}

func (s *subjectSrvc) Permissions(domainId string, subjectId string) (*[]apiv1.Permission, error) {
	log.Infof("list domain(%+v) subject(%+v) permissions", domainId, subjectId)
	input := &DomainSubjectInput{
		DomainId:  domainId,
		SubjectId: subjectId,
	}
	permissions := []apiv1.Permission{}
	err := s.service.query(constant.RuleSubjectPermissions, input, &permissions)
	if err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) subject(%+v) permissions listed:%+v", domainId, subjectId, permissions)

	return &permissions, nil
}

func (s *subjectSrvc) Roles(domainId string, subjectId string) (*[]apiv1.RoleOptions, error) {
	log.Infof("list domain(%+v) subject(%+v) roles", domainId, subjectId)
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldRoles)
	roles := make([]apiv1.RoleOptions, 0)
	if err := s.service.list(path, func(id string, j []byte) error {
		role := apiv1.RoleOptions{}
		if err := json.Unmarshal(j, &role); err != nil {
			return err
		}
		role.Id = id
		roles = append(roles, role)

		return nil
	}); err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) subject(%+v) roles listed:%+v", domainId, subjectId, roles)

	return &roles, nil
}

func (s *subjectSrvc) Denies(domainId string, subjectId string) (*[]apiv1.DenyOptions, error) {
	log.Infof("list domain(%+v) subject(%+v) denies", domainId, subjectId)
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldDenies)
	denies := make([]apiv1.DenyOptions, 0)
	if err := s.service.list(path, func(id string, j []byte) error {
		deny := apiv1.DenyOptions{}
		if err := json.Unmarshal(j, &deny); err != nil {
			return err
		}
		deny.Id = id
		denies = append(denies, deny)

		return nil
	}); err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) subject(%+v) denies listed:%+v", domainId, subjectId, denies)

	return &denies, nil
}

func (s *subjectSrvc) Groups(domainId string, subjectId string) (*[]apiv1.GroupOptions, error) {
	log.Infof("list domain(%+v) subject(%+v) groups", domainId, subjectId)
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldGroups)
	groups := make([]apiv1.GroupOptions, 0)
	if err := s.service.list(path, func(id string, j []byte) error {
		group := apiv1.GroupOptions{}
		if err := json.Unmarshal(j, &group); err != nil {
			return err
		}
		group.Id = id
		groups = append(groups, group)

		return nil
	}); err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) subject(%+v) subjects listed:%+v", domainId, subjectId, groups)

	return &groups, nil
}

func (s *subjectSrvc) DeleteRole(domainId string, subjectId string, roleId string) error {
	log.Infof("delete domain(%+v) subject(%+v) role(%+v)", domainId, subjectId, roleId)
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldRoles, roleId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}
	if err := s.deleteRole(domainId, subjectId, roleId); err != nil {
		return err
	}
	log.Infof("domain(%+v) subject(%+v) role deleted:%+v", domainId, subjectId, roleId)

	return nil
}

func (s *subjectSrvc) DeleteDeny(domainId string, subjectId string, denyId string) error {
	log.Infof("delete domain(%+v) subject(%+v) deny(%+v)", domainId, subjectId, denyId)

	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldDenies, denyId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}
	if err := s.deleteDeny(domainId, subjectId, denyId); err != nil {
		return err
	}

	log.Infof("domain(%+v) subject(%+v) deny deleted:%+v", domainId, subjectId, denyId)

	return nil
}

func (s *subjectSrvc) DeleteGroup(domainId string, subjectId string, groupId string) error {
	log.Infof("delete domain(%+v) subject(%+v) group(%+v)", domainId, subjectId, groupId)

	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldGroups, groupId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	path = BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId, constant.FldMembers, subjectId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	log.Infof("domain(%+v) subject(%+v) member deleted:%+v", domainId, groupId, subjectId)

	return nil
}

func (s *subjectSrvc) AddRole(domainId string, subjectId string, o apiv1.RoleOptions) error {
	log.Infof("add domain(%+v) subject(%+v) role(%+v)", domainId, subjectId, o)
	if err := s.checkRole(domainId, o.Id); err != nil {
		return err
	}
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldRoles, o.Id)
	if err := s.service.create(path, o); err != nil {
		return err
	}

	if err := s.addRole(domainId, subjectId, o.Id); err != nil {
		return err
	}

	log.Infof("domain(%+v) subject(%+v) role added:%+v", domainId, subjectId, o)

	return nil
}

func (s *subjectSrvc) AddDeny(domainId string, subjectId string, o apiv1.DenyOptions) error {
	log.Infof("add domain(%+v) subject(%+v) deny(%+v)", domainId, subjectId, o)
	if err := s.checkDeny(domainId, o.Id); err != nil {
		return err
	}
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldDenies, o.Id)
	if err := s.service.create(path, o); err != nil {
		return err
	}
	if err := s.addDeny(domainId, subjectId, o.Id); err != nil {
		return err
	}
	log.Infof("domain(%+v) subject(%+v) deny added:%+v", domainId, subjectId, o)

	return nil
}

func (s *subjectSrvc) AddGroup(domainId string, subjectId string, o apiv1.GroupOptions) error {
	log.Infof("add domain(%+v) subject(%+v) group(%+v)", domainId, subjectId, o)
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId, constant.FldGroups, o.Id)
	if err := s.service.create(path, o); err != nil {
		return err
	}

	path = BuildPath(constant.FldDomains, domainId, constant.FldGroups, o.Id, constant.FldMembers, subjectId)
	if err := s.service.create(path, o); err != nil {
		return err
	}

	log.Infof("domain(%+v) subject(%+v) member added:%+v", domainId, subjectId, o)

	return nil
}

func (s *subjectSrvc) RolesCanBeGrantedTo(domainId string, subjectId string) (*[]apiv1.RoleOptions, error) {
	log.Infof("list roles can be granted to domain(%+v) subject(%+v)", domainId, subjectId)
	input := &DomainSubjectInput{
		DomainId:  domainId,
		SubjectId: subjectId,
	}
	roles := []apiv1.RoleOptions{}
	err := s.service.query(constant.RuleRolesCanBeGranted, input, &roles)
	if err != nil {
		return nil, err
	}

	log.Debugf("domain(%+v) roles can be granted to subject(%v) listed:%+v", domainId, subjectId, roles)

	return &roles, nil
}

func (s *subjectSrvc) RolesCanBeAccessedBy(domainId string, subjectId string) (*[]apiv1.RoleOptions, error) {
	log.Infof("list roles can be accessed by domain(%+v) subject(%+v)", domainId, subjectId)
	input := &DomainSubjectInput{
		DomainId:  domainId,
		SubjectId: subjectId,
	}
	roles := []apiv1.RoleOptions{}
	err := s.service.query(constant.RuleRolesCanBeAccessed, input, &roles)
	if err != nil {
		return nil, err
	}

	log.Debugf("domain(%+v) roles can be accessed by subject(%v) listed:%+v", domainId, subjectId, roles)

	return &roles, nil
}

func (s *subjectSrvc) DeniesCanBeGrantedTo(domainId string, subjectId string) (*[]apiv1.DenyOptions, error) {
	log.Infof("list denies can be granted to domain(%+v) subject(%+v)", domainId, subjectId)
	input := &DomainSubjectInput{
		DomainId:  domainId,
		SubjectId: subjectId,
	}
	denies := []apiv1.DenyOptions{}
	err := s.service.query(constant.RuleDeniesCanBeGranted, input, &denies)
	if err != nil {
		return nil, err
	}

	log.Debugf("domain(%+v) denies can be granted to subject(%v) listed:%+v", domainId, subjectId, denies)

	return &denies, nil
}

func (s *subjectSrvc) DeniesCanBeAccessedBy(domainId string, subjectId string) (*[]apiv1.DenyOptions, error) {
	log.Infof("list denies can be accessed by domain(%+v) subject(%+v)", domainId, subjectId)
	input := &DomainSubjectInput{
		DomainId:  domainId,
		SubjectId: subjectId,
	}
	denies := []apiv1.DenyOptions{}
	err := s.service.query(constant.RuleDeniesCanBeAccessed, input, &denies)
	if err != nil {
		return nil, err
	}

	log.Debugf("domain(%+v) denies can be accessed by subject(%v) listed:%+v", domainId, subjectId, denies)

	return &denies, nil
}

func (s *subjectSrvc) Has(domainId string, subjectId string) bool {
	path := BuildPath(constant.FldDomains, domainId, constant.FldSubjects, subjectId)

	return s.service.exist(path)
}

func (s *subjectSrvc) checkRole(domainId string, roleId string) error {
	if !s.service.Role().Has(domainId, roleId) {
		return errors.NewErrorC(errors.ErrCdInvalidRole, nil)
	}

	return nil
}

func (s *subjectSrvc) addRole(domainId string, subjectId string, roleId string) error {
	path := BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldRoles, roleId,
		constant.FldSubjects, subjectId)
	if err := s.service.ldb.WriteObject([]byte(path), []byte(subjectId)); err != nil {
		return err
	}

	return nil
}

func (s *subjectSrvc) deleteRole(domainId string, subjectId string, roleId string) error {
	path := BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldRoles, roleId,
		constant.FldSubjects, subjectId)
	if err := s.service.ldb.Delete([]byte(path)); err != nil {
		return err
	}

	return nil
}

func (s *subjectSrvc) checkDeny(domainId string, denyId string) error {
	if !s.service.Deny().Has(domainId, denyId) {
		return errors.NewErrorC(errors.ErrCdInvalidDeny, nil)
	}

	return nil
}

func (s *subjectSrvc) addDeny(domainId string, subjectId string, denyId string) error {
	path := BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldDenies, denyId,
		constant.FldSubjects, subjectId)
	if err := s.service.ldb.WriteObject([]byte(path), []byte(subjectId)); err != nil {
		return err
	}

	return nil
}

func (s *subjectSrvc) deleteDeny(domainId string, subjectId string, denyId string) error {
	path := BuildPath(constant.FldIndexes, constant.FldDomains, domainId, constant.FldDenies, denyId,
		constant.FldSubjects, subjectId)
	if err := s.service.ldb.Delete([]byte(path)); err != nil {
		return err
	}

	return nil
}
