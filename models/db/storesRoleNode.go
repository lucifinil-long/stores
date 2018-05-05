// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresRoleNode is the database model struct
type StoresRoleNode struct {
	RoleId int `json:"role_id" xorm:"not null pk INT(11)"`
	NodeId int `json:"node_id" xorm:"not null pk index INT(11)"`
}
