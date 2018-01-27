package db

// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

import (
	"time"
)

type StoresOpLog struct {
	Id          int64     `json:"id" xorm:"pk autoincr BIGINT(20)"`
	UserId      int       `json:"user_id" xorm:"not null INT(11)"`
	Username    string    `json:"username" xorm:"not null VARCHAR(128)"`
	From        string    `json:"from" xorm:"not null VARCHAR(128)"`
	Action      string    `json:"action" xorm:"not null VARCHAR(128)"`
	Detail      string    `json:"detail" xorm:"TEXT"`
	CreatedTime time.Time `json:"created_time" xorm:"not null DATETIME"`
}
