package models

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models/db"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

// dbUser2ProtoUser translate db user to protocol user
func dbUser2ProtoUser(user *db.StoresUser) *proto.User {
	if user == nil {
		return &proto.User{}
	}

	return &proto.User{
		ID:            user.Id,
		Username:      user.Username,
		Nickname:      user.Nickname,
		Password:      user.Password,
		Mobile:        user.Mobile,
		Remark:        user.Remark,
		Status:        user.Status,
		Level:         user.Level,
		LastLoginTime: user.LastLoginTime.Format("2006-01-02T15:04:05-07:00"),
		CreatedTime:   user.CreatedTime.Format("2006-01-02T15:04:05-07:00"),
	}
}

// protoUser2DBUser translate protocol user to db user
func protoUser2DBUser(user *proto.User) *db.StoresUser {
	if user == nil {
		return &db.StoresUser{}
	}

	loginTime, _ := time.Parse("2006-01-02T15:04:05-07:00", user.LastLoginTime)
	cretedTime, _ := time.Parse("2006-01-02T15:04:05-07:00", user.CreatedTime)
	return &db.StoresUser{
		Id:            user.ID,
		Username:      user.Username,
		Nickname:      user.Nickname,
		Password:      user.Password,
		Mobile:        user.Mobile,
		Remark:        user.Remark,
		Status:        user.Status,
		Level:         user.Level,
		LastLoginTime: loginTime,
		CreatedTime:   cretedTime,
	}
}

// protoUser2DBUser translate protocol user to db user
func protoNewUser2DBUser(user *proto.NewUser) *db.StoresUser {
	if user == nil {
		return &db.StoresUser{}
	}

	return &db.StoresUser{
		Id:            0,
		Username:      user.Username,
		Nickname:      user.Nickname,
		Password:      user.Password,
		Mobile:        user.Mobile,
		Remark:        user.Remark,
		Status:        user.Status,
		CreatedTime:   time.Now(),
		LastLoginTime: time.Now(),
	}
}

// assignAccessToRole assign accesses to user
func assignAccessToRole(session *xorm.Session, rid int, accesses []int) error {
	records := []db.StoresRoleNode{}
	err := session.Table(cTableStoresRoleNode).Where("role_id=?", rid).Find(&records)
	sql, params := session.LastSQL()
	log.Trace("models.assignAccessToRole: query sql: `%v`, parameters: %v", sql, params)
	if err != nil {
		return err
	}

	// check whether need to update access
	if len(accesses) == len(records) {
		same := true
		for _, record := range records {
			found := false
			for _, id := range accesses {
				if id == record.NodeId {
					found = true
					break
				}
			}

			if !found {
				same = false
				break
			}
		}

		if same {
			log.Trace("models.assignAccessToRole returned for user access is not changed.")
			return nil
		}
	}

	// remove current access of user before assign new accesses
	if err = removeAccessOfRole(session, rid); err != nil {
		return err
	}

	inserts := make([]db.StoresRoleNode, 0, len(accesses))
	for _, val := range accesses {
		record := db.StoresRoleNode{RoleId: rid, NodeId: val}
		inserts = append(inserts, record)
	}
	_, err = session.Table(cTableStoresRoleNode).InsertMulti(inserts)
	sql, params = session.LastSQL()
	log.Trace("models.assignAccessToUser: sql: `%v`, parameters: %v", sql, params)

	return err
}

// removeAccessOfRole remove accesses of user
func removeAccessOfRole(session *xorm.Session, rid int) error {
	_, err := session.Table(cTableStoresRoleNode).Where("role_id=?", rid).Delete(db.StoresRoleNode{})
	sql, params := session.LastSQL()
	log.Trace("models.removeAccessOfRole: sql: `%v`, parameters: %v, error: %v", sql, params, err)

	return err
}

// IsSuperAdmin test whether user is super administrator role
func IsSuperAdmin(user *proto.User) bool {
	return user != nil && user.Level == -1
}

// GetUserAccessList get access permissions list
// @param uid is the user id of specified user
// @return (map[string]bool, nil) if get access list successfully, otherwise return (nil, error)
//	Note that if no node that user can access, also return (empty map[string]bool, nil)
func GetUserAccessList(uid int) (map[string]bool, error) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	// 由于节点权限控制设计时只会最终在需要授权的节点生效，因此本处查询仅需查询能访问的叶子节点即可
	nodes, err := getUserAccessList(session, uid, false)
	if err != nil {
		return nil, err
	}

	accesslist := make(map[string]bool, len(nodes))

	for _, node := range nodes {
		if len(node.Path) > 0 {
			accesslist[node.Path] = true
		}
	}

	return accesslist, nil
}

