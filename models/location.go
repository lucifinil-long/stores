package models

import (
	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models/db"
	"github.com/lucifinil-long/stores/proto"
	"github.com/lucifinil-long/stores/utils"
	"github.com/mkideal/log"
)

// Depot is stores location depot struct
type Depot struct {
	db.StoresLocationDepot `xorm:"extends"`
	Shelfs                 []*db.StoresLocationShelf `json:"shelfs" xorm:"-"`
}

func (depot Depot) depot2Proto() proto.Depot {
	shelfs := make([]proto.Shelf, 0, len(depot.Shelfs))
	for _, shelf := range depot.Shelfs {
		shelfs = append(shelfs, dbShelf2Proto(shelf))
	}
	return proto.Depot{
		ID:     depot.Id,
		Name:   depot.Name,
		Detail: depot.Detail,
		Shelfs: shelfs,
	}
}

func (depot *Depot) updateDepotIDOfShelfs() {
	for _, shelf := range depot.Shelfs {
		shelf.DepotId = depot.Id
	}
}

func dbShelf2Proto(shelf *db.StoresLocationShelf) proto.Shelf {
	if shelf == nil {
		return proto.Shelf{}
	}
	return proto.Shelf{
		ID:     shelf.Id,
		Name:   shelf.Name,
		Layers: shelf.Layers,
		Detail: shelf.Detail,
	}
}

func protoShelf2DB(shelf proto.Shelf, needID bool) *db.StoresLocationShelf {
	if needID {
		return &db.StoresLocationShelf{
			Id:     shelf.ID,
			Name:   shelf.Name,
			Layers: shelf.Layers,
			Detail: shelf.Detail,
		}
	}
	return &db.StoresLocationShelf{
		Name:   shelf.Name,
		Layers: shelf.Layers,
		Detail: shelf.Detail,
	}
}

func protoDepot2Model(depot *proto.Depot, needID bool) *Depot {
	if depot == nil {
		return &Depot{
			Shelfs: []*db.StoresLocationShelf{},
		}
	}

	shelfs := make([]*db.StoresLocationShelf, 0, len(depot.Shelfs))
	for _, shelf := range depot.Shelfs {
		shelfs = append(shelfs, protoShelf2DB(shelf, needID))
	}

	ret := &Depot{
		StoresLocationDepot: db.StoresLocationDepot{
			Name:   depot.Name,
			Detail: depot.Detail,
		},
		Shelfs: shelfs,
	}

	if needID {
		ret.Id = depot.ID
	}

	ret.updateDepotIDOfShelfs()
	return ret
}

// GetDepots get depots
func GetDepots(pageIndex, pageSize int, sort string, desc bool, loadShelfs bool) ([]proto.Depot, int64, error) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	depots, count, err := getDepots(session, pageIndex, pageSize, sort, desc, loadShelfs)

	if err != nil {
		return nil, 0, err
	}

	records := make([]proto.Depot, 0, len(depots))
	for _, depot := range depots {
		records = append(records, depot.depot2Proto())
	}

	return records, count, nil
}

func getDepots(session *xorm.Session, pageIndex, pageSize int, sort string, desc bool, loadShelfs bool) ([]*Depot, int64, error) {
	depots := make([]*Depot, 0)
	offset := 0
	if pageIndex > 1 {
		offset = (pageIndex - 1) * pageSize
	}
	count, err := session.Table(cTableStoresLocationDepot).Count()
	if err != nil {
		return nil, 0, FormatMysqlError(err)
	}

	session.Table(cTableStoresLocationDepot).
		Limit(pageSize, offset)
	if len(sort) > 0 {
		if desc {
			session.Desc(sort)
		} else {
			session.Asc(sort)
		}
	}

	err = session.Find(&depots)
	// sql, params := session.LastSQL()
	// log.Trace("models.getAllDepots: query depots sql: `%v`, parameters: %v, error: %v", sql, params, err)

	if err != nil {
		return nil, 0, FormatMysqlError(err)
	}

	if loadShelfs {
		for _, depot := range depots {
			depot.Shelfs = make([]*db.StoresLocationShelf, 0)
			err = session.Table(cTableStoresLocationShelf).Where("depot_id=?", depot.Id).Find(&depot.Shelfs)
			// sql, params := session.LastSQL()
			// log.Trace("models.getAllDepots: query shelfs for depot(id:%v) sql: `%v`, parameters: %v, error: %v",
			// 	depot.Id, sql, params, err)

			if err != nil {
				return nil, 0, FormatMysqlError(err)
			}
		}
	}

	return depots, count, nil
}

