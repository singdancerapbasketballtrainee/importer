package request

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"importer/app/config"
	"importer/app/dao/orm"
	"io/ioutil"
	"net/http"
	"time"
)

// RpsRet rps 数据返回json格式对应该结构体
type RpsRet struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	Data       struct {
		StockList []struct {
			StockCode string `json:"stock_code"`
			RpsMin    int    `json:"rps_min"`
			RpsMax    int    `json:"rps_max"`
		} `json:"stock_list"`
	} `json:"data"`
}

// GetRps 提过接口获取rps数据，并返回orm.Rps切片形式
func GetRps(date string) (rps []orm.Rps, err error) {
	url := fmt.Sprintf("%s?date=%s", config.GetApiConfig().RpsCfg.Url, date)
	reps, err := http.Get(url)
	if err != nil {
		return
	}
	defer reps.Body.Close()
	body, err := ioutil.ReadAll(reps.Body)
	if err != nil {
		return
	}
	var r RpsRet
	err = json.Unmarshal(body, &r)
	if r.StatusCode != 0 {
		return nil, errors.New(r.StatusMsg)
	}
	rps = make([]orm.Rps, len(r.Data.StockList))
	t := time.Now()
	for i, value := range r.Data.StockList {
		rps[i] = orm.Rps{
			StockCode: value.StockCode,
			TradeDate: date,
			Mtime:     t,
			RpsMin:    value.RpsMin,
			RpsMax:    value.RpsMax,
		}
	}
	return
}
