package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models/db"
	"github.com/mkideal/log"

	_ "github.com/mattn/go-sqlite3"
)

func TestModels(t *testing.T) {
	log.Info("Models unittests start...")
	time.Sleep(10 * time.Millisecond)
}

// func TestSqliteEngine(t *testing.T) {
// 	cfg := config.GetConfigs()

// 	if cfg.OrmEngine == nil {
// 		log.Info("Orm Engine is nil")
// 	} else {
// 		session := cfg.OrmEngine.NewSession()
// 		defer session.Close()
// 		users := make([]db.StoresUser, 0)
// 		if err := session.Find(&users); err != nil {
// 			log.Info("Query get error %v", err)
// 		} else {
// 			log.Info("%v", users)
// 		}
// 	}

// }

func TestBuildSqlit(t *testing.T) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()
	orm, err := xorm.NewEngine("sqlite3", "file:../db/test.db?_auth&_auth_user=admin&_auth_pass=admin")
	if err != nil {
		fmt.Println("error:", err)
	} else {
		defer orm.Close()

		orm.Sync(new(db.StoresNode))
		orm.Sync(new(db.StoresUser))
		orm.Sync(new(db.StoresOpLog))
		orm.Sync(new(db.StoresRole))
		orm.Sync(new(db.StoresRoleNode))
		orm.Sync(new(db.StoresRoleUser))
		orm.Sync(new(db.StoresLocationDepot))
		orm.Sync(new(db.StoresLocationShelf))
		orm.Sync(new(db.StoresCommodity))
		orm.Sync(new(db.StoresCommoditySku))
		orm.Sync(new(db.StoresCommoditySpec))
		orm.Sync(new(db.StoresSkuProperty))
		orm.Sync(new(db.StoresSkuPropertyValue))
		orm.Sync(new(db.StoresSkuStock))
		orm.Sync(new(db.StoresSkuStockChange))

		sqliteSess := orm.NewSession()
		defer sqliteSess.Close()

		nodes := make([]*db.StoresNode, 0)
		session.Find(&nodes)
		sqliteSess.InsertMulti(nodes)

		users := make([]*db.StoresUser, 0)
		session.Find(&users)
		sqliteSess.InsertMulti(users)

		roles := make([]*db.StoresRole, 0)
		session.Find(&roles)
		sqliteSess.InsertMulti(roles)
	}
}

// func TestAddSpec(t *testing.T) {
// 	specs, count, err := models.SpecList(0, 100, "", false, false)
// 	newSpec := &proto.SpecEntry{
// 		Name: fmt.Sprintf("%v", count+1),
// 	}

// 	if err != nil {
// 		log.Fatal("models.SpecList failed with %v", err)
// 	}

// 	spec, err := models.AddSpec(newSpec)
// 	if err != nil {
// 		log.Fatal("models.AddSpec failed with %v", err)
// 	}
// 	result, _ := json.MarshalIndent(spec, "", "    ")
// 	log.Info("models.AddSpec returned %v, error: %v", string(result), err)

// 	newSpec = &proto.SpecEntry{
// 		Name:     fmt.Sprintf("%v", count+2),
// 		ParentID: spec.ID,
// 		Amount:   10,
// 	}

// 	spec, err = models.AddSpec(newSpec)
// 	if err != nil {
// 		log.Fatal("models.AddSpec failed with %v", err)
// 	}

// 	specs, count, err = models.SpecList(0, 100, "", false, true)
// 	if err != nil {
// 		log.Fatal("models.SpecList failed with %v", err)
// 	}

// 	result, _ = json.MarshalIndent(specs, "", "    ")
// 	log.Info("models.SpecList returned %v, count: %v, error: %v", string(result), count, err)

// 	time.Sleep(100 * time.Millisecond)
// }

// func TestGetUserInfoByUserIDOrMobile(t *testing.T) {
// 	users, err := models.GetUserInfoByUserIDOrMobile(-1)
// 	result, _ := json.MarshalIndent(users, "", "    ")

// 	log.Info("TestGetUserInfoByUserIDOrMobile: result for user id:\n%v\nerror:%v", string(result), err)

// 	users, err = models.GetUserInfoByUserIDOrMobile(12345678901)
// 	result, _ = json.MarshalIndent(users, "", "    ")

// 	log.Info("TestGetUserInfoByUserIDOrMobile: result for user mobile:\n%v\nerror:%v", string(result), err)

// 	log.Info("TestGetUserInfoByUserIDOrMobile is done.\n\n")
// 	time.Sleep(100 * time.Millisecond)
// }

// func TestGetUserAccessList(t *testing.T) {
// 	accesses, err := models.GetUserAccessList(-1)
// 	log.Info("TestGetUserAccessList: result for super admin:\n%v\nerror:%v", accesses, err)
// 	accesses, err = models.GetUserAccessList(1000)
// 	log.Info("TestGetUserAccessList: result for super admin:\n%v\nerror:%v", accesses, err)

// 	log.Info("TestGetUserAccessList is done.\n\n")
// 	time.Sleep(100 * time.Millisecond)
// }

// func TestGetUserList(t *testing.T) {
// 	users, count, err := models.GetUserList(0, 10, "", false)
// 	result, _ := json.MarshalIndent(users, "", "    ")

