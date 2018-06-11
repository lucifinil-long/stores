// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresSkuStock is the database model struct
type StoresSkuStock struct {
	SkuId   int64 `json:"sku_id" xorm:"not null pk BIGINT(20)"`
	ShelfId int64 `json:"shelf_id" xorm:"not null pk index BIGINT(20)"`
	Amount  int64 `json:"amount" xorm:"not null default 0 BIGINT(20)"`
}
