package servicev1

import (
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type groupSrvc struct {
	service *service
}

func newGroupSrvc(s *service) *groupSrvc {
	return &groupSrvc{
		service: s,
	}
}

func (s *groupSrvc) Create(domainId string, o apiv1.Group) error {
	log.Infof("create domain(%+v) group(%+v):%+v", domainId, o.Id, o)
	type null struct{}

	// check existing
	if s.Has(domainId, o.Id) {
		return errors.NewErrorCf(errors.ErrCdDataExist, nil, "group(%v:%v)", domainId, o.Id)
	}

	// create group
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, o.Id)
	err := s.service.create(path, o)
	if err != nil {
		return err
	}

	// create null assigned members
	path = BuildPath(constant.FldDomains, domainId, constant.FldGroups, o.Id,
		constant.FldMembers)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}

	// create null assigned roles
	path = BuildPath(constant.FldDomains, domainId, constant.FldGroups, o.Id,
		constant.FldRoles)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}

	// create null assigned denies
	path = BuildPath(constant.FldDomains, domainId, constant.FldGroups, o.Id,
		constant.FldDenies)
	err = s.service.create(path, null{})
	if err != nil {
		return err
	}

	log.Infof("domain(%+v) group created:%+v", domainId, o)

	return nil
}

func (s *groupSrvc) Delete(domainId string, groupId string) error {
	log.Infof("update domain(%+v) role(%+v)", domainId, groupId)
	// check existing
	if !s.Has(domainId, groupId) {
		return errors.NewErrorCf(errors.ErrCdDataExist, nil, "group(%v:%v)", domainId, groupId)
	}

	// delete assigned roles
	roles, err := s.Roles(domainId, groupId)
	if err != nil {
		return err
	}
	for _, role := range *roles {
		err = s.DeleteRole(domainId, groupId, role.Id)
		if err != nil {
			return err
		}
	}

	// delete assigned denies
	denies, err := s.Denies(domainId, groupId)
	if err != nil {
		return err
	}
	for _, deny := range *denies {
		err = s.DeleteDeny(domainId, groupId, deny.Id)
		if err != nil {
			return err
		}
	}

	// delete assigned members
	rs, err := s.Members(domainId, groupId)
	if err != nil {
		return err
	}

	for _, member := range *rs {
		err = s.service.Subject().DeleteGroup(domainId, member.Id, groupId)
		if err != nil {
			return err
		}
	}

	// delete group
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId)

	if err = s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	log.Infof("domain(%+v) group deleted:%+v", domainId, groupId)

	return nil
}

func (s *groupSrvc) Update(domainId string, groupId string, o apiv1.Group) error {
	log.Infof("update domain(%+v) group(%+v):%+v", domainId, o.Id, o)

	// check existing
	if !s.Has(domainId, groupId) {
		return errors.NewErrorCf(errors.ErrCdDataExist, nil, "group(%v:%v)", domainId, groupId)
	}

	// update group
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, o.Id)
	err := s.service.update(path, o)
	if err != nil {
		return err
	}

	log.Infof("domain(%+v) group updated:%+v", domainId, groupId)

	return nil
}

func (s *groupSrvc) Get(domainId string, groupId string) (*apiv1.Group, error) {
	log.Infof("get domain(%+v) group(%+v)", domainId, groupId)

	group := apiv1.Group{}
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId)
	err := s.service.get(path, &group)
	if err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) group(%+v) got:%+v", domainId, groupId, group)

	return &group, nil
}

func (s *groupSrvc) List(domainId string) (*[]apiv1.Group, error) {
	log.Infof("list domain(%+v) groups", domainId)

	input := &DomainInput{
		DomainId: domainId,
	}
	groups := []apiv1.Group{}
	err := s.service.query(constant.RuleGroups, input, &groups)
	if err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) groups listed:%+v", domainId, groups)

	return &groups, nil
}

