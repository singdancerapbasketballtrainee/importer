package orm

import "time"

type Fxj struct {
	StockCode string    `gorm:"type:varchar(10);column:stockcode"`
	Market    int       `gorm:"type:int;column:market"`
	TradeDate string    `gorm:"type:date;column:tradedate"`
	Mtime     time.Time `gorm:"type:timestamp;column:mtime"`
	Fxj       float64   `gorm:"type:numeric(16,12);column:fxj"`
}
