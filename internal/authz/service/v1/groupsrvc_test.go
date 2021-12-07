package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
)

func Test_newGroupSrvc(t *testing.T) {
	type args struct {
		s *service
	}
	tests := []struct {
		name string
		args args
		want *groupSrvc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newGroupSrvc(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGroupSrvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupSrvc_Create(t *testing.T) {
	type args struct {
		domainId string
		g        apiv1.Group
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.domainId, tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupSrvc_Delete(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Delete(tt.args.domainId, tt.args.groupId); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupSrvc_Update(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
		g        apiv1.Group
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Update(tt.args.domainId, tt.args.groupId, tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupSrvc_Get(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		want    *apiv1.Group
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Get(tt.args.domainId, tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupSrvc.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupSrvc_List(t *testing.T) {
	type args struct {
		domainId string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		want    *[]apiv1.Group
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.domainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupSrvc.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupSrvc_Permissions(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		want    *[]apiv1.Permission
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Permissions(tt.args.domainId, tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.Permissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupSrvc.Permissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupSrvc_Roles(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		want    *[]apiv1.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Roles(tt.args.domainId, tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.Roles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupSrvc.Roles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupSrvc_Denies(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		want    *[]apiv1.Deny
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Denies(tt.args.domainId, tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.Denies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupSrvc.Denies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupSrvc_Members(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		want    *[]apiv1.Subject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Members(tt.args.domainId, tt.args.groupId)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.Members() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupSrvc.Members() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupSrvc_DeleteRole(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
		roleId   string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.DeleteRole(tt.args.domainId, tt.args.groupId, tt.args.roleId); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.DeleteRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupSrvc_DeleteDeny(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
		denyId   string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.DeleteDeny(tt.args.domainId, tt.args.groupId, tt.args.denyId); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.DeleteDeny() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupSrvc_DeleteMember(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
		memberId string
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.DeleteMember(tt.args.domainId, tt.args.groupId, tt.args.memberId); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.DeleteMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupSrvc_AddRole(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
		r        apiv1.RoleOptions
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddRole(tt.args.domainId, tt.args.groupId, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.AddRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupSrvc_AddDeny(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
		d        apiv1.DenyOptions
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddDeny(tt.args.domainId, tt.args.groupId, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.AddDeny() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_groupSrvc_AddMember(t *testing.T) {
	type args struct {
		domainId string
		groupId  string
		m        apiv1.SubjectOptions
	}
	tests := []struct {
		name    string
		s       *groupSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddMember(tt.args.domainId, tt.args.groupId, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("groupSrvc.AddMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
