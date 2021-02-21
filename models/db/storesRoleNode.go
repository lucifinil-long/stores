// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresRoleNode is the database model struct
type StoresRoleNode struct {
	RoleId int64 `json:"role_id" xorm:"not null pk BIGINT(20)"`
	NodeId int64 `json:"node_id" xorm:"not null pk BIGINT(20)"`
}
