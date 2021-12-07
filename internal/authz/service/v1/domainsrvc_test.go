package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
)

func Test_newDomainSrvc(t *testing.T) {
	type args struct {
		s *service
	}
	tests := []struct {
		name string
		args args
		want *domainSrvc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDomainSrvc(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDomainSrvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_domainSrvc_Create(t *testing.T) {
	type args struct {
		d apiv1.Domain
	}
	tests := []struct {
		name    string
		s       *domainSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("domainSrvc.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domainSrvc_Delete(t *testing.T) {
	type args struct {
		domainId string
	}
	tests := []struct {
		name    string
		s       *domainSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Delete(tt.args.domainId); (err != nil) != tt.wantErr {
				t.Errorf("domainSrvc.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domainSrvc_Update(t *testing.T) {
	type args struct {
		domainId string
		d        apiv1.Domain
	}
	tests := []struct {
		name    string
		s       *domainSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Update(tt.args.domainId, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("domainSrvc.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_domainSrvc_Get(t *testing.T) {
	type args struct {
		domainId string
	}
	tests := []struct {
		name    string
		s       *domainSrvc
		args    args
		want    *apiv1.Domain
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Get(tt.args.domainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("domainSrvc.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("domainSrvc.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_domainSrvc_List(t *testing.T) {
	tests := []struct {
		name    string
		s       *domainSrvc
		want    *[]apiv1.Domain
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("domainSrvc.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("domainSrvc.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
