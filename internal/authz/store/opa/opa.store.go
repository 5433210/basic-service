package opa

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"

	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type OpaStore interface {
	NewTxn(writable bool) (OpaTxn, error)
	Commit(txn OpaTxn) error
	Abort(txn OpaTxn)
	Create(txn OpaTxn, path string, value interface{}) error
	Delete(txn OpaTxn, path string) error
	Update(txn OpaTxn, path string, o interface{}) error
	Read(txn OpaTxn, path string) (interface{}, error)
	Query(rule string, input interface{}, output interface{}) error
	Exist(txn OpaTxn, path string) bool
	DeleteOne(path string) error
}

type opaStore struct {
	mem storage.Store
	cmp *ast.Compiler
}

type OpaTxn interface {
	Ctx() context.Context
	Txn() storage.Transaction
}

type opaTxn struct {
	ctx context.Context
	txn storage.Transaction
}

func (txn opaTxn) Ctx() context.Context {
	return txn.ctx
}

func (txn opaTxn) Txn() storage.Transaction {
	return txn.txn
}

type Options struct {
	RegoPath string
	DataPath string
}

func New(opts *Options) (*opaStore, error) {
	if opts == nil || opts.DataPath == "" || opts.RegoPath == "" {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return nil, err
	}
	s := &opaStore{
		mem: inmem.New(),
		cmp: &ast.Compiler{},
	}

	err := s.Load(opts.RegoPath, opts.DataPath)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *opaStore) Load(regoPath string, dataPath string) error {
	log.Debugf("opa store load data...")
	err := s.loadRego(regoPath)
	if err != nil {
		return err
	}

	err = s.loadData(dataPath)
	if err != nil {
		return err
	}

	log.Debugf("opa store data loaded")

	return nil
}

func (s *opaStore) NewTxn(writable bool) (OpaTxn, error) {
	if s.mem == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return nil, err
	}
	ctx := context.Background()
	param := storage.TransactionParams{Write: writable}
	txn, err := s.mem.NewTransaction(ctx, param)
	if err != nil {
		log.ErrorLog(err)

		return nil, err
	}

	return &opaTxn{
		txn: txn,
		ctx: ctx,
	}, nil
}

func (s *opaStore) Commit(txn OpaTxn) error {
	if s.mem == nil || txn == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	err := s.mem.Commit(txn.Ctx(), txn.Txn())
	if err != nil {
		log.ErrorLog(err)
	}

	return err
}

func (s *opaStore) Abort(txn OpaTxn) {
	if s.mem == nil || txn == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return
	}
	s.mem.Abort(txn.Ctx(), txn.Txn())
}

func (s *opaStore) Create(txn OpaTxn, path string, value interface{}) error {
	log.Debugf("opa create path:%+v, value:%+v", path, value)
	err := s.write(txn, path, value)
	if err != nil {
		return err
	}

	log.Debugf("opa path:%v value:%+v created", path, value)

	return nil
}

func (s *opaStore) Delete(txn OpaTxn, path string) error {
	log.Debugf("opa delete path:%+v", path)
	if s.mem == nil || txn == nil || len(path) == 0 {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	parsedPath, ok := storage.ParsePathEscaped(path)
	if !ok {
		err := errors.NewErrorCf(errors.ErrCdPathParse, nil, "path(%v)", path)
		log.ErrorLog(err)

		return err
	}
	err := s.mem.Write(txn.Ctx(), txn.Txn(), storage.RemoveOp, parsedPath, nil)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	log.Debugf("opa path:%v deleted", path)

	return nil
}

func (s *opaStore) Read(txn OpaTxn, path string) (interface{}, error) {
	log.Debugf("opa read path:%+v", path)
	rs, err := s.read(txn, path)
	if err != nil {
		return nil, err
	}
	log.Debugf("opa path:%v read", path)

	return rs, nil
}

func (s *opaStore) Query(rule string, input interface{}, output interface{}) error {
	if s.mem == nil || s.cmp == nil || len(rule) == 0 || output == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	log.Debugf("opa query rule:%v, input:%+v", rule, input)
	ctx := context.Background()

	rego := rego.New(
		rego.Query(rule),
		rego.Compiler(s.cmp),
		rego.Store(s.mem),
		rego.Input(input),
	)
	rs, err := rego.Eval(ctx)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		err := errors.NewErrorC(errors.ErrCdDataNotFound, nil)
		log.ErrorLog(err)

		return err
	}
	j, err := json.Marshal(rs[0].Expressions[0].Value)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	err = json.Unmarshal(j, output)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	log.Debugf("opa rule:%v, input:%+v query, got output:%+v", rule, input, output)

	return nil
}

func (s *opaStore) Update(txn OpaTxn, path string, value interface{}) error {
	log.Debugf("opa update path:%+v, value:%+v", path, value)
	if s.mem == nil || txn == nil || len(path) == 0 || value == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Struct {
		for k := 0; k < t.NumField(); k++ {
			tag := strings.Split(t.Field(k).Tag.Get("json"), ",")[0]
			// log.Debugf("tag name:%+v", tag)
			tagV := v.Field(k).Interface()
			// log.Debugf("tag value:%+v", tagV)
			subpath := path + "/" + tag
			if err := s.replaceOrCreate(txn, subpath, tagV); err != nil {
				return err
			}
		}
	} else {
		if err := s.replaceOrCreate(txn, path, v); err != nil {
			return err
		}
	}
	log.Debugf("opa path:%+v, value:%+v updated", path, value)

	return nil
}

