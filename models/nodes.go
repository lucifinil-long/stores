package models

import (
	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models/db"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

// GetTreeMenuForUser get tree node struct for user
// @param user is specified user
// @return tree menu nodes
func GetTreeMenuForUser(user *proto.User) []*proto.TreeMenuNode {
	if user == nil {
		return []*proto.TreeMenuNode{}
	}

	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	accessIds := []interface{}{}

	if !IsSuperAdmin(user) {
		nodes, _ := getUserAccessList(session, user.ID, false)
		for _, node := range nodes {
			accessIds = append(accessIds, node.Id)
		}
	}

	return loadTree(session, 0, true, true, accessIds...)
}

// loadTree loads tree according to parent id and levels
// @param session is the database session, can be nil; if nil will use default database session
// @param pid is the top parent node id
// @param menu is the flag of tree menu
// @param recursive indicates whether load tree nodes recursively
// @param accessIds is user access node list
// @return first sub level tree node list of specified parent id node, not includes parent node
func loadTree(session *xorm.Session, pid int, menu, recursive bool, accessIds ...interface{}) []*proto.TreeMenuNode {
	if session == nil {
		session = config.GetConfigs().OrmEngine.NewSession()
		defer session.Close()
	}

	nodes, err := getSubTreeNodes(session, pid, menu, accessIds...)
	if err != nil {
		log.Error("database operation error: %v", err)
		return []*proto.TreeMenuNode{}
	}
	if len(nodes) == 0 {
		return []*proto.TreeMenuNode{}
	}

	treeNodes := make([]*proto.TreeMenuNode, 0, len(nodes))

	for _, v := range nodes {
		node := &proto.TreeMenuNode{}
		node.ID = v.Id
		node.Text = v.Title
		node.PID = pid
		node.IconCls = v.Icon
		if recursive {
			node.Children = loadTree(session, v.Id, menu, recursive, accessIds...)
		} else {
			node.Children = []*proto.TreeMenuNode{}
		}

		node.Attributes.URL = v.Path

		// 仅叶子节点和有子节点的才加入菜单节点列表
		if len(v.Path) > 0 || len(node.Children) > 0 {
			treeNodes = append(treeNodes, node)
		}
	}

	return treeNodes
}

// getSubTreeNodes gets sub tree nodes
// @param session is the database session, can be nil; if nil will use default database session
// @param pid is the id of parent node id
// @param menu is the flag of tree menu
// @param accessIds is id set of access nodes that user can access, include no need auth nodes; if empty, will not check auth
// @return ([]*db.StoresNode, nil) if successful; otherwise return ([]*db.GwAdminNode, error)
func getSubTreeNodes(session *xorm.Session, pid int, menu bool, accessIds ...interface{}) ([]*db.StoresNode, error) {
	if session == nil {
		session = config.GetConfigs().OrmEngine.NewSession()
		defer session.Close()
	}

	records := make([]*db.StoresNode, 0)
	session.Table(cTableStoresNode)
	if len(accessIds) > 0 {
		session.In("id", accessIds...).And("pid=?", pid)
	} else {
		session.Where("pid=?", pid)
	}

	if menu {
		session.And("menu=1")
	}
	err := session.Find(&records)

	// sql, params := session.LastSQL()
	// log.Trace("models.GetSubTreeNodes: query sql: `%v`, parameters: %v", sql, params)

	return records, err
}
