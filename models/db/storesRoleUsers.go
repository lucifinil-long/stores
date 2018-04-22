package db

// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

/*
 *	StoresRoleUsers is a database model struct
 */
type StoresRoleUsers struct {
	RoleId int `json:"role_id" xorm:"not null pk INT(11)"`
	UserId int `json:"user_id" xorm:"not null pk index INT(11)"`
}
