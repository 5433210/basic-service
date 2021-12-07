package dbstore

import (
	"time"

	"gorm.io/gorm"

	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/log"
)

type StoreType string

type Transaction struct {
	Handler *gorm.DB
}

type Handler struct {
	Handler *gorm.DB
}

type Store interface {
	Open() error
	Close()
	MaxIdleConns(n int)
	MaxOpenConns(n int)
	ConnMaxLifetime(d time.Duration)
	ConnMaxIdletime(d time.Duration)
	BeginTxn() Transaction
	CommitTxn(tx Transaction) Transaction
	RollbackTxn(tx Transaction) Transaction
	DB() Transaction
	Transaction(tx *Transaction) Transaction
}

type storeImpl interface {
	Open() (*gorm.DB, error)
	Close()
	DB() *gorm.DB
}

type store struct {
	impl storeImpl
}

type Options struct {
	DSN string
}

var _ Store = &store{}

func New(opts *Options) (Store, error) {
	storeType, err := getStoreTypeByDSN(opts.DSN)
	if err != nil {
		return nil, err
	}
	switch storeType {
	case constant.DsMysql:
		return &store{impl: newDbMysql(opts.DSN)}, nil
	default:
		return nil, err
	}
}

func getStoreTypeByDSN(dsn string) (StoreType, error) {
	return constant.DsMysql, nil
}

func (s *store) Open() error {
	_, err := s.impl.Open()

	return err
}

func (s *store) Close() {
	s.impl.Close()
}

func (s *store) MaxIdleConns(n int) {
	db, err := s.impl.DB().DB()
	if err != nil {
		db.SetMaxIdleConns(n)
	}
}

func (s *store) MaxOpenConns(n int) {
	db, err := s.impl.DB().DB()
	if err != nil {
		db.SetMaxOpenConns(n)
	}
}

func (s *store) ConnMaxLifetime(d time.Duration) {
	db, err := s.impl.DB().DB()
	if err != nil {
		db.SetConnMaxLifetime(d)
	}
}

func (s *store) ConnMaxIdletime(d time.Duration) {
	db, err := s.impl.DB().DB()
	if err != nil {
		db.SetConnMaxIdleTime(d)
	}
}

func (s *store) BeginTxn() Transaction {
	log.Debug("begin tTxnransaction...")

	return Transaction{Handler: s.impl.DB().Begin()}
}

func (s *store) CommitTxn(tx Transaction) Transaction {
	log.Debug("commit transaction")

	return Transaction{Handler: tx.Handler.Commit()}
}

func (s *store) RollbackTxn(tx Transaction) Transaction {
	log.Debug("rollback transaction")

	return Transaction{Handler: tx.Handler.Rollback()}
}

func (s *store) DB() Transaction {
	log.Debug("onetime operation")

	return Transaction{Handler: s.impl.DB()}
}

func (s *store) Transaction(tx *Transaction) Transaction {
	if tx == nil {
		return s.DB()
	}

	return *tx
}
