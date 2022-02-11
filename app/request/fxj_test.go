package request

import (
	"encoding/json"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var FxjTest = "{\"version\":0,\"attrs\":[],\"descrs\":[{\"name\":\"THS_CODE\",\"type\":4,\"attrs\":[]},{\"name\":\"653092_000_00_0_3\",\"type\":6,\"attrs\":[]}],\"rows\":{\"000001.SH\":[\"000001.SH\",\"26.420009191777\"]}}"

func TestJsonDecode(t *testing.T) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(FxjTest), &m)
	assert.NoError(t, err, "except no error")
	assert.Equal(t, "000001.SH", m["rows"].(map[string]interface{})["000001.SH"].([]interface{})[0].(string))
	assert.Equal(t, "26.420009191777", m["rows"].(map[string]interface{})["000001.SH"].([]interface{})[1].(string))
}

func TestXml(t *testing.T) {
	output, err := xml.MarshalIndent(&indexXml, "", "    ")
	assert.NoError(t, err, "except no error")
	log.Println(xml.Header + string(output))
	output, err = xml.MarshalIndent(&fundXml, "", "    ")
	assert.NoError(t, err, "except no error")
	log.Println(xml.Header + string(output))
	output, err = xml.MarshalIndent(&bondXml, "", "    ")
	assert.NoError(t, err, "except no error")
	log.Println(xml.Header + string(output))
}
