package db

// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

type StoresUserNode struct {
	Id     int `json:"id" xorm:"not null pk autoincr INT(11)"`
	UserId int `json:"user_id" xorm:"not null index INT(11)"`
	NodeId int `json:"node_id" xorm:"not null index INT(11)"`
}
