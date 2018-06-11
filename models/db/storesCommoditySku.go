// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresCommoditySku is the database model struct
type StoresCommoditySku struct {
	Id      int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Name    string `json:"name" xorm:"not null VARCHAR(64)"`
	Barcode string `json:"barcode" xorm:"not null default '' unique VARCHAR(64)"`
	SpecId  int64  `json:"spec_id" xorm:"not null index BIGINT(20)"`
}
