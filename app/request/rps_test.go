package request

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetRps(t *testing.T) {
	rps, err := GetRps("2020-01-27")
	assert.NoError(t, err, "expected no error")
	log.Println(rps)
}
