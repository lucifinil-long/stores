// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresLocationDepot is the database model struct
type StoresLocationDepot struct {
	Id     int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Name   string `json:"name" xorm:"not null unique VARCHAR(64)"`
	Detail string `json:"detail" xorm:"not null default '' VARCHAR(64)"`
}
