package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/log"
)

func Test_Subject(t *testing.T) {
	log.Init(log.OptLevel(log.DebugLevel))
	s, _ := New(dbPath, regoPath, dataPath)
	defer s.Release()
	t.Run("", func(t *testing.T) {
		domain := apiv1.Domain{Id: "mydomain", Name: "我的域：wailik.com"}
		ds := s.Domain()
		subjectSrvc := s.Subject()
		permissionSrvc := s.Permission()
		denySrvc := s.Deny()
		roleSrvc := s.Role()

		ds.Create(domain)

		subject := apiv1.Subject{Id: "subjectId"}
		subjectSrvc.Create(domain.Id, subject)

		permissionSrvc.Create(domain.Id, apiv1.Permission{Id: "readabook", Options: "permission option"})
		permissionSrvc.Create(domain.Id, apiv1.Permission{Id: "buyabook", Options: "permission option"})
		permissionSrvc.Create(domain.Id, apiv1.Permission{Id: "borrowabook", Options: "permission option"})

		roleSrvc.Create(domain.Id, apiv1.Role{Id: "roleId", Permissions: apiv1.PermissionsInRole{Included: func() map[string]interface{} {
			ps := make(map[string]interface{})
			ps["readabook"] = "day,night"
			ps["buyabook"] = "day,night"
			ps["borrowabook"] = "day"
			return ps
		}()}, Scopes: apiv1.ScopesInRole{CanBeGranted: func() []string {
			scopes := make([]string, 2)
			scopes[0] = "groupId*"
			return scopes
		}(), CanBeAccessed: func() []string {
			scopes := make([]string, 2)
			scopes[0] = "groupId*"
			return scopes
		}()}})

		err := denySrvc.Create(domain.Id, apiv1.Deny{Id: "denyId", Permissions: apiv1.PermissionsInDeny{Included: func() []string {
			ps := make([]string, 1)
			ps[0] = "borrowabook"
			return ps
		}()}, Scopes: apiv1.ScopesInDeny{CanBeGranted: func() []string {
			scopes := make([]string, 1)
			scopes[0] = "groupId*"
			return scopes
		}(), CanBeAccessed: func() []string {
			scopes := make([]string, 1)
			scopes[0] = "groupId*"
			return scopes
		}()}})

		if err != nil {
			t.Error(err)

		}
		p, _ := subjectSrvc.Get(domain.Id, subject.Id)
		if !reflect.DeepEqual(subject, *p) {
			t.Errorf("not equal between create(%v) and get(%v)", subject, *p)

		}

		r, _ := subjectSrvc.List(domain.Id)

		if !reflect.DeepEqual(subject, (*r)[0]) {
			t.Errorf("not equal between list(%v) and get(%v)", subject, (*r)[0])

		}

		do := apiv1.DenyOptions{Id: "denyId", Options: "deny option"}
		subjectSrvc.AddDeny(domain.Id, subject.Id, do)

		c, _ := s.Dump("/")
		log.Debugf("%v", string(c))
		dr, _ := subjectSrvc.Denies(domain.Id, subject.Id)

		if !reflect.DeepEqual(do, (*dr)[0]) {
			t.Errorf("not equal between list(%+v) and get(%+v)", do, (*dr)[0])

		}

		subjectSrvc.DeleteDeny(domain.Id, subject.Id, do.Id)

		dr, _ = subjectSrvc.Denies(domain.Id, subject.Id)

		if len(*dr) != 0 {
			t.Errorf("delete deny option error")

		}

		ro := apiv1.RoleOptions{Id: "roleId", Options: "role option"}
		subjectSrvc.AddRole(domain.Id, subject.Id, ro)
		rr, _ := subjectSrvc.Roles(domain.Id, subject.Id)

		if !reflect.DeepEqual(ro, (*rr)[0]) {
			t.Errorf("not equal between list(%+v) and get(%+v)", ro, (*rr)[0])

		}

		subjectSrvc.DeleteRole(domain.Id, subject.Id, ro.Id)

		rr, _ = subjectSrvc.Roles(domain.Id, subject.Id)

		if len(*rr) != 0 {
			t.Errorf("delete role option error")

		}

		gpo := apiv1.GroupOptions{Id: "groupId", Options: "group option"}
		subjectSrvc.AddGroup(domain.Id, subject.Id, gpo)
		gr, _ := subjectSrvc.Groups(domain.Id, subject.Id)

		if !reflect.DeepEqual(gpo, (*gr)[0]) {
			t.Errorf("not equal between list(%+v) and get(%+v)", gpo, (*gr)[0])

		}

		subjectSrvc.DeleteGroup(domain.Id, subject.Id, gpo.Id)

		gr, _ = subjectSrvc.Groups(domain.Id, subject.Id)

		if len(*gr) != 0 {
			t.Errorf("delete subject option error")
		}

		subjectSrvc.AddRole(domain.Id, subject.Id, apiv1.RoleOptions{Id: "roleId", Options: "options"})

		err = subjectSrvc.AddDeny(domain.Id, subject.Id, apiv1.DenyOptions{Id: "denyId", Options: "options"})
		if err != nil {
			t.Error(err)

		}

		gp, _ := subjectSrvc.Permissions(domain.Id, subject.Id)

		if gp == nil || !((*gp)[0].Id == "readabook" || (*gp)[1].Id == "readabook") {
			t.Errorf("group permissions error:%+v", gp)
		}

		groupSrvc := s.Group()
		group := apiv1.Group{Id: "groupId"}
		groupSrvc.Create(domain.Id, group)
		// groupSrvc.AddMember(domain.Id, group.Id, apiv1.SubjectOptions{Id: subject.Id})
		subjectSrvc.AddGroup(domain.Id, subject.Id, apiv1.GroupOptions{Id: group.Id})

		c, _ = s.Dump("/")
		log.Debugf(string(c))

		do2 := apiv1.DenyOptions{Id: "denyId"}
		da, _ := subjectSrvc.DeniesCanBeAccessedBy(domain.Id, subject.Id)
		if !reflect.DeepEqual(do2, (*da)[0]) {
			t.Errorf("not equal between list(%+v) and get(%+v)", do2, (*da)[0])
		}
		da, _ = subjectSrvc.DeniesCanBeGrantedTo(domain.Id, subject.Id)
		if !reflect.DeepEqual(do2, (*da)[0]) {
			t.Errorf("not equal between list(%+v) and get(%+v)", do2, (*da)[0])
		}

		dr2 := apiv1.RoleOptions{Id: "roleId"}
		ra, _ := subjectSrvc.RolesCanBeAccessedBy(domain.Id, subject.Id)
		if !reflect.DeepEqual(dr2, (*ra)[0]) {
			t.Errorf("not equal between list(%+v) and get(%+v)", dr2, (*ra)[0])
		}
		ra, _ = subjectSrvc.RolesCanBeGrantedTo(domain.Id, subject.Id)
		if !reflect.DeepEqual(dr2, (*ra)[0]) {
			t.Errorf("not equal between list(%+v) and get(%+v)", dr2, (*ra)[0])
		}

		subject.Options = "new"
		subjectSrvc.Update(domain.Id, subject.Id, subject)
		p, _ = subjectSrvc.Get(domain.Id, subject.Id)
		if !reflect.DeepEqual(subject, *p) {
			t.Errorf("not equal between update(%v) and get(%v)", subject, *p)
		}

		subjectSrvc.Delete(domain.Id, subject.Id)
		p, err = subjectSrvc.Get(domain.Id, subject.Id)
		if err == nil {
			t.Errorf("delete error")
		}
	})
}
