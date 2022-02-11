package request

import (
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

var (
	IndexXmlStr string
	FundXmlStr  string
	BondXmlStr  string
)

const (
	indexCode  = "000001"
	indexCode_ = "000001.SH"
	fundCode   = "160127"
	fundCode_  = "160127.OF"
	bondCode   = "110065"
	bondCode_  = "110065.SH"
)

func init() {
	baseParams := []Param{
		{Name: "FD0", Value: "111", System: "true"},
	}

	indexXml = XMLRequest{Thscodes: indexCode_}
	item := Item{Name: "06971_000_00_0_1", Params: baseParams}
	indexXml.Items = append(indexXml.Items, item)

	fundXml = XMLRequest{Thscodes: fundCode_}
	item = Item{Name: "00038_000_00_0_5", Params: baseParams}
	item.Params = append(item.Params, Param{Name: "FT", Value: "100", System: "true"})
	fundXml.Items = append(fundXml.Items, item)

	bondXml = XMLRequest{Thscodes: bondCode_}
	item = Item{Name: "03548_000_00_0_11", Params: baseParams}
	item.Params = append(item.Params, Param{Name: "FJ", Value: "103", System: "true"})
	bondXml.Items = append(bondXml.Items, item)
}

func init() {
	//IndexXmlStr = `<?xml version="1.0" encoding="UTF-8"?>
	//<request>
	//<thscodes>000001.SH</thscodes>
	//<items>
	//<item name="06971_000_00_0_1">
	//<params>
	//<param name="FD0" value="111" system="true"/>
	//</params>
	//</item>
	//</items>
	//</request>`
	output, _ := xml.MarshalIndent(&indexXml, "", "")
	IndexXmlStr = string(output)

	//FundXmlStr = `<?xml version="1.0" encoding="UTF-8"?>
	//<request>
	//<thscodes>160201.OF</thscodes>
	//<items>
	//<item name="00038_000_00_0_5">
	//<params>
	//<param name="FD0" value="111" system="true"/>
	//<param name="FT" value="100" system="true"/>
	//</params>
	//</item>
	//</items>
	//</request>`
	output, _ = xml.MarshalIndent(&fundXml, "", "")
	FundXmlStr = string(output)

	//BondXmlStr = `<?xml version="1.0" encoding="UTF-8"?>
	//<request>
	//<thscodes>110065.SH</thscodes>
	//<items>
	//<item name="03548_000_00_0_11">
	//<params>
	//<param name="FD0" value="111" system="true"/>
	//<param name="FJ" value="103" system="true"/>
	//</params>
	//</item>
	//</items>
	//</request>`
	output, _ = xml.MarshalIndent(&bondXml, "", "")
	BondXmlStr = string(output)
}

// GetFxj 发行价的接口设计写得比较死，三支代码只能一个个来
func GetFxj() ([]orm.Fxj, error) {
	fxj := make([]orm.Fxj, 0)
	req, err := http.NewRequest("POST", config.GetApiConfig().FxjCfg.Url, nil)
	if err != nil {
		return fxj, err
	}
	// 设置http头
	req.Header.Set("X-Arsenal-Auth", config.GetApiConfig().AppName)

	codes := []string{indexCode_, fundCode_, bondCode_}
	for _, code := range codes {
		f, err := getFxjData(req, code)
		if err != nil {
			log.Log.Error(fmt.Sprintf("get %s data error: %s", code, err.Error()))
		} else {
			fxj = append(fxj, f)
		}
	}
	return fxj, nil
}

func getFxjData(req *http.Request, code string) (fxj orm.Fxj, err error) {
	params := make(url.Values)
	params.Set("xml_request", getXmlStr(code))
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
	// 根据接口特点，此处不需要判断 ()[1] 是否存在
	if m["rows"].(map[string]interface{})[code].([]interface{})[1] == nil {
		err = fmt.Errorf("get fxl failed")
		return
	}
	strFxj := m["rows"].(map[string]interface{})[code].([]interface{})[1].(string)
	floatFxj, err := strconv.ParseFloat(strFxj, 64)
	if err != nil {
		return
	}
	return orm.Fxj{
		StockCode: getHqCode(code),
		Market:    getMarket(code),
		TradeDate: getTradedate(),
		Mtime:     time.Now(),
		Fxj:       floatFxj,
	}, nil
}

// 接口返回不带市场，提过该函数获取市场，因为只有三种类型，指数、基金、债券，所以直接用switch，多了再和getHqCode一起优化
func getMarket(code string) int {
	switch code {
	case indexCode_:
		return 20
	case fundCode_:
		return 35
	case bondCode_:
		return 105
	default:
		return 0
	}
}

// 根据接口返回的带后缀的代码获取不带后缀的原代码，因为只有三种类型，指数、基金、债券，所以直接用switch，多了再和getMarket一起优化
func getHqCode(code string) string {
	switch code {
	case indexCode_:
		return indexCode
	case fundCode_:
		return fundCode
	case bondCode_:
		return bondCode
	default:
		return ""
	}
}

// 同上，先写死，数据接口优化了这里一起优化
func getXmlStr(code string) string {
	switch code {
	case indexCode_:
		return IndexXmlStr
	case fundCode_:
		return FundXmlStr
	case bondCode_:
		return BondXmlStr
	default:
		return ""
	}
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
