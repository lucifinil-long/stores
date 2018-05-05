package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lucifinil-long/stores/models"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

func TestModels(t *testing.T) {
	log.Info("Models unittests start...")
	time.Sleep(10 * time.Millisecond)
}

func TestGetUserInfoByUserIDOrMobile(t *testing.T) {
	users, err := models.GetUserInfoByUserIDOrMobile(-1)
	result, _ := json.MarshalIndent(users, "", "    ")

	log.Info("TestGetUserInfoByUserIDOrMobile: result for user id:\n%v\nerror:%v", string(result), err)

	users, err = models.GetUserInfoByUserIDOrMobile(12345678901)
	result, _ = json.MarshalIndent(users, "", "    ")

	log.Info("TestGetUserInfoByUserIDOrMobile: result for user mobile:\n%v\nerror:%v", string(result), err)

	log.Info("TestGetUserInfoByUserIDOrMobile is done.\n\n")
	time.Sleep(100 * time.Millisecond)
}

func TestGetUserAccessList(t *testing.T) {
	accesses, err := models.GetUserAccessList(-1)
	log.Info("TestGetUserAccessList: result for super admin:\n%v\nerror:%v", accesses, err)
	accesses, err = models.GetUserAccessList(1000)
	log.Info("TestGetUserAccessList: result for super admin:\n%v\nerror:%v", accesses, err)

	log.Info("TestGetUserAccessList is done.\n\n")
	time.Sleep(100 * time.Millisecond)
}

func TestGetUserList(t *testing.T) {
	users, count, err := models.GetUserList(0, 10, "", false)
	result, _ := json.MarshalIndent(users, "", "    ")

	log.Info("TestGetUserList: result count %v:\n%v\nerror:%v", count, string(result), err)

	log.Info("TestGetUserList is done.\n\n")
	time.Sleep(100 * time.Millisecond)
}

func TestGetTreeMenuForUser(t *testing.T) {
	user := &proto.User{
		ID: -1,
	}

	nodes := models.GetTreeMenuForUser(user)
	result, _ := json.MarshalIndent(nodes, "", "    ")

	log.Info("TestGetTreeMenuForUser: result for super admin:\n%v", string(result))

	user.ID = 1000
	nodes = models.GetTreeMenuForUser(user)
	result, _ = json.MarshalIndent(nodes, "", "    ")

	log.Info("TestGetTreeMenuForUser: result for normal admin:\n%v", string(result))

	log.Info("TestGetTreeMenuForUser is done.\n\n")
	time.Sleep(100 * time.Millisecond)
}

func TestGetAccessTree(t *testing.T) {
	nodes, err := models.GetAccessTree(0)
	result, _ := json.MarshalIndent(nodes, "", "    ")

	log.Info("TestGetAccessTree:\n%v\nerror:%v", string(result), err)

	log.Info("TestGetAccessTree is done.\n\n")
	time.Sleep(100 * time.Millisecond)
}
