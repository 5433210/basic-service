package servicev1

import (
	"reflect"
	"testing"

	apiv1 "wailik.com/internal/pkg/api/v1"
)

func Test_newSubjectSrvc(t *testing.T) {
	type args struct {
		s *service
	}
	tests := []struct {
		name string
		args args
		want *subjectSrvc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSubjectSrvc(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSubjectSrvc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_Create(t *testing.T) {
	type args struct {
		domainId string
		g        apiv1.Subject
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Create(tt.args.domainId, tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_Delete(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Delete(tt.args.domainId, tt.args.subjectId); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_Update(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
		g         apiv1.Subject
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Update(tt.args.domainId, tt.args.subjectId, tt.args.g); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_Get(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *apiv1.Subject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Get(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_List(t *testing.T) {
	type args struct {
		domainId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Subject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.domainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_Permissions(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Permission
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Permissions(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.Permissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.Permissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_Roles(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Roles(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.Roles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.Roles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_Denies(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Deny
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Denies(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.Denies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.Denies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_Groups(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Group
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Groups(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.Groups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.Groups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_DeleteRole(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
		roleId    string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.DeleteRole(tt.args.domainId, tt.args.subjectId, tt.args.roleId); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.DeleteRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_DeleteDeny(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
		denyId    string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.DeleteDeny(tt.args.domainId, tt.args.subjectId, tt.args.denyId); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.DeleteDeny() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_DeleteGroup(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
		groupId   string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.DeleteGroup(tt.args.domainId, tt.args.subjectId, tt.args.groupId); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_AddRole(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
		r         apiv1.RoleOptions
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddRole(tt.args.domainId, tt.args.subjectId, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.AddRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_AddDeny(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
		d         apiv1.DenyOptions
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddDeny(tt.args.domainId, tt.args.subjectId, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.AddDeny() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_AddGroup(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
		m         apiv1.GroupOptions
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddGroup(tt.args.domainId, tt.args.subjectId, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.AddGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subjectSrvc_RolesCanBeGrantedTo(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.RolesCanBeGrantedTo(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.RolesCanBeGrantedTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.RolesCanBeGrantedTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_RolesCanBeAccessedBy(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Role
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.RolesCanBeAccessedBy(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.RolesCanBeAccessedBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.RolesCanBeAccessedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_DeniesCanBeGrantedTo(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Deny
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.DeniesCanBeGrantedTo(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.DeniesCanBeGrantedTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.DeniesCanBeGrantedTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_subjectSrvc_DeniesCanBeAccessedBy(t *testing.T) {
	type args struct {
		domainId  string
		subjectId string
	}
	tests := []struct {
		name    string
		s       *subjectSrvc
		args    args
		want    *[]apiv1.Deny
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.DeniesCanBeAccessedBy(tt.args.domainId, tt.args.subjectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("subjectSrvc.DeniesCanBeAccessedBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subjectSrvc.DeniesCanBeAccessedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
