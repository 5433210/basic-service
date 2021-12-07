package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
)

func Test_newPermissionSrvc(t *testing.T) {
	type args struct {
		s *service
	}
	tests := []struct {
		name string
		args args
		want *permissionSrvc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPermissionSrvc(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPermissionSrvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_permissionSrvc_Create(t *testing.T) {
	type args struct {
		domainId string
		p        apiv1.Permission
	}
	tests := []struct {
		name    string
		s       *permissionSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.domainId, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("permissionSrvc.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_permissionSrvc_Delete(t *testing.T) {
	type args struct {
		domainId     string
		permissionId string
	}
	tests := []struct {
		name    string
		s       *permissionSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Delete(tt.args.domainId, tt.args.permissionId); (err != nil) != tt.wantErr {
				t.Errorf("permissionSrvc.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_permissionSrvc_Update(t *testing.T) {
	type args struct {
		domainId     string
		permissionId string
		p            apiv1.Permission
	}
	tests := []struct {
		name    string
		s       *permissionSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Update(tt.args.domainId, tt.args.permissionId, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("permissionSrvc.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_permissionSrvc_Get(t *testing.T) {
	type args struct {
		domainId     string
		permissionId string
	}
	tests := []struct {
		name    string
		s       *permissionSrvc
		args    args
		want    *apiv1.Permission
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Get(tt.args.domainId, tt.args.permissionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("permissionSrvc.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("permissionSrvc.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_permissionSrvc_List(t *testing.T) {
	type args struct {
		domainId string
	}
	tests := []struct {
		name    string
		s       *permissionSrvc
		args    args
		want    *[]apiv1.Permission
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.domainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("permissionSrvc.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("permissionSrvc.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
