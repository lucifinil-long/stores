package db

// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

/*
 *	StoresRoleNodes is a database model struct
 */
type StoresRoleNodes struct {
	RoleId int `json:"role_id" xorm:"not null pk INT(11)"`
	NodeId int `json:"node_id" xorm:"not null pk index INT(11)"`
}
