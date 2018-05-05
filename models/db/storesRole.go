// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresRole is the database model struct
type StoresRole struct {
	Id        int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	RoleName  string `json:"role_name" xorm:"not null unique VARCHAR(128)"`
	Remark    string `json:"remark" xorm:"VARCHAR(512)"`
	Deletable int    `json:"deletable" xorm:"not null default 1 TINYINT(4)"`
}
