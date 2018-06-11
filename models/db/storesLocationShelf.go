// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresLocationShelf is the database model struct
type StoresLocationShelf struct {
	Id      int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	DepotId int64  `json:"depot_id" xorm:"not null unique(idx_uni) BIGINT(20)"`
	Name    string `json:"name" xorm:"not null unique(idx_uni) VARCHAR(64)"`
	Layers  int    `json:"layers" xorm:"not null default 1 TINYINT(4)"`
	Detail  string `json:"detail" xorm:"not null default '' VARCHAR(64)"`
}
