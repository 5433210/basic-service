package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/log"
)

func Test_Role(t *testing.T) {
	log.Init(log.OptLevel(log.DebugLevel))
	s, _ := New(dbPath, regoPath, dataPath)
	defer s.Release()
	t.Run("", func(t *testing.T) {
		domain := apiv1.Domain{Id: "mydomain", Name: "我的域：wailik.com"}
		ds := s.Domain()
		roleSvrc := s.Role()

		ds.Create(domain)

		role := apiv1.Role{Id: "roleId", Name: "我的角色"}

		roleSvrc.Create(domain.Id, role)

		p, _ := roleSvrc.Get(domain.Id, role.Id)
		if !reflect.DeepEqual(role, *p) {
			t.Errorf("not equal between create(%v) and get(%v", role, *p)
		}

		role.Name = "新的角色"
		role.Permissions = apiv1.PermissionsInRole{Included: func() map[string]interface{} {
			ps := make(map[string]interface{})
			ps["read a book"] = "day,night"
			return ps
		}(), Excluded: func() []string {
			ps := make([]string, 0)
			ps = append(ps, "borrow a book")
			return ps
		}()}
		role.Scopes = apiv1.ScopesInRole{CanBeGranted: func() []string {
			ss := make([]string, 0)
			ss = append(ss, "/a1")
			ss = append(ss, "/a2")
			return ss
		}()}
		role.ExclusiveRoles = func() []string {
			ss := make([]string, 0)
			ss = append(ss, "/r1")
			ss = append(ss, "/r2")
			return ss
		}()
		role.DynamicallyIsolated = true
		roleSvrc.Update(domain.Id, role.Id, role)

		p, _ = roleSvrc.Get(domain.Id, role.Id)
		if !reflect.DeepEqual(role, *p) {
			t.Errorf("not equal between update(%v) and get(%v", role, *p)
		}

		role.DynamicallyIsolated = false
		role.ExclusiveRoles = nil
		roleSvrc.Update(domain.Id, role.Id, role)

		p, _ = roleSvrc.Get(domain.Id, role.Id)
		if !reflect.DeepEqual(role, *p) {
			t.Errorf("not equal between update(%v) and get(%v", role, *p)
		}

		r, _ := roleSvrc.List(domain.Id)

		if !reflect.DeepEqual(role, (*r)[0]) {
			t.Errorf("not equal between list(%v) and get(%v", role, (*r)[0])
		}
		c, _ := s.Dump("/")
		log.Debugf(string(c))

		roleSvrc.Delete(domain.Id, role.Id)

		p, err := roleSvrc.Get(domain.Id, role.Id)
		if err == nil {
			t.Errorf("delete error")
		}

	})
}
