package leveldb

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
	"wailik.com/internal/pkg/log"
)

const (
	DBPath = "/opt/rbac/db"
)

func TestNew(t *testing.T) {
	type args struct {
		opt *Options
	}
	tests := []struct {
		name    string
		args    args
		want    *leveldbStore
		wantErr bool
	}{
		// TODO: Add test cases.
		{"db options is null", args{opt: nil}, nil, true},
		{"db path is null", args{opt: &Options{DBPath: ""}}, nil, true},
		{"db path is invalid", args{opt: &Options{DBPath: "/i/am/invalid/path"}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.opt)
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

// func Test_leveldbStore_Write(t *testing.T) {
// 	db, _ := leveldb.OpenFile(DBPath, nil)
// 	defer db.Close()

// 	type fields struct {
// 		db *leveldb.DB
// 	}
// 	type args struct {
// 		key   []byte
// 		value []byte
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{"db write with null db object", fields{nil}, args{[]byte("/i/am/a/key"), []byte("value")}, true},
// 		{"db write with null key", fields{db}, args{nil, []byte("value")}, true},
// 		{"db write with null value", fields{db}, args{[]byte("/i/am/a/key"), nil}, false},
// 		{"db write with null key & null value", fields{db}, args{nil, nil}, true},
// 		{"db write with empty key string", fields{db}, args{[]byte(""), []byte("value")}, true},
// 		{"db write with empty value string", fields{db}, args{[]byte("/i/am/a/key"), []byte("")}, false},
// 		{"db write with key & value", fields{db}, args{[]byte("/i/am/a/key"), []byte("value")}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &leveldbStore{
// 				db: tt.fields.db,
// 			}
// 			if err := s.Write(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
// 				t.Errorf("leveldbStore.Write() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}

// }

func Test_leveldbStore_Delete(t *testing.T) {
	db, _ := leveldb.OpenFile(DBPath, nil)
	defer db.Close()
	type fields struct {
		db *leveldb.DB
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"db write with null db object", fields{nil}, args{[]byte("/i/am/a/key")}, true},
		{"db write with null key", fields{db}, args{nil}, true},
		{"db write with empty key string", fields{db}, args{[]byte("")}, true},
		{"db write with key & value", fields{db}, args{[]byte("/i/am/a/key")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &leveldbStore{
				db: tt.fields.db,
			}
			if err := s.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("leveldbStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_leveldbStore_WriteObject(t *testing.T) {
	db, _ := leveldb.OpenFile(DBPath, nil)
	defer db.Close()
	type object struct {
		name    string
		pObject *object
	}
	type fields struct {
		db *leveldb.DB
	}
	type args struct {
		key    []byte
		object interface{}
	}

	obj := object{
		name: "parent",
		pObject: &object{
			name:    "child",
			pObject: nil,
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"db write object with null db object", fields{nil}, args{[]byte("/i/am/a/key"), obj}, true},
		{"db write object  null key", fields{db}, args{nil, obj}, true},
		{"db write object with null value", fields{db}, args{[]byte("/i/am/a/key"), nil}, false},
		{"db write object with null key & null value", fields{db}, args{nil, nil}, true},
		{"db write object with empty key string", fields{db}, args{[]byte(""), obj}, true},
		{"db write object with empty value string", fields{db}, args{[]byte("/i/am/a/key"), []byte("")}, false},
		{"db write object with key & value", fields{db}, args{[]byte("/i/am/a/key"), obj}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &leveldbStore{
				db: tt.fields.db,
			}
			if err := s.WriteObject(tt.args.key, tt.args.object); (err != nil) != tt.wantErr {
				t.Errorf("leveldbStore.WriteObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.Close()
}

func Test_leveldbStore_TraverseAll(t *testing.T) {
	db, _ := leveldb.OpenFile(DBPath, nil)
	defer db.Close()
	var f Operation = func(key []byte, value interface{}) error {
		fmt.Printf("%+v:%+v", key, value)

		return nil
	}
	type fields struct {
		db *leveldb.DB
	}
	type args struct {
		oper Operation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"db travel with null db object", fields{nil}, args{f}, true},
		{"db travel with null func", fields{db}, args{nil}, true},
		{"db travel with func", fields{db}, args{f}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &leveldbStore{
				db: tt.fields.db,
			}
			if err := s.TraverseAll(tt.args.oper); (err != nil) != tt.wantErr {
				t.Errorf("leveldbStore.TraverseAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_leveldbStore_Close(t *testing.T) {
	db, _ := leveldb.OpenFile(DBPath, nil)
	defer db.Close()
	type fields struct {
		db *leveldb.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{"close db with null db object", fields{nil}, false},
		{"close db", fields{db: db}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &leveldbStore{
				db: tt.fields.db,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("leveldbStore.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_leveldbStore_DeleteWithPrefix(t *testing.T) {
	log.Init(log.OptLevel(log.DebugLevel))
	db, _ := leveldb.OpenFile(DBPath, nil)
	defer db.Close()

	db.Put([]byte("/1"), []byte("v1"), nil)
	db.Put([]byte("/1/2"), []byte("v2"), nil)
	db.Put([]byte("/1/2/3"), []byte("v3"), nil)

	type fields struct {
		db *leveldb.DB
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"ok", fields{db: db}, args{key: []byte("/1")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &leveldbStore{
				db: tt.fields.db,
			}
			if err := s.DeleteWithPrefix(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("leveldbStore.BatchDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "ok" {
				existed, err := db.Has([]byte("/1/2"), nil)
				if existed || err != nil {
					t.Errorf("leveldbStore.BatchDelete() not work")
				}
			}
		})
	}
}
