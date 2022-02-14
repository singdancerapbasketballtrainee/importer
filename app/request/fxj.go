package request

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"importer/app/config"
	"importer/app/dao/orm"
	"importer/app/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type XMLRequest struct {
	XMLName  xml.Name `xml:"request"`
	Items    []Item   `xml:"items>item"`
	Thscodes string   `xml:"thscodes"`
}

type Item struct {
	XMLName xml.Name `xml:"item"`
	Name    string   `xml:"name,attr"`
	Params  []Param  `xml:"params>param"`
}

type Param struct {
	XMLName xml.Name `xml:"param"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
	System  string   `xml:"system,attr"`
}

var indexXml, fundXml, bondXml XMLRequest

type _type int

// type
const (
	index _type = iota
	fund
	bond
)

func init() {
	baseParams := []Param{
		{Name: "FD0", Value: "111", System: "true"},
	}

	//IndexXmlStr = `<?xml version="1.0" encoding="UTF-8"?>
	//<request>
	//<items>
	//<item name="06971_000_00_0_1">
	//<params>
	//<param name="FD0" value="111" system="true"/>
	//</params>
	//</item>
	//</items>
	//<thscodes>code1,code2...</thscodes>
	//</request>`
	indexXml = XMLRequest{}
	item := Item{Name: "06971_000_00_0_1", Params: baseParams}
	indexXml.Items = append(indexXml.Items, item)

	//FundXmlStr = `<?xml version="1.0" encoding="UTF-8"?>
	//<request>
	//<items>
	//<item name="00038_000_00_0_5">
	//<params>
	//<param name="FD0" value="111" system="true"/>
	//<param name="FT" value="100" system="true"/>
	//</params>
	//</item>
	//</items>
	//<thscodes>code1,code2,...</thscodes>
	//</request>`
	fundXml = XMLRequest{}
	item = Item{Name: "00038_000_00_0_5", Params: baseParams}
	item.Params = append(item.Params, Param{Name: "FT", Value: "100", System: "true"})
	fundXml.Items = append(fundXml.Items, item)

	//BondXmlStr = `<?xml version="1.0" encoding="UTF-8"?>
	//<request>
	//<items>
	//<item name="03548_000_00_0_11">
	//<params>
	//<param name="FD0" value="111" system="true"/>
	//<param name="FJ" value="103" system="true"/>
	//</params>
	//</item>
	//</items>
	//<thscodes>code1,code2,...</thscodes>
	//</request>`
	bondXml = XMLRequest{}
	item = Item{Name: "03548_000_00_0_11", Params: baseParams}
	item.Params = append(item.Params, Param{Name: "FJ", Value: "103", System: "true"})
	bondXml.Items = append(bondXml.Items, item)
}

// GetFxj 发行价的接口设计写得比较死，三支代码只能一个个来
func GetFxj() ([]orm.Fxj, error) {
	fxj := make([]orm.Fxj, 0)
	apiCfg := config.GetApiConfig()
	req, err := http.NewRequest("POST", apiCfg.FxjCfg.Url, nil)
	if err != nil {
		return fxj, err
	}
	// 设置http头
	req.Header.Set("X-Arsenal-Auth", apiCfg.AppName)

	stockMarkets := map[_type][]uint8{
		index: apiCfg.FxjCfg.Markets.Index,
		fund:  apiCfg.FxjCfg.Markets.Fund,
		bond:  apiCfg.FxjCfg.Markets.Bond,
	}

	for stockType, markets := range stockMarkets {
		for _, market := range markets {
			f, err := getFxjData(req, market, stockType)
			if err != nil {
				log.Log.Error(fmt.Sprintf("get %d data error: %s", market, err.Error()))
			} else {
				fxj = append(fxj, f...)
			}
		}

	}
	return fxj, nil
}

func getFxjData(req *http.Request, market uint8, stockType _type) (fxjs []orm.Fxj, err error) {
	params := make(url.Values)
	codes, code2hqcode, err := getMarketCodes(market)
	if err != nil {
		return
	}
	params.Set("xml_request", getXmlStr(codes, stockType))
	req.URL.RawQuery = params.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(s, &m)
	if err != nil {
		return
	}
	fxjs = make([]orm.Fxj, 0)
	mtime := time.Now()
	tradedate := getTradedate()
	idx := 0
	rows := m["rows"].(map[string]interface{})
	for code, data := range rows {
		if data.([]interface{})[1] == nil || data.([]interface{})[1].(string) == "" {
			continue
		}
		strFxj := data.([]interface{})[1].(string)
		floatFxj, _ := strconv.ParseFloat(strFxj, 64)
		fxjs = append(fxjs, orm.Fxj{
			StockCode: code2hqcode[code],
			Market:    int(market),
			TradeDate: tradedate,
			Mtime:     mtime,
			Fxj:       floatFxj,
		})
		idx++
	}

	return
}

// 和业务方确认过，有且仅有仅该三支类型
func getBaseXml(stockType _type) XMLRequest {
	switch stockType {
	case index:
		return indexXml
	case fund:
		return fundXml
	case bond:
		return bondXml
	default:
		return XMLRequest{}
	}
}

func getMarketCodes(market uint8) (string, map[string]string, error) {
	db := getIfindPg()
	sql := fmt.Sprintf("select thscode_hq,thscode from pub205 where thsmarket_code_hq = '%d' "+
		"and thscode is not null and thscode_hq is not null", marketUint8toInt8(market))
	rows, err := db.Query(sql)
	log.Log.Info("SQL query: " + sql)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()
	m := make(map[string]string)
	codes := new(bytes.Buffer)
	var code, hqcode string
	for rows.Next() {
		err = rows.Scan(&hqcode, &code)
		if err != nil {
			return "", nil, err
		}
		m[code] = hqcode
		codes.WriteString(code)
		codes.WriteString(",")
	}
	return strings.Trim(codes.String(), ","), m, nil
}

func getXmlStr(codes string, stockType _type) string {
	x := getBaseXml(stockType)
	x.Thscodes = codes
	output, _ := xml.MarshalIndent(&x, "", "")
	return string(output)
}

// 针对发行价的接口，一般在15:03左右更新数据
func getTradedate() string {
	t := time.Now()
	h := t.Hour()
	if h < 15 {
		t.AddDate(0, 0, -1)
	}
	y, m, d := t.Date()
	return fmt.Sprintf("%4d-%02d-%02d", y, m, d)
}
