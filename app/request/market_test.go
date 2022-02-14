package request

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarketTransform(t *testing.T) {
	uMarket := 177
	market := marketUint8toInt8(uint8(uMarket))
	assert.Equal(t, int8(uMarket), market)
	assert.Equal(t, "-79", fmt.Sprintf("%d", market))
}
