package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/log"
)

func Test_Permission(t *testing.T) {
	log.Init(log.OptLevel(log.DebugLevel))
	s, _ := New(dbPath, regoPath, dataPath)
	defer s.Release()
	t.Run("", func(t *testing.T) {
		domain := apiv1.Domain{Id: "mydomain", Name: "我的域：wailik.com"}
		ds := s.Domain()
		ps := s.Permission()

		ds.Create(domain)

		permission := apiv1.Permission{Id: "permissionId", Options: "hello"}

		ps.Create(domain.Id, permission)

		p, _ := ps.Get(domain.Id, permission.Id)
		log.Debugf("%+v:%+v", permission, *p)
		if !reflect.DeepEqual(permission, *p) {
			t.Errorf("not equal between create(%v) and get(%v", permission, *p)
		}

		permission.Options = "world"
		ps.Update(domain.Id, permission.Id, permission)

		p, _ = ps.Get(domain.Id, permission.Id)
		if !reflect.DeepEqual(permission, *p) {
			t.Errorf("not equal between update(%v) and get(%v", permission, *p)
		}

		r, _ := ps.List(domain.Id)

		log.Debugf("list:%+v", r)

		if !reflect.DeepEqual(permission, (*r)[0]) {
			t.Errorf("not equal between list(%v) and get(%v", permission, (*r)[0])
		}

		ps.Delete(domain.Id, permission.Id)
		p, err := ps.Get(domain.Id, permission.Id)
		if err == nil {
			t.Errorf("delete error")
		}

	})
}