// getUserAccessList get the node list that specified user can access
// @param session is the database session, can be nil; if nil will use default database session
// @param uid is the user id of specified user
// @param onlyAuthNode indicates whether only auth need nodes are returned
// @return ([]*db.StoresNode, nil) if successful; otherwise return ([]*db.StoresNode, error)
func getUserAccessList(session *xorm.Session, uid int, onlyAuthNode bool) ([]*db.StoresNode, error) {
	if session == nil {
		session = config.GetConfigs().OrmEngine.NewSession()
		defer session.Close()
	}

	records := make([]*db.StoresNode, 0)
	session.Table(cTableStoresNode).
		Where("id in (select node_id from stores_user_node where user_id=?)", uid)
	if onlyAuthNode {
		session.And("auth = 1")
	} else {
		session.Or("auth = 0")
	}
	err := session.Find(&records)

	sql, params := session.LastSQL()
	log.Trace("models.GetUserAccessList: sql: `%v`, parameters: %v", sql, params)

	return records, err
}

// GetUserInfoByUsername validates login info
// @param session is database xorm session
// @param username is the username of user that be retrived information
// @return (*proto.User, nil) if get user infomation successfuly, otherwise return (nil, error)
func GetUserInfoByUsername(username string) (*proto.User, error) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	user, err := getUserInfoByUsername(session, username)
	if err != nil {
		return nil, err
	}

	return dbUser2ProtoUser(user), nil
}

// getUserInfoByUsername validates login info
// @param session is database xorm session
// @param username is the username of user that be retrived information
// @return (*db.StoresUser, nil) if get user infomation successfuly, otherwise return (nil, error)
func getUserInfoByUsername(session *xorm.Session, username string) (*db.StoresUser, error) {
	user := &db.StoresUser{}
	found, err := session.Where("username=?", username).Get(user)

	sql, params := session.LastSQL()
	log.Trace("models.checkUserInfo: sql: `%v`, parameters: %v", sql, params)

	if err != nil {
		return nil, err
	}

	if !found || user.Id == 0 {
		return nil, proto.ErrUserNotExist
	}

	if user.Status == 0 {
		return nil, proto.ErrUserDisabled
	}

	return user, nil
}

// UpdateUserLoginTime updates user last login time
// @param uid is the user id of specified user
func UpdateUserLoginTime(uid int) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	updateUserLoginTime(session, uid)
}

// updateUserLoginTime updates user last login time to DB
// @param session is the database session, can not be nil
// @param uid is the user id of specified user
func updateUserLoginTime(session *xorm.Session, uid int) {
	user := db.StoresUser{Id: uid, LastLoginTime: time.Now()}
	updates, err := session.Table(user).Where("id=?", uid).Cols("last_login_time").Update(user)
	if err != nil || updates != 1 {
		log.Error("models.updateUserLoginTime: updated record with error: %v, update count: %v", err, updates)
	} else {
		log.Trace("models.updateUserLoginTime: updated user (ID:%v) login time", uid)
	}
}

// UpdateUserPassword update user's password
// @param uid is the user id of the user will be updated
// @param newPwd is new password for the user
// @return nil if successfully, otherwise return error
func UpdateUserPassword(uid int, newPwd string) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	return updateUserPassword(session, uid, newPwd)
}

// updateUserPassword update user's password to DB
// @param session is the database session, can not be nil
// @param uid is the user id of the user will be updated
// @param newPwd is new password for the user
// @return nil if successfully, otherwise return error
func updateUserPassword(session *xorm.Session, uid int, newPwd string) error {
	user := db.StoresUser{
		Id:       uid,
		Password: newPwd,
	}

	_, err := session.Table(user).Where("id = ?", uid).Cols("password").Update(user)

	return err
}

// GetUserList get user list
// @param page is the current page index
// @param pageSize is the page size
// @param sort is the table column name that used to be sorted
// @param desc indicates whether sort data desc
// @return ([]proto.User, total record count, nil) if successfully, otherwise return (empty []proto.User, 0, error)
func GetUserList(pageIndex, pageSize int, sort string, desc bool) ([]proto.User, int64, error) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	return getUserList(session, pageIndex, pageSize, sort, desc)
}

