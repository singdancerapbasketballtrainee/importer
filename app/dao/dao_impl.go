package dao

import (
	"fmt"
	"github.com/pkg/errors"
	"importer/app/dao/orm"
	"importer/app/log"
	"importer/app/request"
	"time"
)

// ImportRps is a dao.将rps数据批量导入数据库
func (d *dao) ImportRps(date string) error {
	// 校验表是否存在，不存在则建表
	if !d.db.Migrator().HasTable(&orm.Rps{}) {
		log.Log.Info("table rps not exist, create table")
		err := d.db.AutoMigrate(&orm.Rps{})
		if err != nil {
			return err
		}
		log.Log.Info("table rps create success")
	}
	rps, err := request.GetRps(date)
	if err != nil {
		return err
	}
	if len(rps) == 0 {
		log.Log.Warn(fmt.Sprintf("lack of rps data in %s", date))
		return errors.New(fmt.Sprintf("lack of rps data in %s", date))
	}
	d.db.Delete(orm.Rps{}, "Tradedate = ?", date)
	d.db.Create(rps)
	return d.db.Error
}

// ImportFxj is a dao.将发行价数据批量导入数据库
func (d *dao) ImportFxj() error {
	// 校验表是否存在，不存在则建表
	if !d.db.Migrator().HasTable(&orm.Fxj{}) {
		log.Log.Info("table fxj not exist, create table")
		err := d.db.AutoMigrate(&orm.Fxj{})
		if err != nil {
			return err
		}
		log.Log.Info("create table fxj success")
	}
	t := time.Now()
	fxj, err := request.GetFxj()
	if err != nil {
		return err
	}
	d.db.Create(fxj)
	d.db.Delete(orm.Fxj{}, "mtime < ?", t)
	return d.db.Error
}
