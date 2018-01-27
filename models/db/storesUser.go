package db

// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

import (
	"time"
)

type StoresUser struct {
	Id            int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Username      string    `json:"username" xorm:"not null default '' unique VARCHAR(128)"`
	Password      string    `json:"password" xorm:"not null default '' VARCHAR(64)"`
	Nickname      string    `json:"nickname" xorm:"not null default '' VARCHAR(128)"`
	Mobile        string    `json:"mobile" xorm:"default '' VARCHAR(128)"`
	Remark        string    `json:"remark" xorm:"VARCHAR(512)"`
	Status        int       `json:"status" xorm:"not null default 1 TINYINT(4)"`
	Level         int       `json:"level" xorm:"not null default 1 TINYINT(4)"`
	LastLoginTime time.Time `json:"last_login_time" xorm:"DATETIME"`
	CreatedTime   time.Time `json:"created_time" xorm:"not null DATETIME"`
}
