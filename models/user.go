package models

import (
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models/db"
	"github.com/lucifinil-long/stores/proto"
	"github.com/lucifinil-long/stores/utils"
	"github.com/mkideal/log"
)

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
	log.Trace("models.GetUserAccessList: query sql: `%v`, parameters: %v", sql, params)

	return records, err
}

// dbUser2ProtoUser translate db user to protocol user
func dbUser2ProtoUser(user *db.StoresUser) *proto.User {
	return &proto.User{
		ID:            user.Id,
		Username:      user.Username,
		Nickname:      user.Nickname,
		Mobile:        user.Mobile,
		Remark:        user.Remark,
		Status:        user.Status,
		Level:         user.Level,
		LastLoginTime: user.LastLoginTime.Format("2006-01-02T15:04:05-07:00"),
		CreatedTime:   user.CreatedTime.Format("2006-01-02T15:04:05-07:00"),
	}
}

// CheckUserInfo validates login info
// @param session is database xorm session
// @param username is the username that be checked
// @param password is the correspond password that be checked
// @return (*proto.User, nil) if validate user infomation successfuly, otherwise return (nil, error)
func CheckUserInfo(username, password string) (*proto.User, error) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	user, err := checkUserInfo(session, username, password)
	if err != nil {
		return nil, err
	}

	return dbUser2ProtoUser(user), nil
}

// CheckUserInfo validates login info
// @param session is database xorm session
// @param username is the username that be checked
// @param password is the correspond password that be checked
// @return (*db.StoresUser, nil) if validate user infomation successfuly, otherwise return (nil, error)
func checkUserInfo(session *xorm.Session, username, password string) (*db.StoresUser, error) {
	if session == nil {
		session = config.GetConfigs().OrmEngine.NewSession()
		defer session.Close()
	}

	user := &db.StoresUser{}
	found, err := session.Where("username=?", username).Get(user)

	sql, params := session.LastSQL()
	log.Trace("models.checkUserInfo: query sql: `%v`, parameters: %v", sql, params)

	if err != nil {
		return nil, err
	}

	if !found || user.Id == 0 {
		return nil, ErrUserNotExist
	}

	if user.Status == 0 {
		return nil, ErrUserDisabled
	}

	passwordMd5 := strings.ToLower(utils.String2MD5(password))
	password = strings.ToLower(password)
	userDBPwd := strings.ToLower(user.Password)

	if userDBPwd != password && userDBPwd != passwordMd5 {
		return nil, ErrUserWrongPwd
	}

	return user, nil
}
