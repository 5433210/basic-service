package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/log"
	// servicev1 "wailik.com/internal/server/service/v1"
)

func Test_Deny(t *testing.T) {
	log.Init(log.OptLevel(log.DebugLevel))
	s, _ := NewService(dbPath, regoPath, dataPath)
	defer s.Release()
	t.Run("", func(t *testing.T) {
		domain := apiv1.Domain{Id: "mydomain", Name: "我的域：wailik.com"}
		ds := s.Domain()
		denySrvc := s.Deny()

		ds.Create(domain)

		deny := apiv1.Deny{Id: "roleId", Name: "我的禁止"}

		denySrvc.Create(domain.Id, deny)

		p, _ := denySrvc.Get(domain.Id, deny.Id)
		if !reflect.DeepEqual(deny, *p) {
			t.Errorf("not equal between create(%v) and get(%v", deny, *p)
		}

		deny.Name = "新的禁止"
		deny.Permissions = apiv1.PermissionsInDeny{Included: func() []string {
			ps := make([]string, 0)
			ps = append(ps, "read a book")
			return ps
		}(), Excluded: func() []string {
			ps := make([]string, 0)
			ps = append(ps, "borrow a book")
			return ps
		}()}

		deny.Scopes = apiv1.ScopesInDeny{CanBeGranted: func() []string {
			ss := make([]string, 0)
			ss = append(ss, "/a1")
			ss = append(ss, "/a2")
			return ss
		}()}

		denySrvc.Update(domain.Id, deny.Id, deny)

		p, _ = denySrvc.Get(domain.Id, deny.Id)
		if !reflect.DeepEqual(deny, *p) {
			t.Errorf("not equal between update(%v) and get(%v", deny, *p)
		}

		r, _ := denySrvc.List(domain.Id)

		if !reflect.DeepEqual(deny, (*r)[0]) {
			t.Errorf("not equal between list(%v) and get(%v", deny, (*r)[0])
		}

		denySrvc.Delete(domain.Id, deny.Id)
		p, err := denySrvc.Get(domain.Id, deny.Id)
		if err == nil {
			t.Errorf("delete error")
		}

	})
}
