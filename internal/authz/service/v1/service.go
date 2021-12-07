package servicev1

import (
	"bytes"
	stdjson "encoding/json"

	jsoniter "github.com/json-iterator/go"

	"wailik.com/internal/authz/store/leveldb"
	"wailik.com/internal/authz/store/opa"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Service interface {
	Domain() *domainSrvc
	Permission() *permissionSrvc
	Group() *groupSrvc
	Role() *roleSrvc
	Deny() *denySrvc
	Subject() *subjectSrvc
	LoadData() error
	Release()
	Dump(path string) ([]byte, error)
}

type service struct {
	ldb leveldb.LeveldbStore
	odb opa.OpaStore
}

// input structs for opa rule query.
type DomainGroupInput struct {
	DomainId     string `json:"domain_id"`
	GroupPattern string `json:"group_pattern"`
}

type DomainInput struct {
	DomainId string `json:"domain_id"`
}

type DomainRoleInput struct {
	DomainId string `json:"domain_id"`
	RoleId   string `json:"role_id"`
}

type DomainDenyInput struct {
	DomainId string `json:"domain_id"`
	DenyId   string `json:"deny_id"`
}

type DomainSubjectInput struct {
	DomainId  string `json:"domain_id"`
	SubjectId string `json:"subject_id"`
}

var _ Service = &service{}

func (s *service) Domain() *domainSrvc         { return newDomainSrvc(s) }
func (s *service) Permission() *permissionSrvc { return newPermissionSrvc(s) }
func (s *service) Group() *groupSrvc           { return newGroupSrvc(s) }
func (s *service) Role() *roleSrvc             { return newRoleSrvc(s) }
func (s *service) Deny() *denySrvc             { return newDenySrvc(s) }
func (s *service) Subject() *subjectSrvc       { return newSubjectSrvc(s) }

func NewService(dbPath string, regoPath string, dataPath string) (Service, error) {
	log.Info("new service...")
	log.Infof("db path:%v", dbPath)
	log.Infof("rego path:%v", regoPath)
	log.Infof("data path:%v", dataPath)

	var leveldbStore leveldb.LeveldbStore
	leveldbOpts := &leveldb.Options{DBPath: dbPath}
	leveldbStore, err := leveldb.New(leveldbOpts)
	if err != nil {
		return nil, err
	}

	var opaStore opa.OpaStore
	opaOpts := &opa.Options{
		DataPath: dataPath,
		RegoPath: regoPath,
	}
	opaStore, err = opa.New(opaOpts)
	if err != nil {
		return nil, err
	}

	var s Service = &service{
		ldb: leveldbStore,
		odb: opaStore,
	}

	log.Info("service created")

	return s, nil
}

func (s *service) Dump(path string) ([]byte, error) {
	rs, err := s.read(path)
	if err != nil {
		return nil, err
	}

	j, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}

	var str bytes.Buffer
	err = stdjson.Indent(&str, j, "", "\t")
	if err != nil {
		return nil, err
	}

	return str.Bytes(), nil
}

func (s *service) LoadData() error {
	if s.ldb == nil || s.odb == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	action := func(txn opa.OpaTxn) (interface{}, error) {
		opAdd := func(path []byte, value interface{}) error {
			return s.odb.Create(txn, string(path), value)
		}

		err := s.ldb.TraverseWithPrefix([]byte("/"+constant.FldDomains), opAdd)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := s.doActionWithTxn(true, action)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Release() {
	s.ldb.Close()
}

type Action func(txn opa.OpaTxn) (interface{}, error)

func (s *service) doActionWithTxn(write bool, action Action) (interface{}, error) {
	if action == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return nil, err
	}

	txn, err := s.odb.NewTxn(write)
	if err != nil {
		return nil, err
	}

	output, err := action(txn)
	if err != nil {
		s.odb.Abort(txn)

		return nil, err
	}

	err = s.odb.Commit(txn)
	if err != nil {
		return nil, err
	}

	return output, nil
}

type element func(id string, jsonStr []byte) error

func (s *service) list(path string, elem element) error {
	if elem == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	m, err := s.read(path)
	if err != nil {
		return err
	}

	for key, value := range m.(map[string]interface{}) {
		j, err := json.Marshal(value)
		if err != nil {
			return err
		}
		err = elem(key, j)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) get(path string, output interface{}) error {
	if output == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	m, err := s.read(path)
	if err != nil {
		return err
	}
	err = transform(m, &output)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) create(path string, object interface{}) error {
	action := func(txn opa.OpaTxn) (interface{}, error) {
		if err := s.odb.Create(txn, path, object); err != nil {
			return nil, err
		}

		if err := s.ldb.WriteObject([]byte(path), object); err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := s.doActionWithTxn(true, action)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) delete(path string) error {
	action := func(txn opa.OpaTxn) (interface{}, error) {
		if err := s.odb.Delete(txn, path); err != nil {
			return nil, err
		}

		if err := s.ldb.Delete([]byte(path)); err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := s.doActionWithTxn(true, action)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) deleteWithPrefix(pathPrefix string) error {
	action := func(txn opa.OpaTxn) (interface{}, error) {
		if err := s.odb.Delete(txn, pathPrefix); err != nil {
			return nil, err
		}

		if err := s.ldb.DeleteWithPrefix([]byte(pathPrefix)); err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := s.doActionWithTxn(true, action)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) read(path string) (interface{}, error) {
	action := func(txn opa.OpaTxn) (interface{}, error) {
		return s.odb.Read(txn, path)
	}
	r, err := s.doActionWithTxn(false, action)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) update(path string, object interface{}) error {
	action := func(txn opa.OpaTxn) (interface{}, error) {
		if err := s.odb.Update(txn, path, object); err != nil {
			return nil, err
		}

		if err := s.ldb.WriteObject([]byte(path), object); err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := s.doActionWithTxn(true, action)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) query(rule string, input interface{}, output interface{}) error {
	if output == nil {
		err := errors.NewErrorC(errors.ErrCdArgIsNull, nil)
		log.ErrorLog(err)

		return err
	}

	mi := make(map[string]interface{})
	if input != nil {
		err := struct2map(input, &mi)
		if err != nil {
			return err
		}
	}

	err := s.odb.Query(rule, mi, output)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) exist(path string) bool {
	if _, err := s.read(path); err != nil {
		return false
	}

	return true
}

// func map2struct(src *map[string]interface{}, dst interface{}) error {
// 	return transform(src, dst)
// }

func struct2map(src interface{}, dst *map[string]interface{}) error {
	return transform(src, dst)
}

func transform(src interface{}, dst interface{}) error {
	raw, err := json.Marshal(src)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	if err = json.Unmarshal(raw, dst); err != nil {
		log.ErrorLog(err)

		return err
	}

	return nil
}

func BuildPath(fields ...string) string {
	var path string
	for _, v := range fields {
		path += constant.FldSlash
		path += v
	}

	return path
}
