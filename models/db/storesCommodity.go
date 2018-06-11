// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresCommodity is the database model struct
type StoresCommodity struct {
	Id   int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Name string `json:"name" xorm:"not null VARCHAR(128)"`
}
