package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/log"
)

func Test_Domain(t *testing.T) {
	log.Init(log.OptLevel(log.DebugLevel))
	s, _ := New(dbPath, regoPath, dataPath)
	defer s.Release()
	t.Run("", func(t *testing.T) {
		domain := apiv1.Domain{Id: "mydomain", Name: "我的域：wailik.com"}
		gots := make([]apiv1.Domain, 0)

		ds := s.Domain()
		// d, _ := ds.Get(constant.DefaultDomainId)
		// gots = append(gots, *d)

		ds.Create(domain)
		d, _ := ds.Get(domain.Id)
		if !reflect.DeepEqual(domain, *d) {
			t.Errorf("not equal between create(%v) and get(%v", domain, d)
		}

		gots = append(gots, *d)

		list, _ := ds.List()

		if len(*list) != 1 {
			t.Errorf("list error")
		}

		domain.Name = "变更域名了"
		ds.Update(domain.Id, domain)

		d, _ = ds.Get(domain.Id)

		if d.Name != domain.Name {
			t.Errorf("not equal between gots(%v) and updated(%v)", d, domain)
		}

		ds.Delete(domain.Id)
		d, err := ds.Get(domain.Id)

		if err == nil {
			t.Errorf("domain(%v) not deleted", domain.Id)
		}

	})
}
