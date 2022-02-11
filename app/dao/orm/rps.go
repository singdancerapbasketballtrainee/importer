package orm

import (
	"time"
)

type Rps struct {
	StockCode string    `gorm:"type:varchar(10);column:stockcode"`
	TradeDate string    `gorm:"date;column:tradedate"`
	Mtime     time.Time `gorm:"type:timestamp;column:mtime"`
	RpsMin    int       `gorm:"type:int;column:rps_min"`
	RpsMax    int       `gorm:"type:int;column:rps_max"`
}
