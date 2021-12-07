package dbstore

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbMysql struct {
	dsn string
	db  *gorm.DB
}

func newDbMysql(dsn string) storeImpl {
	return &dbMysql{dsn: dsn}
}

func (impl *dbMysql) Open() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(impl.dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, err
	}
	impl.db = db

	return db, nil
}

func (impl *dbMysql) Close() {
}

func (impl *dbMysql) DB() *gorm.DB {
	return impl.db
}
