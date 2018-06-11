// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package db

// StoresRoleUser is the database model struct
type StoresRoleUser struct {
	RoleId int64 `json:"role_id" xorm:"not null pk BIGINT(20)"`
	UserId int64 `json:"user_id" xorm:"not null pk index BIGINT(20)"`
}
