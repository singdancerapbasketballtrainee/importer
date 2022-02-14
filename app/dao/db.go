package dao

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"importer/app/config"
)

func NewDB() (db *gorm.DB, cf func(), err error) {
	pgCfg := config.GetPgConfig()
	db, err = makePgConn(pgCfg)
	cf = Close
	return
}

func Close() {
}

/**
 * @Description: 通过dsn创建pg连接对象
 * @param dsn
 * @return db
 * @return err
 */
func makePgConn(pgCfg config.PgConfig) (db *gorm.DB, err error) {
	level := logger.Warn
	db, err = gorm.Open(postgres.Open(pgCfg.DefaultDSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		return
	}

	// 设置请求超时时间
	db.Exec(fmt.Sprintf("set statement_timeout to %d", pgCfg.QueryTimeout))
	// 连接设置初始化
	d, err := db.DB()
	if err != nil {
		return
	}
	d.SetMaxIdleConns(pgCfg.MaxIdleConns)
	d.SetMaxOpenConns(pgCfg.MaxOpenConns)
	return
}