func (s *groupSrvc) Permissions(domainId string, groupId string) (*[]apiv1.PermissionOptions, error) {
	log.Infof("list domain(%+v) group(%+v) permissions", domainId, groupId)

	input := &DomainGroupInput{
		DomainId:     domainId,
		GroupPattern: groupId,
	}
	permissions := []apiv1.PermissionOptions{}
	err := s.service.query(constant.RuleGroupPermissions, input, &permissions)
	if err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) group(%+v) permissions listed:%+v", domainId, groupId, permissions)

	return &permissions, nil
}

func (s *groupSrvc) Roles(domainId string, groupId string) (*[]apiv1.RoleOptions, error) {
	log.Infof("list domain(%+v) group(%+v) roles", domainId, groupId)

	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldRoles)
	roles := make([]apiv1.RoleOptions, 0)
	if err := s.service.list(path, func(id string, j []byte) error {
		role := apiv1.RoleOptions{}
		if err := json.Unmarshal(j, &role); err != nil {
			return err
		}
		roles = append(roles, role)

		return nil
	}); err != nil {
		return nil, err
	}

	log.Infof("domain(%+v) group(%+v) roles listed:%+v", domainId, groupId, roles)

	return &roles, nil
}

func (s *groupSrvc) Denies(domainId string, groupId string) (*[]apiv1.DenyOptions, error) {
	log.Infof("list domain(%+v) group(%+v) denies", domainId, groupId)

	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldDenies)
	denies := make([]apiv1.DenyOptions, 0)
	if err := s.service.list(path, func(id string, j []byte) error {
		deny := apiv1.DenyOptions{}
		if err := json.Unmarshal(j, &deny); err != nil {
			return err
		}
		denies = append(denies, deny)

		return nil
	}); err != nil {
		return nil, err
	}
	log.Infof("domain(%+v) group(%+v) denies listed:%+v", domainId, groupId, denies)

	return &denies, nil
}

func (s *groupSrvc) Members(domainId string, groupId string) (*[]apiv1.SubjectOptions, error) {
	log.Infof("list domain(%+v) group(%+v) members", domainId, groupId)

	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldMembers)
	subjects := make([]apiv1.SubjectOptions, 0)
	if err := s.service.list(path, func(id string, j []byte) error {
		subject := apiv1.SubjectOptions{}
		if err := json.Unmarshal(j, &subject); err != nil {
			return err
		}
		subjects = append(subjects, subject)

		return nil
	}); err != nil {
		return nil, err
	}
	log.Infof("domain(%+v) group(%+v) subjects listed:%+v", domainId, groupId, subjects)

	return &subjects, nil
}

func (s *groupSrvc) DeleteRole(domainId string, groupId string, roleId string) error {
	log.Infof("delete domain(%+v) group(%+v) role(%+v)", domainId, groupId, roleId)
	// delete assigned role
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldRoles, roleId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	// delete related index
	if err := s.deleteRole(domainId, groupId, roleId); err != nil {
		return err
	}

	log.Infof("domain(%+v) group(%+v) role deleted:%+v", domainId, groupId, roleId)

	return nil
}

func (s *groupSrvc) DeleteDeny(domainId string, groupId string, denyId string) error {
	log.Infof("delete domain(%+v) group(%+v) deny(%+v)", domainId, groupId, denyId)
	// delete assigned deny
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldDenies, denyId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	// delete related index
	if err := s.deleteDeny(domainId, groupId, denyId); err != nil {
		return err
	}

	log.Infof("domain(%+v) group(%+v) deny deleted:%+v", domainId, groupId, denyId)

	return nil
}

func (s *groupSrvc) DeleteMember(domainId string, groupId string, memberId string) error {
	log.Infof("delete domain(%+v) group(%+v) member(%+v)", domainId, groupId, memberId)
	// delete assigned member
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldMembers, memberId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	// delete this group assigned to members
	path = BuildPath(constant.FldDomains, domainId, constant.FldSubjects, memberId,
		constant.FldGroups, groupId)
	if err := s.service.deleteWithPrefix(path); err != nil {
		return err
	}

	log.Infof("domain(%+v) group(%+v) member deleted:%+v", domainId, groupId, memberId)

	return nil
}

