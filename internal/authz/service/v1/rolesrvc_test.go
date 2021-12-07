package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
)

func Test_newRoleSrvc(t *testing.T) {
	type args struct {
		s *service
	}
	tests := []struct {
		name string
		args args
		want *roleSrvc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRoleSrvc(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRoleSrvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roleSrvc_Create(t *testing.T) {
	type args struct {
		domainId string
		r        apiv1.Role
	}
	tests := []struct {
		name    string
		s       *roleSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.domainId, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("roleSrvc.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_roleSrvc_Delete(t *testing.T) {
	type args struct {
		domainId string
		roleId   string
	}
	tests := []struct {
		name    string
		s       *roleSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Delete(tt.args.domainId, tt.args.roleId); (err != nil) != tt.wantErr {
				t.Errorf("roleSrvc.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_roleSrvc_Update(t *testing.T) {
	type args struct {
		domainId string
		roleId   string
		g        apiv1.Role
	}
	tests := []struct {
		name    string
		s       *roleSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Update(tt.args.domainId, tt.args.roleId, tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("roleSrvc.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_roleSrvc_Get(t *testing.T) {
	type args struct {
		domainId string
		roleId   string
	}
	tests := []struct {
		name    string
		s       *roleSrvc
		args    args
		want    *apiv1.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Get(tt.args.domainId, tt.args.roleId)
			if (err != nil) != tt.wantErr {
				t.Errorf("roleSrvc.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roleSrvc.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roleSrvc_List(t *testing.T) {
	type args struct {
		domainId string
	}
	tests := []struct {
		name    string
		s       *roleSrvc
		args    args
		want    *[]apiv1.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.domainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("roleSrvc.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roleSrvc.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
