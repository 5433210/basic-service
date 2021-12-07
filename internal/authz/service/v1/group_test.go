package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/log"
)

func Test_Group(t *testing.T) {
	log.Init(log.OptLevel(log.DebugLevel))
	s, _ := NewService(dbPath, regoPath, dataPath)
	defer s.Release()
	t.Run("", func(t *testing.T) {
		var err error
		domain := apiv1.Domain{Id: "mydomain", Name: "我的域：wailik.com"}
		ds := s.Domain()
		groupSrvc := s.Group()
		permissionSrvc := s.Permission()
		subjectSrvc := s.Subject()

		ds.Create(domain)

		group := apiv1.Group{Id: "groupId"}

		groupSrvc.Create(domain.Id, group)

		subjectSrvc.Create(domain.Id, apiv1.Subject{Id: "subjectId"})

		permissionSrvc.Create(domain.Id, apiv1.Permission{Id: "readabook", Options: "permission option"})

		roleSrvc := s.Role()
		roleSrvc.Create(domain.Id, apiv1.Role{Id: "roleId", Permissions: apiv1.PermissionsInRole{Included: func() map[string]interface{} {
			ps := make(map[string]interface{})
			ps["readabook"] = "day,night"
			return ps
		}()}})

		denySrvc := s.Deny()
		denySrvc.Create(domain.Id, apiv1.Deny{Id: "denyId"})

		p, _ := groupSrvc.Get(domain.Id, group.Id)
		if !reflect.DeepEqual(group, *p) {
			t.Errorf("not equal between create(%v) and get(%v)", group, *p)
		}

		r, _ := groupSrvc.List(domain.Id)

		if !reflect.DeepEqual(group, (*r)[0]) {
			t.Errorf("not equal between list(%v) and get(%v)", group, (*r)[0])
		}

		do := apiv1.DenyOptions{Id: "denyId", Options: "deny option"}
		err = groupSrvc.AddDeny(domain.Id, group.Id, do)
		if err != nil {
			t.Error(err)
			return
		}
		dr, _ := groupSrvc.Denies(domain.Id, group.Id)

		if !reflect.DeepEqual(do, (*dr)[0]) {
			t.Errorf("not equal between list(%v) and get(%v)", do, (*dr)[0])
		}

		groupSrvc.DeleteDeny(domain.Id, group.Id, do.Id)

		dr, _ = groupSrvc.Denies(domain.Id, group.Id)

		if len(*dr) != 0 {
			t.Errorf("delete deny option error")
		}

		ro := apiv1.RoleOptions{Id: "roleId", Options: "role option"}
		groupSrvc.AddRole(domain.Id, group.Id, ro)
		rr, _ := groupSrvc.Roles(domain.Id, group.Id)

		if !reflect.DeepEqual(ro, (*rr)[0]) {
			t.Errorf("not equal between list(%v) and get(%v)", ro, (*rr)[0])
		}

		groupSrvc.DeleteRole(domain.Id, group.Id, ro.Id)

		rr, _ = groupSrvc.Roles(domain.Id, group.Id)

		if len(*rr) != 0 {
			t.Errorf("delete role option error")
		}

		so := apiv1.SubjectOptions{Id: "subjectId", Options: "subject option"}
		groupSrvc.AddMember(domain.Id, group.Id, so)
		sr, _ := groupSrvc.Members(domain.Id, group.Id)

		if !reflect.DeepEqual(so, (*sr)[0]) {
			t.Errorf("not equal between list(%v) and get(%v)", so, (*sr)[0])
		}

		groupSrvc.DeleteMember(domain.Id, group.Id, so.Id)

		sr, _ = groupSrvc.Members(domain.Id, group.Id)

		if len(*sr) != 0 {
			t.Errorf("delete subject option error")
		}

		groupSrvc.AddRole(domain.Id, group.Id, apiv1.RoleOptions{Id: "roleId", Options: "options"})

		c, _ := s.Dump("/")

		log.Debugf(string(c))

		gp, _ := groupSrvc.Permissions(domain.Id, group.Id)

		if gp == nil || (*gp)[0].Id != "readabook" {
			t.Errorf("group permissions error")
		}

		group.Options = "new"
		groupSrvc.Update(domain.Id, group.Id, group)
		p, _ = groupSrvc.Get(domain.Id, group.Id)
		if !reflect.DeepEqual(group, *p) {
			t.Errorf("not equal between update(%v) and get(%v)", group, *p)
		}

		groupSrvc.Delete(domain.Id, group.Id)
		p, err = groupSrvc.Get(domain.Id, group.Id)
		if err == nil {
			t.Errorf("delete error")
		}
	})
}
