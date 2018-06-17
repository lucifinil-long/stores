// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresCommoditySku is the database model struct
type StoresCommoditySku struct {
	Id          int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Name        string `json:"name" xorm:"not null VARCHAR(64)"`
	Barcode     string `json:"barcode" xorm:"not null default '' unique VARCHAR(64)"`
	SpecId      int64  `json:"spec_id" xorm:"not null BIGINT(20)"`
	Profit      int    `json:"profit" xorm:"not null default 20 INT(11)"`
	MaxDiscount int    `json:"max_discount" xorm:"default 0 TINYINT(4)"`
}