// 	log.Info("TestGetUserList: result count %v:\n%v\nerror:%v", count, string(result), err)

// 	log.Info("TestGetUserList is done.\n\n")
// 	time.Sleep(100 * time.Millisecond)
// }

// func TestGetTreeMenuForUser(t *testing.T) {
// 	user := &proto.User{
// 		ID: -1,
// 	}

// 	nodes := models.GetTreeMenuForUser(user)
// 	result, _ := json.MarshalIndent(nodes, "", "    ")

// 	log.Info("TestGetTreeMenuForUser: result for super admin:\n%v", string(result))

// 	user.ID = 1000
// 	nodes = models.GetTreeMenuForUser(user)
// 	result, _ = json.MarshalIndent(nodes, "", "    ")

// 	log.Info("TestGetTreeMenuForUser: result for normal admin:\n%v", string(result))

// 	log.Info("TestGetTreeMenuForUser is done.\n\n")
// 	time.Sleep(100 * time.Millisecond)
// }

// func TestGetAccessTree(t *testing.T) {
// 	nodes, err := models.GetAccessTree(0)
// 	result, _ := json.MarshalIndent(nodes, "", "    ")

// 	log.Info("TestGetAccessTree:\n%v\nerror:%v", string(result), err)

// 	log.Info("TestGetAccessTree is done.\n\n")
// 	time.Sleep(100 * time.Millisecond)
// }

// func TestAddDepot(t *testing.T) {
// 	depot := &proto.Depot{
// 		ID:     5,
// 		Name:   "1号仓库",
// 		Detail: "1号仓库",
// 		Shelfs: []proto.Shelf{
// 			proto.Shelf{
// 				ID:     22,
// 				Name:   "1号货架",
// 				Layers: 5,
// 				Detail: "1号货架",
// 			},
// 		},
// 	}
// 	err := models.AddDepot(depot)

// 	log.Info("models.AddDepot is done with error %v.\n\n", err)
// 	time.Sleep(100 * time.Millisecond)
// }

// func TestGetAllDepotsAndUpdate(t *testing.T) {
// 	records, count, err := models.GetDepots(0, 50, "", false)

// 	for _, record := range records {
// 		record.Detail = fmt.Sprintf("%v %v", record.Name, time.Now())
// 		err = models.UpdateDepotProperties(&record)
// 		log.Info("models.UpdateDepotProperties is done with error %v.\n\n", err)
// 	}

// 	records, count, err = models.GetDepots(0, 50, "", false)
// 	result, _ := json.MarshalIndent(records, "", "    ")

// 	log.Info("models.GetDepots: %v records in first page\n%v\nerror: %v", count, string(result), err)

// 	log.Info("TestGetAllDepotsAndUpdate is done.\n\n")

// 	time.Sleep(100 * time.Millisecond)
// }

// func TestDeleteDepots(t *testing.T) {
// 	records, err := models.GetAllDepots()

// 	ids := make([]int64, 0, len(records))
// 	for _, record := range records {
// 		ids = append(ids, record.ID)
// 	}

// 	err = models.DeleteDepots(ids)

// 	log.Info("models.DeleteDepots(%v) is done with error (%v).\n\n", ids, err)

// 	time.Sleep(100 * time.Millisecond)
// }

// func TestAddShelfs(t *testing.T) {
// 	Shelfs := []proto.Shelf{
// 		proto.Shelf{
// 			ID:     22,
// 			Name:   "2号货架",
// 			Layers: 6,
// 			Detail: "2号货架",
// 		},
// 		proto.Shelf{
// 			ID:     22,
// 			Name:   "3号货架",
// 			Layers: 6,
// 			Detail: "3号货架",
// 		},
// 	}
// 	err := models.AddShelfs(3, Shelfs)

// 	log.Info("models.AddShelfs is done with error %v.\n\n", err)
// 	time.Sleep(100 * time.Millisecond)
// }

// func TestGetAllDepotsAndUpdate(t *testing.T) {
// 	records, err := models.GetAllDepots()

// 	if len(records) > 0 {
// 		shelfs := records[0].Shelfs
// 		for _, shelf := range shelfs {
// 			shelf.Detail = fmt.Sprintf("%v %v", shelf.Name, time.Now())
// 			shelf.Layers = shelf.Layers + 1
// 			err := models.UpdateShelf(shelf)
// 			log.Info("models.UpdateShelf is done with error %v.\n\n", err)
// 		}
// 	}

// 	records, err = models.GetAllDepots()
// 	result, _ := json.MarshalIndent(records, "", "    ")

// 	log.Info("Updated depot:\n%v\nerror: %v", string(result), err)

// 	log.Info("models.GetAllDepots is done.\n\n")

// 	time.Sleep(100 * time.Millisecond)
// }

// func TestDeleteShelfs(t *testing.T) {
// 	records, err := models.GetAllDepots()

// 	ids := make([]int64, 0, len(records))
// 	for _, record := range records {
// 		for _, shelf := range record.Shelfs {

// 			ids = append(ids, shelf.ID)
// 		}
// 	}

// 	err = models.DeleteShelfs(ids)

// 	log.Info("models.DeleteShelfs(%v) is done with error (%v).\n\n", ids, err)

// 	time.Sleep(100 * time.Millisecond)
// }

func TestModelsDone(t *testing.T) {
	log.Info("All tests for models are done.")
}
