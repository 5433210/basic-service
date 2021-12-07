package opa

import (
	"context"
	"reflect"
	"testing"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/storage"
	apiv1 "wailik.com/internal/pkg/api/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/log"
)

var regoPath = "/opt/rbac/data/rbac.rego"
var dataPath = "/opt/rbac/data/data.json"

func Test_opaTxn_Ctx(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		ctx context.Context
		txn storage.Transaction
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		// TODO: Create test cases.
		{"context is null", fields{nil, nil}, nil},
		{"context is ok", fields{ctx, nil}, ctx},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := opaTxn{
				ctx: tt.fields.ctx,
				txn: tt.fields.txn,
			}
			if got := txn.Ctx(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opaTxn.Ctx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opaTxn_Txn(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	ctx := context.Background()
	txn, _ := db.mem.NewTransaction(ctx, storage.TransactionParams{Write: true})
	type fields struct {
		ctx context.Context
		txn storage.Transaction
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Transaction
	}{
		// TODO: Create test cases.
		{"txn is null", fields{nil, nil}, nil},
		{"txn is ok", fields{nil, txn}, txn},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := opaTxn{
				ctx: tt.fields.ctx,
				txn: tt.fields.txn,
			}
			if got := txn.Txn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opaTxn.Txn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opaStore_Commit(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	txn, _ := db.NewTxn(true)
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		txn OpaTxn
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"db commit with null mem store", fields{nil, nil}, args{txn: txn}, true},
		{"db commit with null txn", fields{db.mem, nil}, args{txn: nil}, true},
		{"db commit with txn", fields{db.mem, nil}, args{txn: txn}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.Commit(tt.args.txn); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_opaStore_Abort(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	txn, _ := db.NewTxn(true)
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		txn OpaTxn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Create test cases.
		{"db abort with null mem store", fields{nil, nil}, args{txn: txn}},
		{"db abort with null txn", fields{db.mem, nil}, args{txn: nil}},
		{"db abort with txn", fields{db.mem, nil}, args{txn: txn}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			s.Abort(tt.args.txn)
		})
	}
}

func Test_opaStore_Create(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})

	txn, _ := db.NewTxn(true)
	path := "/domains/newdomainid"
	value := "value"
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		txn   OpaTxn
		path  string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"db add with null mem store", fields{nil, nil}, args{txn: txn, path: path, value: value}, true},
		{"db add with null txn", fields{db.mem, nil}, args{txn: nil, path: path, value: value}, true},
		{"db add with txn", fields{db.mem, nil}, args{txn: txn, path: path, value: value}, false},
		{"db add with txn repeatly", fields{db.mem, nil}, args{txn: txn, path: path, value: value}, false},
		{"db add with txn and null path", fields{db.mem, nil}, args{txn: txn, path: "", value: value}, true},
		{"db add with txn and null value", fields{db.mem, nil}, args{txn: txn, path: path, value: nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.Create(tt.args.txn, tt.args.path, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_opaStore_Replace(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	txn, _ := db.NewTxn(true)
	path := "/domains/78839721-a274-4a01-a2be-2725903bcf82"
	value := "value"
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		txn   OpaTxn
		path  string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"db replace with null mem store", fields{nil, nil}, args{txn: txn, path: path, value: value}, true},
		{"db replace with null txn", fields{db.mem, nil}, args{txn: nil, path: path, value: value}, true},
		{"db replace with txn and path", fields{db.mem, nil}, args{txn: txn, path: path, value: value}, false},
		{"db replace with txn and path repeatly", fields{db.mem, nil}, args{txn: txn, path: path, value: value}, false},
		{"db replace with txn and null path", fields{db.mem, nil}, args{txn: txn, path: "", value: value}, true},
		{"db replace with txn and null value", fields{db.mem, nil}, args{txn: txn, path: path, value: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.replace(tt.args.txn, tt.args.path, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func Test_opaStore_Delete(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	txn, _ := db.NewTxn(true)
	path := "/domains/78839721-a274-4a01-a2be-2725903bcf82"
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		txn  OpaTxn
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"db delete with null mem store", fields{nil, nil}, args{txn: txn, path: path}, true},
		{"db delete with null txn", fields{db.mem, nil}, args{txn: nil, path: path}, true},
		{"db delete with txn and null path", fields{db.mem, nil}, args{txn: txn, path: ""}, true},
		{"db delete with txn and path", fields{db.mem, nil}, args{txn: txn, path: path}, false},
		{"db delete with txn and path repeatly", fields{db.mem, nil}, args{txn: txn, path: path}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.Delete(tt.args.txn, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_opaStore_NewTxn(t *testing.T) {
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		writable bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    OpaTxn
		wantErr bool
	}{
		// TODO: Create test cases.
		{"new transaction with null mem store", fields{nil, nil}, args{writable: true}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			got, err := s.NewTxn(tt.args.writable)
			if (err != nil) != tt.wantErr {
				t.Errorf("opaStore.NewTxn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opaStore.NewTxn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opaStore_Read(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	txn, _ := db.NewTxn(true)
	path := "/domains/78839721-a274-4a01-a2be-2725903bcf82/name"
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		txn  OpaTxn
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Create test cases.
		{"db read with null mem store", fields{nil, nil}, args{txn: txn, path: path}, nil, true},
		{"db read with null txn", fields{db.mem, nil}, args{txn: nil, path: path}, nil, true},
		{"db read with txn", fields{db.mem, nil}, args{txn: txn, path: path}, "财务管理系统", false},
		{"db read with txn repeatly", fields{db.mem, nil}, args{txn: txn, path: path}, "财务管理系统", false},
		{"db read with txn and null path", fields{db.mem, nil}, args{txn: txn, path: ""}, nil, true},
		{"db read with txn and invalid path", fields{db.mem, nil}, args{txn: txn, path: "/this/is/an/invalid/path"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			got, err := s.Read(tt.args.txn, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("opaStore.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opaStore.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_opaStore_Write(t *testing.T) {
// 	type object struct {
// 		name    string
// 		pObject *object
// 	}
// 	obj := object{
// 		name: "parent",
// 		pObject: &object{
// 			name:    "child",
// 			pObject: nil,
// 		},
// 	}
// 	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
// 	txn, _ := db.NewTxn(true)
// 	path := "/domains/78839721-a274-4a01-a2be-2725903bcf82"
// 	value := obj
// 	type fields struct {
// 		mem storage.Store
// 		cmp *ast.Compiler
// 	}
// 	type args struct {
// 		txn  OpaTxn
// 		path string
// 		o    interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Create test cases.
// 		{"db write with null mem store", fields{nil, nil}, args{txn: txn, path: path, o: value}, true},
// 		{"db write with null txn", fields{db.mem, nil}, args{txn: nil, path: path, o: value}, true},
// 		{"db write with txn ", fields{db.mem, nil}, args{txn: txn, path: path, o: value}, false},
// 		{"db write with txn repeatly", fields{db.mem, nil}, args{txn: txn, path: path, o: value}, false},
// 		{"db write with txn and null path", fields{db.mem, nil}, args{txn: txn, path: "", o: value}, true},
// 		{"db write with txn and null value", fields{db.mem, nil}, args{txn: txn, path: path, o: nil}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &opaStore{
// 				mem: tt.fields.mem,
// 				cmp: tt.fields.cmp,
// 			}
// 			if err := s.Write(tt.args.txn, tt.args.path, tt.args.o); (err != nil) != tt.wantErr {
// 				t.Errorf("opaStore.Write() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func Test_opaStore_ReadOne(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	path := "/domains/78839721-a274-4a01-a2be-2725903bcf82"

	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Create test cases.
		{"db readone with null mem store", fields{nil, nil}, args{path: path}, nil, true},
		{"db readone with null path", fields{db.mem, nil}, args{path: ""}, nil, true},
		{"db readone with invalid path", fields{db.mem, nil}, args{path: "/this/is/an/invalid/path"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			got, err := s.ReadOne(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("opaStore.ReadOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("opaStore.ReadOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opaStore_WriteOne(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	path := "/domains/78839721-a274-4a01-a2be-2725903bcf82"
	type object struct {
		name    string
		pObject *object
	}
	obj := object{
		name: "parent",
		pObject: &object{
			name:    "child",
			pObject: nil,
		},
	}
	value := obj
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		path string
		o    interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"db writeone with null mem store", fields{nil, nil}, args{path: path, o: value}, true},
		{"db writeone with path and value", fields{db.mem, nil}, args{path: path, o: value}, false},
		{"db writeone with null path", fields{db.mem, nil}, args{path: "", o: value}, true},
		{"db writeone with null value", fields{db.mem, nil}, args{path: path, o: nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.CreateOne(tt.args.path, tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.WriteOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_opaStore_Query(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	rule := constant.RuleDomains
	rullNull := constant.RuleNULL
	out := make([]apiv1.Domain, 0)
	rs, _ := db.ReadOne("/domains/78839721-a274-4a01-a2be-2725903bcf82")
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		rule   string
		input  interface{}
		output interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"db query with null mem store and null compile module", fields{nil, nil}, args{rule: rule, input: nil, output: out}, true},
		{"db query with mem store and null compile module", fields{db.mem, nil}, args{rule: rule, input: nil, output: out}, true},
		{"db query with null mem store and compile module", fields{nil, db.cmp}, args{rule: rule, input: nil, output: out}, true},
		{"db query with null rule", fields{db.mem, db.cmp}, args{rule: "", input: nil, output: out}, true},
		{"db query with null output", fields{db.mem, db.cmp}, args{rule: rule, input: nil, output: nil}, true},
		{"db query with invalid rule", fields{db.mem, db.cmp}, args{rule: "invalid rule", input: nil, output: out}, true},
		{"db query with null rule", fields{db.mem, db.cmp}, args{rule: rullNull, input: nil, output: &out}, false},
		{"db query with rule, input, output", fields{db.mem, db.cmp}, args{rule: rule, input: nil, output: &out}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.Query(tt.args.rule, tt.args.input, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.Query() error = %v, wantErr %v, rs %+v", err, tt.wantErr, rs)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts *Options
	}
	tests := []struct {
		name    string
		args    args
		want    *opaStore
		wantErr bool
	}{
		// TODO: Create test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opaStore_load(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})

	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		regoPath string
		dataPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null mem store", fields{nil, nil}, args{regoPath: regoPath, dataPath: dataPath}, true},
		{"null regoPath", fields{db.mem, nil}, args{regoPath: "", dataPath: dataPath}, true},
		{"null dataPath", fields{db.mem, nil}, args{regoPath: regoPath, dataPath: ""}, true},
		{"null regoPath & dataPath", fields{db.mem, nil}, args{regoPath: "", dataPath: ""}, true},
		{"invalid regoPath", fields{db.mem, nil}, args{regoPath: "/invalid/rego/path", dataPath: dataPath}, true},
		{"invalid dataPath", fields{db.mem, nil}, args{regoPath: regoPath, dataPath: "/invalid/data/path"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.Load(tt.args.regoPath, tt.args.dataPath); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_opaStore_loadRego(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})

	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		regoPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null mem store", fields{nil, nil}, args{regoPath: regoPath}, true},
		{"null regoPath", fields{db.mem, nil}, args{regoPath: ""}, true},
		{"invalid regoPath", fields{db.mem, nil}, args{regoPath: "/invalid/rego/path"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.loadRego(tt.args.regoPath); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.readRego() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_opaStore_loadData(t *testing.T) {
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})

	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		dataPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null mem store", fields{nil, nil}, args{dataPath: dataPath}, true},
		{"null dataPath", fields{db.mem, nil}, args{dataPath: ""}, true},
		{"invalid dataPath", fields{db.mem, nil}, args{dataPath: "/invalid/data/path"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.loadData(tt.args.dataPath); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.readData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_opaStore_Update(t *testing.T) {
	log.Init(log.OptLevel(log.DebugLevel))
	db, _ := New(&Options{RegoPath: regoPath, DataPath: dataPath})
	txn, _ := db.NewTxn(true)
	path := "/domains/78839721-a274-4a01-a2be-2725903bcf82"
	value := "value"
	value2 := apiv1.Role{}
	type fields struct {
		mem storage.Store
		cmp *ast.Compiler
	}
	type args struct {
		txn   OpaTxn
		path  string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"db Update with null mem store", fields{nil, nil}, args{txn: txn, path: path, value: value}, true},
		{"db Update with null txn", fields{db.mem, nil}, args{txn: nil, path: path, value: value}, true},
		{"db Update with txn and path", fields{db.mem, nil}, args{txn: txn, path: path, value: value}, false},
		{"db Update with txn and path, struct", fields{db.mem, nil}, args{txn: txn, path: path, value: value2}, false},
		{"db Update with txn and path repeatly", fields{db.mem, nil}, args{txn: txn, path: path, value: value}, false},
		{"db Update with txn and null path", fields{db.mem, nil}, args{txn: txn, path: "", value: value}, true},
		{"db Update with txn and null value", fields{db.mem, nil}, args{txn: txn, path: path, value: nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &opaStore{
				mem: tt.fields.mem,
				cmp: tt.fields.cmp,
			}
			if err := s.Update(tt.args.txn, tt.args.path, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("opaStore.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
