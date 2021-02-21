// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresSkuPropertyValue is the database model struct
type StoresSkuPropertyValue struct {
	SkuId      int64  `json:"sku_id" xorm:"not null pk BIGINT(20)"`
	PropertyId int64  `json:"property_id" xorm:"not null pk BIGINT(20)"`
	Value      string `json:"value" xorm:"not null VARCHAR(64)"`
}