// getUserList get user list from DB
// @param session is the database session, can not be nil
// @param pageIndex is the current page index
// @param pageSize is the page size
// @param sort is the table column name that used to be sorted
// @param desc indicates whether sort data desc
// @return ([]proto.User, total record count, nil) if successfully, otherwise return (empty []proto.User, 0, error)
func getUserList(session *xorm.Session, pageIndex, pageSize int, sort string, desc bool) ([]proto.User, int64, error) {
	records := make([]*db.StoresUser, 0)
	offset := 0
	if pageIndex > 1 {
		offset = (pageIndex - 1) * pageSize
	}
	session.Table(cTableStoresUser).
		Cols("id,username,nickname,mobile,remark,status,last_login_time,created_time").
		Limit(pageSize, offset)
	if len(sort) > 0 {
		if desc {
			session.Desc(sort)
		} else {
			session.Asc(sort)
		}
	}

	err := session.Find(&records)
	sql, params := session.LastSQL()
	log.Trace("models.getUserList: sql: `%v`, parameters: %v", sql, params)

	if err != nil {
		return []proto.User{}, 0, err
	}

	count, err := session.Table(cTableStoresUser).Count()
	if err != nil {
		return []proto.User{}, 0, err
	}

	ret := make([]proto.User, 0, len(records))
	for _, record := range records {
		user := dbUser2ProtoUser(record)
		ret = append(ret, *user)
	}

	return ret, count, err
}

// AddUser handles add admin user to database request
// @param user is admin user information
// @return nil if successful; otherwise return an error
func AddUser(user *proto.NewUser) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return err
	}

	err = addUser(session, user)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.AddUser: rollback failed with %v", rollbackErr)
		}
		return err
	}

	return session.Commit()
}

// addUser adds db user record and related access
func addUser(session *xorm.Session, user *proto.NewUser) error {
	dbUser := protoNewUser2DBUser(user)

	err := addDBUser(session, dbUser)
	if err != nil {
		return err
	}

	if dbUser.Id <= 0 {
		log.Error("models.addUser: user id (%v) is invalid after add user to database.", dbUser.Id)
		return proto.ErrCommonInternalError
	}

	return nil
}

// addDBUser handles add user to database request
// @param session is database session, can be nil; if nil will use default database session
// @param user is db user information
func addDBUser(session *xorm.Session, user *db.StoresUser) error {
	if len(user.Username) == 0 {
		return proto.ErrCommonInvalidParam
	}

	_, err := session.Insert(user)
	sql, params := session.LastSQL()
	log.Trace("models.addDBUser: query sql: `%v`, parameters: %v", sql, params)

	return err
}

// DeleteUser delete specified user
// @param uid is the user id in database
// @return nil if successful; otherwise return an error
func DeleteUser(uid int) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return err
	}

	err = deleteUser(session, uid)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.AddUser: rollback failed with %v", rollbackErr)
		}
		return err
	}

	return session.Commit()
}

// deleteUser delete specified user from user table and
// @param session is the database session, can be nil; if nil will use default database session
// @param uid is the user id in database
// @return nil if successful; otherwise return an error
func deleteUser(session *xorm.Session, uid int) error {
	err := deleteDBUser(session, uid)
	if err != nil {
		return err
	}

	return nil
}

// deleteDBUser marks specified user in user table is deleted
// @param session is the database session, can be nil; if nil will use default database session
// @param uid is the user id in database
// @return nil if successful; otherwise return an error
func deleteDBUser(session *xorm.Session, uid int) error {
	user := &db.StoresUser{Id: uid, Deleted: 1}
	_, err := session.Table(user).
		Where("id=?", uid).
		And("deleted=1").
		Cols("deleted").
		Update(user)
	return err
}

// GetUser get an user information from database
// @param uid is the user id in database
// @param columns is user database columns name to be queried
// @return (*db.StoresUser, nil) if successful; otherwise return (nil, error)
func GetUser(uid int, columns ...string) (*db.StoresUser, error) {

	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	user := &db.StoresUser{Id: uid}
	session.Table(user).Where("id=?", uid)
	if len(columns) > 0 {
		session.Cols(columns...)
	}
	found, err := session.Get(user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, proto.ErrUserNotExist
	}

	return user, nil
}