// AddDepot adds a depot
func AddDepot(depot *proto.Depot) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return FormatMysqlError(err)
	}

	err = addDepot(session, depot)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.AddDepot: rollback failed with %v", rollbackErr)
		}
		return err
	}

	return FormatMysqlError(session.Commit())
}

func addDepot(session *xorm.Session, depot *proto.Depot) error {
	model := protoDepot2Model(depot, false)
	count, err := session.Table(cTableStoresLocationDepot).Insert(&model.StoresLocationDepot)
	if err != nil {
		return err
	}

	if count < 1 {
		return proto.ErrDatabase
	}

	model.updateDepotIDOfShelfs()

	_, err = session.Table(cTableStoresLocationShelf).InsertMulti(&model.Shelfs)

	return FormatMysqlError(err)
}

// UpdateDepotProperties update depot properties
func UpdateDepotProperties(depot *proto.Depot) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	return updateDepotProperties(session, depot)
}

func updateDepotProperties(session *xorm.Session, depot *proto.Depot) error {
	model := protoDepot2Model(depot, true)
	_, err := session.Table(cTableStoresLocationDepot).
		Where("id=?", depot.ID).
		Cols("name,detail").
		Update(&model.StoresLocationDepot)

	return FormatMysqlError(err)
}

// DeleteDepots delete depots
func DeleteDepots(ids []int64) error {
	if len(ids) == 0 {
		return proto.ErrCommonInvalidParam
	}
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return FormatMysqlError(err)
	}

	err = deleteDepots(session, ids)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.DeleteDepots: rollback failed with %v", rollbackErr)
		}
		return err
	}

	return FormatMysqlError(session.Commit())
}

func deleteDepots(session *xorm.Session, ids []int64) error {
	params := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		params = append(params, id)
	}

	amount, err := getStockForDepots(session, params)
	if err != nil {
		return err
	}
	if amount > 0 {
		return proto.ErrCanNotDeleteDepotsForStock
	}

	_, err = session.Table(cTableStoresLocationShelf).
		In("depot_id", params...).
		Delete(&db.StoresLocationShelf{})
	sql, params := session.LastSQL()
	log.Trace("models.deleteDepots: delete shelf sql: `%v`, parameters: %v, error: %v", sql, params, err)

	if err != nil {
		return FormatMysqlError(err)
	}

	_, err = session.Table(cTableStoresLocationDepot).
		In("id", params...).
		Delete(&db.StoresLocationDepot{})
	sql, params = session.LastSQL()
	log.Trace("models.deleteDepots: delete depots sql: `%v`, parameters: %v, error: %v", sql, params, err)

	return FormatMysqlError(err)
}

func getStockForDepots(session *xorm.Session, ids []interface{}) (int64, error) {
	param := utils.JoinArray(",", ids...)
	result, err := session.Query("SELECT sum(amount) as amount FROM stores_sku_stock where shelf_id in (select id from stores_location_shelf where depot_id in (?))",
		param)

	if err != nil {
		return 0, FormatMysqlError(err)
	}

	amount := int64(0)
	if 1 == len(result) {
		amount, _ = strconv.ParseInt(string(result[0]["amount"][:]), 10, 64)
	}
	return amount, nil
}

// GetShelfsOfDepot get shelfs of specified depot
func GetShelfsOfDepot(depotID int64, pageIndex, pageSize int, sort string, desc bool) ([]proto.Shelf, int64, error) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	shelfs, count, err := getShelfs(session, depotID, pageIndex, pageSize, sort, desc)

	if err != nil {
		return nil, 0, err
	}

	records := make([]proto.Shelf, 0, len(shelfs))
	for _, shelf := range shelfs {
		records = append(records, dbShelf2Proto(shelf))
	}

	return records, count, nil
}

