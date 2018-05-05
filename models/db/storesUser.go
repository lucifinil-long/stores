// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

import (
	"time"
)

// StoresUser is the database model struct
type StoresUser struct {
	Id            int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Mobile        int64     `json:"mobile" xorm:"not null unique BIGINT(20)"`
	Nickname      string    `json:"nickname" xorm:"not null default '' VARCHAR(128)"`
	Password      string    `json:"password" xorm:"not null default '' VARCHAR(64)"`
	Remark        string    `json:"remark" xorm:"VARCHAR(512)"`
	Deletable     int       `json:"deletable" xorm:"not null default 1 TINYINT(4)"`
	Deleted       int       `json:"deleted" xorm:"not null default 0 TINYINT(4)"`
	LastLoginTime time.Time `json:"last_login_time" xorm:"DATETIME"`
	CreatedTime   time.Time `json:"created_time" xorm:"not null DATETIME"`
}