func (s *groupSrvc) AddRole(domainId string, groupId string, o apiv1.RoleOptions) error {
	log.Infof("add domain(%+v) group(%+v) role(%+v):%+v", domainId, groupId, o.Id, o)

	// check validity of role
	if err := s.checkRole(domainId, o.Id); err != nil {
		return err
	}

	// assign role
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldRoles, o.Id)
	if err := s.service.create(path, o); err != nil {
		return err
	}

	// create related index
	if err := s.addRole(domainId, groupId, o.Id); err != nil {
		return err
	}

	log.Infof("domain(%+v) group(%+v) role added:%+v", domainId, groupId, o)

	return nil
}

func (s *groupSrvc) AddDeny(domainId string, groupId string, o apiv1.DenyOptions) error {
	log.Infof("add domain(%+v) group(%+v) deny(%+v):%+v", domainId, groupId, o.Id, o)

	// check validity of deny
	if err := s.checkDeny(domainId, o.Id); err != nil {
		return err
	}

	// assign deny
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldDenies, o.Id)
	if err := s.service.create(path, o); err != nil {
		return err
	}

	// create related index
	if err := s.addDeny(domainId, groupId, o.Id); err != nil {
		return err
	}

	log.Infof("domain(%+v) group(%+v) deny added:%+v", domainId, groupId, o)

	return nil
}

func (s *groupSrvc) AddMember(domainId string, groupId string, o apiv1.SubjectOptions) error {
	log.Infof("add domain(%+v) group(%+v) member(%+v):%+v", domainId, groupId, o.Id, o)

	// check validity of deny
	if err := s.checkMember(domainId, o.Id); err != nil {
		return err
	}

	// assign member to group
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId,
		constant.FldMembers, o.Id)
	if err := s.service.create(path, o); err != nil {
		return err
	}

	// assign group to member
	path = BuildPath(constant.FldDomains, domainId, constant.FldSubjects, o.Id,
		constant.FldGroups, groupId)
	if err := s.service.create(path, o); err != nil {
		return err
	}

	log.Infof("domain(%+v) group(%+v) member added:%+v", domainId, groupId, o)

	return nil
}

func (s *groupSrvc) Has(domainId string, groupId string) bool {
	path := BuildPath(constant.FldDomains, domainId, constant.FldGroups, groupId)

	return s.service.exist(path)
}

func (s *groupSrvc) checkRole(domainId string, roleId string) error {
	if !s.service.Role().Has(domainId, roleId) {
		return errors.NewErrorC(errors.ErrCdInvalidRole, nil)
	}

	return nil
}

func (s *groupSrvc) addRole(domainId string, groupId string, roleId string) error {
	path := BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldRoles, roleId, constant.FldGroups, groupId)
	if err := s.service.ldb.WriteObject([]byte(path), []byte(groupId)); err != nil {
		return err
	}

	return nil
}

func (s *groupSrvc) deleteRole(domainId string, groupId string, roleId string) error {
	path := BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldRoles, roleId, constant.FldGroups, groupId)
	if err := s.service.ldb.Delete([]byte(path)); err != nil {
		return err
	}

	return nil
}

func (s *groupSrvc) checkDeny(domainId string, denyId string) error {
	if !s.service.Deny().Has(domainId, denyId) {
		return errors.NewErrorC(errors.ErrCdInvalidDeny, nil)
	}

	return nil
}

func (s *groupSrvc) addDeny(domainId string, groupId string, denyId string) error {
	path := BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldDenies, denyId, constant.FldGroups, groupId)
	if err := s.service.ldb.WriteObject([]byte(path), []byte(groupId)); err != nil {
		return err
	}

	return nil
}

func (s *groupSrvc) deleteDeny(domainId string, groupId string, denyId string) error {
	path := BuildPath(constant.FldIndexes, constant.FldDomains, domainId,
		constant.FldDenies, denyId, constant.FldGroups, groupId)
	if err := s.service.ldb.Delete([]byte(path)); err != nil {
		return err
	}

	return nil
}

func (s *groupSrvc) checkMember(domainId string, memberId string) error {
	if !s.service.Subject().Has(domainId, memberId) {
		return errors.NewErrorC(errors.ErrCdInvalidSubject, nil)
	}

	return nil
}
