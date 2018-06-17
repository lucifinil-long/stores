// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresSkuProperty is the database model struct
type StoresSkuProperty struct {
	Id       int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Property string `json:"property" xorm:"not null VARCHAR(64)"`
}
