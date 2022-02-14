package request

// todo request 和 dao 进一步整理后将pg连接统一管理
import (
	"database/sql"
	_ "github.com/lib/pq"
	"importer/app/config"
	"importer/app/log"
)

var ifindPg *sql.DB

func ifindPgInit() {
	pgCfg := config.GetPgConfig()
	db, err := sql.Open("postgres", pgCfg.IfindPgDSN)
	if err != nil {
		log.Log.Fatal(err.Error())
	}
	ifindPg = db
}

func getIfindPg() *sql.DB {
	if ifindPg == nil {
		ifindPgInit()
	}
	return ifindPg
}
