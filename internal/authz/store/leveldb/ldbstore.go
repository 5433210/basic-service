package leveldb

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type LeveldbStore interface {
	// Write(key []byte, value []byte) error
	WriteObject(key []byte, object interface{}) error
	Delete(key []byte) error
	DeleteWithPrefix(prefix []byte) error
	TraverseAll(oper Operation) error
	TraverseWithPrefix(prefix []byte, oper Operation) error
	Close() error
}

type leveldbStore struct {
	db *leveldb.DB
}

type Options struct {
	DBPath string
}

var _ LeveldbStore = &leveldbStore{}

func New(opts *Options) (*leveldbStore, error) {
	if opts == nil {
		err := errors.NewErrorC(errors.ErrCdOptionIsNull, nil)
		log.ErrorLog(err)

		return nil, err
	}

	ldb, err := leveldb.OpenFile(opts.DBPath, nil)
	if err != nil {
		log.ErrorLog(err)

		return nil, err
	}

	log.Infof("leveldb(%v) opened", &ldb)

	return &leveldbStore{
		db: ldb,
	}, nil
}

type Operation func(key []byte, value interface{}) error

func (s *leveldbStore) TraverseAll(oper Operation) error {
	log.Debugf("ldb traverse all key...")

	if s.db == nil {
		err := errors.NewErrorC(errors.ErrCdNeedInit, nil)
		log.ErrorLog(err)

		return err
	}
	if oper == nil {
		err := errors.NewErrorC(errors.ErrCdFuncIsNull, nil)
		log.ErrorLog(err)

		return err
	}

	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		var v interface{}
		err := json.Unmarshal(value, &v)
		if err != nil {
			// log.ErrorLog(err, log.ByteString("key", key), log.ByteString("value", value))
			err = oper(key, value)
			if err != nil {
				return err
			}
		} else {
			err = oper(key, v)
			if err != nil {
				return err
			}
		}
	}

	log.Debugf("ldb all keys traversed")

	return iter.Error()
}

func (s *leveldbStore) TraverseWithPrefix(prefix []byte, oper Operation) error {
	log.Debugf("ldb traverse keys with prefix:%v ...", string(prefix))
	if s.db == nil {
		err := errors.NewErrorC(errors.ErrCdNeedInit, nil)
		log.ErrorLog(err)

		return err
	}
	if oper == nil {
		err := errors.NewErrorC(errors.ErrCdFuncIsNull, nil)
		log.ErrorLog(err)

		return err
	}

	iter := s.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		log.Debugf("k:%+v", string(key))
		log.Debugf("v:%+v", string(value))
		var v interface{}
		err := json.Unmarshal(value, &v)
		if err != nil {
			// log.ErrorLog(err, log.ByteString("key", key), log.ByteString("value", value))
			err = oper(key, value)
			if err != nil {
				return err
			}
		} else {
			err = oper(key, v)
			if err != nil {
				return err
			}
		}
	}
	log.Debugf("ldb keys with prefix:%v traversed", string(prefix))

	return iter.Error()
}

// func (s *leveldbStore) Write(key []byte, value []byte) error {
// 	log.Debugf("ldb write key:%v value:%v", string(key), string(value))
// 	if s.db == nil {
// 		err := errors.NewErrorC(errors.ErrCdNeedInit, nil)
// 		log.ErrorLog(err)

// 		return err
// 	}
// 	if len(key) == 0 {
// 		err := errors.NewErrorC(errors.ErrCdKeyIsNull, nil)
// 		log.ErrorLog(err)

// 		return err
// 	}
// 	err := s.db.Put(key, value, nil)
// 	if err != nil {
// 		log.ErrorLog(err)

// 		return err
// 	}

// 	log.Debugf("ldb key:%v writed", string(key))

// 	return nil
// }

func (s *leveldbStore) Delete(key []byte) error {
	log.Debugf("ldb delete key:%v", string(key))
	if s.db == nil {
		err := errors.NewErrorC(errors.ErrCdNeedInit, nil)
		log.ErrorLog(err)

		return err
	}
	if len(key) == 0 {
		err := errors.NewErrorC(errors.ErrCdKeyIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	if err := s.db.Delete(key, nil); err != nil {
		log.ErrorLog(err)

		return err
	}
	log.Debugf("ldb key:%v deleted", string(key))

	return nil
}

func (s *leveldbStore) DeleteWithPrefix(prefix []byte) error {
	log.Debugf("ldb delete keys with prefix:%v", string(prefix))
	if s.db == nil {
		err := errors.NewErrorC(errors.ErrCdNeedInit, nil)
		log.ErrorLog(err)

		return err
	}
	if len(prefix) == 0 {
		err := errors.NewErrorC(errors.ErrCdKeyIsNull, nil)
		log.ErrorLog(err)

		return err
	}
	batch := new(leveldb.Batch)
	iter := s.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		batch.Delete(key)
		log.Debugf("delete key:%+v", string(key))
	}
	if err := iter.Error(); err != nil {
		return err
	}
	if err := s.db.Write(batch, nil); err != nil {
		return err
	}
	log.Debugf("ldb keys with prefix:%v deleted", string(prefix))

	return nil
}

func (s *leveldbStore) WriteObject(key []byte, object interface{}) error {
	if s.db == nil {
		err := errors.NewErrorC(errors.ErrCdNeedInit, nil)
		log.ErrorLog(err)

		return err
	}
	if len(key) == 0 {
		err := errors.NewErrorC(errors.ErrCdKeyIsNull, nil)
		log.ErrorLog(err)

		return err
	}

	value, err := json.Marshal(object)
	if err != nil {
		log.ErrorLog(err)

		return err
	}
	err = s.db.Put(key, value, nil)
	if err != nil {
		log.ErrorLog(err)

		return err
	}

	return nil
}

func (s *leveldbStore) Close() error {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			return err
		}

		log.Infof("leveldb(%v) closed", &s.db)

		s.db = nil
	}

	return nil
}
