// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

import (
	"time"
)

// StoresSkuStockChange is the database model struct
type StoresSkuStockChange struct {
	Id          int64     `json:"id" xorm:"pk autoincr BIGINT(20)"`
	SkuId       int64     `json:"sku_id" xorm:"not null BIGINT(20)"`
	ShelfId     int64     `json:"shelf_id" xorm:"not null BIGINT(20)"`
	Layer       string    `json:"layer" xorm:"VARCHAR(64)"`
	Price       float64   `json:"price" xorm:"not null DOUBLE(10,2)"`
	Amount      int64     `json:"amount" xorm:"not null default 0 BIGINT(20)"`
	Type        int       `json:"type" xorm:"not null TINYINT(4)"`
	OperatorId  int64     `json:"operator_id" xorm:"not null BIGINT(20)"`
	CreatedTime time.Time `json:"created_time" xorm:"not null index DATETIME"`
	Detail      string    `json:"detail" xorm:"not null default '' VARCHAR(1024)"`
}
