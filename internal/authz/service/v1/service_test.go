package servicev1

import (
	"fmt"
	"reflect"
	"testing"

	"wailik.com/internal/authz/store/leveldb"
	"wailik.com/internal/authz/store/opa"
)

// var regoPath = "/opt/rbac/data/rbac.rego"
// var dataPath = "/opt/rbac/data/data.json"
// var dbPath = "/opt/rbac/db"
var errAction = func(txn opa.OpaTxn) (interface{}, error) { return nil, fmt.Errorf("action error") }
var okAction = func(txn opa.OpaTxn) (interface{}, error) { return "ok", nil }
var errElement = func(id string, json []byte) error { return fmt.Errorf("element error") }
var okElement = func(id string, json []byte) error { return nil }
var testPath = "/path"
var jsonValue = `{"sp1":{"f1":"v1"}, "sp2":{"f2":"v2"}}`

type p struct {
	Sp1 sp1 `json:"sp1"`
	Sp2 sp2 `json:"sp2"`
}

type sp1 struct {
	F1 string `json:"f1"`
}

type sp2 struct {
	F2 string `json:"f2"`
}

func TestNewService(t *testing.T) {
	type args struct {
		dbPath   string
		regoPath string
		dataPath string
	}
	tests := []struct {
		name    string
		args    args
		want    Service
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null args", args{dbPath: "", regoPath: "", dataPath: ""}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.args.dbPath, tt.args.regoPath, tt.args.dataPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_LoadData(t *testing.T) {
	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})

	defer ldb.Close()
	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
	}

	tests := []struct {
		name   string
		fields fields

		wantErr bool
	}{
		// TODO: Create test cases.
		{"null args", fields{nil, nil}, true},
		{"null ldb", fields{nil, odb}, true},
		{"null odb", fields{ldb, nil}, true},
		// {"null regoPath", fields{ldb, odb}, args{regoPath: "", dataPath: dataPath}, true},
		// {"null dataPath", fields{ldb, odb}, args{regoPath: regoPath, dataPath: ""}, true},
		// {"invalid regoPath", fields{ldb, odb}, args{regoPath: "/invalid/path", dataPath: dataPath}, true},
		// {"invalid dataPath", fields{ldb, odb}, args{regoPath: regoPath, dataPath: "/invalid/path"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}
			if err := s.LoadData(); (err != nil) != tt.wantErr {
				t.Errorf("service.LoadData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_doActionWithTxn(t *testing.T) {
	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})
	defer ldb.Close()
	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
	}
	type args struct {
		write  bool
		action Action
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null action", fields{ldb: ldb, odb: odb}, args{true, nil}, nil, true},
		{"error action", fields{ldb: ldb, odb: odb}, args{true, errAction}, nil, true},
		{"write ok action ", fields{ldb: ldb, odb: odb}, args{true, okAction}, "ok", false},
		{"read ok action", fields{ldb: ldb, odb: odb}, args{false, okAction}, "ok", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}
			got, err := s.doActionWithTxn(tt.args.write, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.doActionWithTxn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.doActionWithTxn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_list(t *testing.T) {
	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})

	mapValue := make(map[string]interface{}, 0)
	txn, _ := odb.NewTxn(true)
	json.Unmarshal([]byte(jsonValue), &mapValue)
	odb.Create(txn, testPath, mapValue)
	ldb.WriteObject([]byte(testPath), mapValue)
	odb.Commit(txn)

	defer func() {
		txn, _ := odb.NewTxn(true)
		odb.Delete(txn, testPath)
		ldb.Delete([]byte(testPath))
		odb.Commit(txn)
		ldb.Close()
	}()

	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
	}
	type args struct {
		path string
		elem element
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null path", fields{ldb, odb}, args{"", okElement}, true},
		{"invalid path", fields{ldb, odb}, args{"/i/am/invalid/path", okElement}, true},
		{"null element", fields{ldb, odb}, args{testPath + "/sp1", nil}, true},
		{"error element", fields{ldb, odb}, args{testPath + "/sp1", errElement}, true},
		{"ok", fields{ldb, odb}, args{testPath + "/sp1", okElement}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}
			if err := s.list(tt.args.path, tt.args.elem); (err != nil) != tt.wantErr {
				t.Errorf("service.list() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_get(t *testing.T) {
	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})

	mapValue := make(map[string]interface{}, 0)
	txn, _ := odb.NewTxn(true)
	json.Unmarshal([]byte(jsonValue), &mapValue)
	odb.Create(txn, testPath, mapValue)
	ldb.WriteObject([]byte(testPath), mapValue)
	odb.Commit(txn)

	defer func() {
		txn, _ := odb.NewTxn(true)
		odb.Delete(txn, testPath)
		ldb.Delete([]byte(testPath))
		odb.Commit(txn)
		ldb.Close()
	}()
	p := &p{}
	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
	}
	type args struct {
		path   string
		output interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null path", fields{ldb, odb}, args{"", p}, true},
		{"invalid path", fields{ldb, odb}, args{"/i/am/invalid/path", p}, true},
		{"null element", fields{ldb, odb}, args{testPath, nil}, true},
		{"element not match", fields{ldb, odb}, args{testPath + "/sp1", p}, false},
		{"ok", fields{ldb, odb}, args{testPath, p}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}
			if err := s.get(tt.args.path, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("service.get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_create(t *testing.T) {
	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})

	defer func() {
		txn, _ := odb.NewTxn(true)
		odb.Delete(txn, testPath)
		ldb.Delete([]byte(testPath))
		odb.Commit(txn)
		ldb.Close()
	}()
	po := &p{Sp1: sp1{F1: "fld1"}, Sp2: sp2{F2: "fld2"}}

	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
	}
	type args struct {
		path   string
		object interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null path", fields{ldb, odb}, args{"", po}, true},
		{"invalid path", fields{ldb, odb}, args{"/i/am/invalid/path", po}, true},
		{"ok", fields{ldb, odb}, args{testPath, po}, false},
		{"null element", fields{ldb, odb}, args{testPath, nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}

			if err := s.create(tt.args.path, tt.args.object); (err != nil) != tt.wantErr {
				t.Errorf("service.create() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func Test_service_update(t *testing.T) {
	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})

	mapValue := make(map[string]interface{}, 0)
	txn, _ := odb.NewTxn(true)
	json.Unmarshal([]byte(jsonValue), &mapValue)
	odb.Create(txn, testPath, mapValue)
	ldb.WriteObject([]byte(testPath), mapValue)
	odb.Commit(txn)

	defer func() {
		txn, _ := odb.NewTxn(true)
		odb.Delete(txn, testPath)
		ldb.Delete([]byte(testPath))
		odb.Commit(txn)
		ldb.Close()
	}()
	po := p{Sp1: sp1{F1: "fld1"}, Sp2: sp2{F2: "fld2"}}
	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
	}
	type args struct {
		path   string
		object interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null path", fields{ldb, odb}, args{"", po}, true},
		{"invalid path", fields{ldb, odb}, args{"/i/am/invalid/path", po}, true},
		{"ok", fields{ldb, odb}, args{testPath, po}, false},
		{"null element", fields{ldb, odb}, args{testPath, nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}

			if err := s.update(tt.args.path, tt.args.object); (err != nil) != tt.wantErr {
				t.Errorf("service.update() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func Test_service_delete(t *testing.T) {
	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})

	mapValue := make(map[string]interface{}, 0)
	txn, _ := odb.NewTxn(true)
	json.Unmarshal([]byte(jsonValue), &mapValue)
	odb.Create(txn, testPath, mapValue)
	ldb.WriteObject([]byte(testPath), mapValue)
	odb.Commit(txn)

	defer func() {
		txn, _ := odb.NewTxn(true)
		odb.Delete(txn, testPath)
		ldb.Delete([]byte(testPath))
		odb.Commit(txn)
		ldb.Close()
	}()
	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Create test cases.
		{"null path", fields{ldb, odb}, args{""}, true},
		{"invalid path", fields{ldb, odb}, args{"/i/am/invalid/path"}, true},
		{"ok", fields{ldb, odb}, args{testPath}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}
			if err := s.delete(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("service.delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_read(t *testing.T) {
	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})

	mapValue := make(map[string]interface{}, 0)
	txn, _ := odb.NewTxn(true)
	json.Unmarshal([]byte(jsonValue), &mapValue)
	odb.Create(txn, testPath, mapValue)
	ldb.WriteObject([]byte(testPath), mapValue)
	odb.Commit(txn)

	defer func() {
		txn, _ := odb.NewTxn(true)
		odb.Delete(txn, testPath)
		ldb.Delete([]byte(testPath))
		odb.Commit(txn)
		ldb.Close()
	}()

	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
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
		{"null path", fields{ldb, odb}, args{""}, nil, true},
		{"invalid path", fields{ldb, odb}, args{"/i/am/invalid/path"}, nil, true},
		{"ok", fields{ldb, odb}, args{testPath}, mapValue, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}
			got, err := s.read(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.read() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_service_write(t *testing.T) {
// 	var odb, _ = opa.New(&opa.Options{RegoPath: regoPath, DataPath: dataPath})
// 	var ldb, _ = leveldb.New(&leveldb.Options{DBPath: dbPath})

// 	defer func() {
// 		txn, _ := odb.NewTxn(true)
// 		odb.Delete(txn, testPath)
// 		ldb.Delete([]byte(testPath))
// 		odb.Commit(txn)
// 		ldb.Close()
// 	}()
// 	po := &p{Sp1: sp1{F1: "fld1"}, Sp2: sp2{F2: "fld2"}}

// 	type fields struct {
// 		ldb leveldb.LeveldbStore
// 		odb opa.OpaStore
// 	}
// 	type args struct {
// 		path   string
// 		object interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Create test cases.
// 		{"null path", fields{ldb, odb}, args{"", po}, true},
// 		{"invalid path", fields{ldb, odb}, args{"/i/am/invalid/path", po}, true},
// 		{"ok", fields{ldb, odb}, args{testPath, po}, false},
// 		{"null element", fields{ldb, odb}, args{testPath, nil}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &service{
// 				ldb: tt.fields.ldb,
// 				odb: tt.fields.odb,
// 			}
// 			if err := s.write(tt.args.path, tt.args.object); (err != nil) != tt.wantErr {
// 				t.Errorf("service.write() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func Test_service_query(t *testing.T) {
	type fields struct {
		ldb leveldb.LeveldbStore
		odb opa.OpaStore
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				ldb: tt.fields.ldb,
				odb: tt.fields.odb,
			}
			if err := s.query(tt.args.rule, tt.args.input, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("service.query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func Test_map2struct(t *testing.T) {
// 	mapValue := make(map[string]interface{}, 0)
// 	json.Unmarshal([]byte(jsonValue), &mapValue)
// 	o := p{}
// 	type args struct {
// 		m interface{}
// 		s interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Create test cases.
// 		{"null map & null struct", args{nil, nil}, true},
// 		{"null map & ok struct", args{nil, &o}, false},
// 		{"ok map & null struct", args{&mapValue, nil}, true},
// 		{"ok map & ok struct", args{&mapValue, &o}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := map2struct(tt.args.m, tt.args.s); (err != nil) != tt.wantErr {
// 				t.Errorf("map2struct() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			t.Logf("%+v", tt.args.s)
// 		})
// 		// t.Logf("%+v", o)
// 	}
// }

// func Test_struct2map(t *testing.T) {
// 	mapValue := make(map[string]interface{}, 0)
// 	// json.Unmarshal([]byte(jsonValue), &mapValue)
// 	o := p{}
// 	type args struct {
// 		s interface{}
// 		m interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Create test cases.
// 		{"null struct & null map", args{nil, nil}, false},
// 		{"null struct & ok map", args{nil, &mapValue}, false},
// 		{"ok struct & null map", args{&o, nil}, false},
// 		{"null struct & null map", args{&o, mapValue}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := struct2map(tt.args.s, tt.args.m); (err != nil) != tt.wantErr {
// 				t.Errorf("struct2map() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			t.Logf("%+v", tt.args.m)
// 		})
// 	}
// }
