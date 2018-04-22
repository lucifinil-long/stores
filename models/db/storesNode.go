package db

// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

/*
 *	StoresNode is a database model struct
 */
type StoresNode struct {
	Id     int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	Title  string `json:"title" xorm:"not null default '' VARCHAR(100)"`
	Path   string `json:"path" xorm:"not null default '' VARCHAR(256)"`
	Level  int    `json:"level" xorm:"not null default 1 INT(11)"`
	Pid    int    `json:"pid" xorm:"not null default 0 index INT(11)"`
	Menu   int    `json:"menu" xorm:"not null default 0 TINYINT(4)"`
	Auth   int    `json:"auth" xorm:"default 1 TINYINT(4)"`
	Icon   string `json:"icon" xorm:"VARCHAR(256)"`
	Remark string `json:"remark" xorm:"VARCHAR(200)"`
}
