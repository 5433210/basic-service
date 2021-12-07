package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
)

func Test_newDenySrvc(t *testing.T) {
	type args struct {
		s *service
	}
	tests := []struct {
		name string
		args args
		want *denySrvc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDenySrvc(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDenySrvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_denySrvc_Create(t *testing.T) {
	type args struct {
		domainId string
		r        apiv1.Deny
	}
	tests := []struct {
		name    string
		s       *denySrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.domainId, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("denySrvc.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_denySrvc_Delete(t *testing.T) {
	type args struct {
		domainId string
		denyId   string
	}
	tests := []struct {
		name    string
		s       *denySrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Delete(tt.args.domainId, tt.args.denyId); (err != nil) != tt.wantErr {
				t.Errorf("denySrvc.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_denySrvc_Update(t *testing.T) {
	type args struct {
		domainId string
		denyId   string
		d        apiv1.Deny
	}
	tests := []struct {
		name    string
		s       *denySrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Update(tt.args.domainId, tt.args.denyId, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("denySrvc.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_denySrvc_Get(t *testing.T) {
	type args struct {
		domainId string
		denyId   string
	}
	tests := []struct {
		name    string
		s       *denySrvc
		args    args
		want    *apiv1.Deny
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Get(tt.args.domainId, tt.args.denyId)
			if (err != nil) != tt.wantErr {
				t.Errorf("denySrvc.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("denySrvc.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_denySrvc_List(t *testing.T) {
	type args struct {
		domainId string
	}
	tests := []struct {
		name    string
		s       *denySrvc
		args    args
		want    *[]apiv1.Deny
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.domainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("denySrvc.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("denySrvc.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
