package dao

import (
	"fmt"
	"gorm.io/gorm"
	"importer/app/config"
	"importer/app/cron"
	"importer/app/log"
	"time"
)
import "github.com/google/wire"

type Dao interface {
	Start() error
	Close()
	ImportRps(date string) error
	ImportFxj() error
}

type dao struct {
	db *gorm.DB
}

var Provider = wire.NewSet(New, NewDB)

func New(db *gorm.DB) (d Dao, cf func(), err error) {
	return newDao(db)
}

func newDao(db *gorm.DB) (d *dao, cf func(), err error) {
	d = &dao{
		db,
	}
	cf = d.Close
	return
}

func (d *dao) Start() (err error) {
	return d.addImportCron()
}

func (d *dao) Close() {

}

func (d *dao) addImportCron() error {
	apiCfg := config.GetApiConfig()

	err := d.addImportRpsCron(apiCfg.RpsCfg.Cron)
	if err != nil {
		return err
	}
	err = d.addImportFxjCron(apiCfg.FxjCfg.Cron)
	return err
}

func (d *dao) addImportRpsCron(spec string) error {
	rpsCron := func() {
		log.Log.Info("rps cron start")
		t := time.Now()
		year, mon, day := t.Date()
		err := d.ImportRps(fmt.Sprintf("%4d-%02d-%02d", year, mon, day))
		if err != nil {
			log.Log.Error("import rps error:" + err.Error())
		}
	}
	return cron.AddFunc(spec, rpsCron)
}

func (d *dao) addImportFxjCron(spec string) error {
	fxjCron := func() {
		log.Log.Info("fxj cron start")
		err := d.ImportFxj()
		if err != nil {
			log.Log.Error("import fxj error:" + err.Error())
		}
	}
	return cron.AddFunc(spec, fxjCron)
}