func getShelfs(session *xorm.Session, depotID int64, pageIndex, pageSize int, sort string, desc bool) ([]*db.StoresLocationShelf, int64, error) {
	shelfs := make([]*db.StoresLocationShelf, 0)
	offset := 0
	if pageIndex > 1 {
		offset = (pageIndex - 1) * pageSize
	}
	count, err := session.Table(cTableStoresLocationShelf).Where("depot_id=?", depotID).Count()
	if err != nil {
		return nil, 0, FormatMysqlError(err)
	}

	session.Table(cTableStoresLocationShelf).
		Limit(pageSize, offset)
	if len(sort) > 0 {
		if desc {
			session.Desc(sort)
		} else {
			session.Asc(sort)
		}
	}

	err = session.Where("depot_id=?", depotID).Find(&shelfs)
	// sql, params := session.LastSQL()
	// log.Trace("models.GetShelfsOfDepot: query shelfs sql: `%v`, parameters: %v, error: %v", sql, params, err)

	if err != nil {
		return nil, 0, FormatMysqlError(err)
	}

	return shelfs, count, nil
}

// AddShelfs adds shelfs to specified depot
func AddShelfs(depotID int64, shelfs []proto.Shelf) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return FormatMysqlError(err)
	}

	err = addShelfs(session, depotID, shelfs)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.AddShelfs: rollback failed with %v", rollbackErr)
		}
		return err
	}

	return FormatMysqlError(session.Commit())
}

func addShelfs(session *xorm.Session, depotID int64, shelfs []proto.Shelf) error {
	record := &db.StoresLocationDepot{}
	has, err := session.Table(cTableStoresLocationDepot).Where("id=?", depotID).Get(record)
	if err != nil {
		return FormatMysqlError(err)
	}
	if !has {
		return proto.ErrDepotIsNotExist
	}

	protoDepot := &proto.Depot{
		Shelfs: shelfs,
	}
	model := protoDepot2Model(protoDepot, false)
	model.Id = depotID
	model.updateDepotIDOfShelfs()
	_, err = session.Table(cTableStoresLocationShelf).InsertMulti(&model.Shelfs)

	return FormatMysqlError(err)
}

// UpdateShelf update shelf infomation
func UpdateShelf(shelf proto.Shelf) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	return updateShelf(session, shelf)
}

func updateShelf(session *xorm.Session, shelf proto.Shelf) error {
	if len(shelf.Name) == 0 {
		return proto.ErrEmptyShelfName
	}
	model := protoShelf2DB(shelf, true)
	_, err := session.Table(cTableStoresLocationShelf).
		Cols("name,layers,detail").
		Where("id=?", shelf.ID).
		Update(model)
	return FormatMysqlError(err)
}

// DeleteShelfs delete shelfs
func DeleteShelfs(ids []int64) error {
	if len(ids) == 0 {
		return proto.ErrCommonInvalidParam
	}
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return FormatMysqlError(err)
	}

	err = deleteShelfs(session, ids)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.DeleteShelfs: rollback failed with %v", rollbackErr)
		}
		return err
	}

	return FormatMysqlError(session.Commit())
}

func deleteShelfs(session *xorm.Session, ids []int64) error {
	params := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		params = append(params, id)
	}

	amount, err := getStockForShelfs(session, params)
	if err != nil {
		return err
	}
	if amount > 0 {
		return proto.ErrCanNotDeleteShelfsForStock
	}

	_, err = session.Table(cTableStoresLocationShelf).
		In("id", params...).
		Delete(&db.StoresLocationShelf{})
	sql, params := session.LastSQL()
	log.Trace("models.deleteShelfs: delete shelf sql: `%v`, parameters: %v, error: %v", sql, params, err)

	return FormatMysqlError(err)
}

func getStockForShelfs(session *xorm.Session, ids []interface{}) (int64, error) {
	param := utils.JoinArray(",", ids...)
	result, err := session.Query("SELECT sum(amount) as amount FROM stores_sku_stock where shelf_id in (?)",
		param)

	if err != nil {
		return 0, FormatMysqlError(err)
	}

	amount := int64(0)
	if 1 == len(result) {
		amount, _ = strconv.ParseInt(string(result[0]["amount"][:]), 10, 64)
	}
	return amount, nil
}