func (s *opaStore) Exist(txn OpaTxn, path string) bool {
	if rs := s.exist(txn, path); rs {
		log.Debugf("path(%v) existed in opa store", path)

		return true
	}

	log.Debugf("path(%v) not existed in opa store", path)

	return false
}

func (s *opaStore) ReadOne(path string) (interface{}, error) {
	if s.mem == nil || len(path) == 0 {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return nil, err
	}
	txn, err := s.NewTxn(false)
	if err != nil {
		return nil, err
	}

	r, err := s.Read(txn, path)
	if err != nil {
		s.Abort(txn)

		return nil, err
	}

	err = s.Commit(txn)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *opaStore) CreateOne(path string, value interface{}) error {
	if s.mem == nil || len(path) == 0 || value == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	txn, err := s.NewTxn(true)
	if err != nil {
		return err
	}

	err = s.Create(txn, path, value)
	if err != nil {
		s.Abort(txn)

		return err
	}

	err = s.Commit(txn)
	if err != nil {
		return err
	}

	return nil
}

func (s *opaStore) DeleteOne(path string) error {
	if s.mem == nil || len(path) == 0 {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	txn, err := s.NewTxn(true)
	if err != nil {
		return err
	}

	err = s.Delete(txn, path)
	if err != nil {
		s.Abort(txn)

		return err
	}

	err = s.Commit(txn)
	if err != nil {
		return err
	}

	return nil
}

func (s *opaStore) write(txn OpaTxn, path string, value interface{}) error {
	if s.mem == nil || txn == nil || len(path) == 0 || value == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	parsedPath, ok := storage.ParsePathEscaped(path)
	if !ok {
		err := errors.NewErrorCf(errors.ErrCdPathParse, nil, "path(%v)", path)
		log.ErrorLog(err)

		return err
	}

	// log.Debugf("parsedPath:%+v", parsedPath)

	err := s.mem.Write(txn.Ctx(), txn.Txn(), storage.AddOp, parsedPath, value)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	return nil
}

func (s *opaStore) read(txn OpaTxn, path string) (interface{}, error) {
	if s.mem == nil || txn == nil || len(path) == 0 {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return nil, err
	}
	parsedPath, ok := storage.ParsePathEscaped(path)
	if !ok {
		err := errors.NewErrorCf(errors.ErrCdPathParse, nil, "path(%v)", path)
		log.ErrorLog(err)

		return nil, err
	}

	mem := s.mem
	r, err := mem.Read(txn.Ctx(), txn.Txn(), parsedPath)
	if err != nil {
		if storage.IsNotFound(err) {
			err = errors.NewErrorCf(errors.ErrCdDataNotExist, err, "path(%v)", path)
		} else {
			log.ErrorLog(err)
		}

		return nil, err
	}

	return r, nil
}

func (s *opaStore) replace(txn OpaTxn, path string, value interface{}) error {
	if s.mem == nil || txn == nil || len(path) == 0 {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}

	parsedPath, ok := storage.ParsePathEscaped(path)
	if !ok {
		err := errors.NewErrorCf(errors.ErrCdPathParse, nil, "path(%v)", path)
		log.ErrorLog(err)

		return err
	}
	err := s.mem.Write(txn.Ctx(), txn.Txn(), storage.ReplaceOp, parsedPath, value)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	return nil
}

func (s *opaStore) replaceOrCreate(txn OpaTxn, path string, value interface{}) error {
	if s.exist(txn, path) {
		err := s.replace(txn, path, value)
		if err != nil {
			return err
		}
	} else {
		err := s.write(txn, path, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *opaStore) exist(txn OpaTxn, path string) bool {
	if _, err := s.read(txn, path); err != nil {
		return false
	}

	return true
}

func (s *opaStore) loadRego(regoPath string) error {
	if s.mem == nil || regoPath == "" {
		err := errors.NewErrorCf(errors.ErrCdArgIsNull, nil, "store:%+v, regoPath:%v", s.mem, regoPath)

		log.ErrorLog(err)

		return err
	}
	file, err := os.Open(regoPath)
	if err != nil {
		return err
	}
	defer file.Close()

	regoContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	s.cmp, err = ast.CompileModules(map[string]string{
		constant.RegoFilename: string(regoContent),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *opaStore) loadData(dataPath string) error {
	if s.mem == nil || dataPath == "" {
		err := errors.NewErrorCf(errors.ErrCdArgIsNull, nil, "store:%+v, dataPath:%v", s.mem, dataPath)
		log.ErrorLog(err)

		return err
	}
	file, err := os.Open(dataPath)
	if err != nil {
		return err
	}
	defer file.Close()
	s.mem = inmem.NewFromReader(file)

	return nil
}
