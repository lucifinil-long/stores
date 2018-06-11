package models

import (
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models/db"
)

func getUserRoles(session *xorm.Session, uid int64) (string, error) {
	if session == nil {
		session = config.GetConfigs().OrmEngine.NewSession()
		defer session.Close()
	}

	records := make([]*db.StoresRole, 0)

	err := session.Table(cTableStoresRole).
		Join("inner", cTableStoresRoleUser, "stores_role.id=stores_role_user.role_id").
		Where("stores_role_user.user_id=?", uid).
		Cols("role_name").
		Find(&records)

	if err != nil {
		return "", err
	}

	names := make([]string, 0, len(records))
	for _, role := range records {
		names = append(names, role.RoleName)
	}

	return strings.Join(names, ", "), nil
}
