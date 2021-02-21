package models

import (
	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models/db"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

func protoSpec2Model(entry *proto.SpecEntry) *db.StoresCommoditySpec {
	if entry == nil {
		return &db.StoresCommoditySpec{}
	}
	return &db.StoresCommoditySpec{
		Id:          entry.ID,
		SpecName:    entry.Name,
		Detail:      entry.Detail,
		ParentId:    entry.ParentID,
		Segmentable: entry.Segmentable,
		SubId:       entry.SubID,
		SubAmount:   entry.SubAmount,
	}
}

func modelSpec2Proto(spec *db.StoresCommoditySpec, session *xorm.Session, loadSub bool) *proto.Specification {
	if spec == nil {
		return &proto.Specification{}
	}
	ret := &proto.Specification{
		ID:          spec.Id,
		Name:        spec.SpecName,
		Detail:      spec.Detail,
		Segmentable: spec.Segmentable,
		SubAmount:   spec.SubAmount,
	}

	if spec.ParentId > 0 {
		record, _ := getSpec(session, spec.ParentId)
		ret.Parent = modelSpec2Proto(record, session, false)
	}
	if spec.Segmentable == 1 && loadSub {
		if session == nil {
			session = config.GetConfigs().OrmEngine.NewSession()
			defer session.Close()
		}
		ret.Sub, _ = subSpec(session, spec.Id, loadSub)
	}

	return ret
}

// AddSpec adds a spec entry
func AddSpec(entry *proto.SpecEntry) (*proto.Specification, error) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return nil, FormatMysqlError(err)
	}

	spec, err := addSpec(session, entry)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.AddSpec: rollback failed with %v", rollbackErr)
		}
		return nil, FormatMysqlError(err)
	}

	return spec, FormatMysqlError(session.Commit())
}

func addSpec(session *xorm.Session, entry *proto.SpecEntry) (*proto.Specification, error) {
	model := protoSpec2Model(entry)

	model.Id = 0
	_, err := session.Table(cTableStoresSpecification).Insert(model)
	if err != nil {
		return nil, err
	}

	if entry.ParentID > 0 {
		parent, err := getSpec(session, entry.ParentID)
		if err == proto.ErrRecordIsNotExist {
			return nil, proto.ErrParentSpecIsNotExist
		}
		if err != nil {
			return nil, err
		}

		if hasSub(parent) {
			return nil, proto.ErrParentHasChild
		}

		assignSubSpecInDB(session, parent, model, entry.Amount)
	}

	if entry.SubID > 0 {
		sub, err := getSpec(session, entry.SubID)
		if err == proto.ErrRecordIsNotExist {
			return nil, proto.ErrChildSpecIsNotExist
		}
		if err != nil {
			return nil, err
		}

		if hasParent(sub) {
			return nil, proto.ErrChildHasParent
		}

		assignSubSpecInDB(session, model, sub, entry.SubAmount)
	}

	return modelSpec2Proto(model, session, false), nil
}

func getSpec(session *xorm.Session, sid int64) (*db.StoresCommoditySpec, error) {
	spec := &db.StoresCommoditySpec{}
	found, err := session.Table(cTableStoresSpecification).Where("id=?", sid).Get(spec)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, proto.ErrRecordIsNotExist
	}

	return spec, nil
}

func hasParent(spec *db.StoresCommoditySpec) bool {
	return spec.ParentId > 0
}

func hasSub(spec *db.StoresCommoditySpec) bool {
	return spec.SubId > 0
}

func assignSubSpecInDB(session *xorm.Session, parent, sub *db.StoresCommoditySpec, subAmount int) error {
	if parent.Id <= 0 || sub.Id <= 0 {
		return proto.ErrRecordIsNotExist
	}
	if subAmount <= 0 {
		return proto.ErrCommonInvalidParam
	}

	parent.SubId = sub.Id
	parent.Segmentable = 1
	parent.SubAmount = subAmount
	_, err := session.Table(cTableStoresSpecification).
		Where("id=?", parent.Id).
		Cols("sub_id,segmentable,sub_amount").Update(parent)
	if err != nil {
		return err
	}

	sub.ParentId = parent.Id
	_, err = session.Table(cTableStoresSpecification).
		Where("id=?", sub.Id).
		Cols("parent_id").Update(sub)
	if err != nil {
		return err
	}

	return nil
}

// AssignSubSpec assigns sub to parent
func AssignSubSpec(pid, sid int64, subAmount int) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return FormatMysqlError(err)
	}

	err = assignSubSpec(session, pid, sid, subAmount)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.AssignSubSpec: rollback failed with %v", rollbackErr)
		}
		return FormatMysqlError(err)
	}

	return FormatMysqlError(session.Commit())
}

func assignSubSpec(session *xorm.Session, pid, sid int64, subAmount int) error {
	parent, err := getSpec(session, pid)
	if err != nil {
		return err
	}
	sub, err := getSpec(session, sid)
	if err != nil {
		return err
	}

	if hasSub(parent) {
		return proto.ErrParentHasChild
	}
	if hasParent(sub) {
		return proto.ErrChildHasParent
	}

	return assignSubSpecInDB(session, parent, sub, subAmount)
}

// SpecList return the paged list of specification
func SpecList(pageIndex, pageSize int,
	sort string,
	desc, loadSub bool) ([]*proto.Specification, int64, error) {

	log.Trace("entried modles.SpecList with parameters: page: %v, page size: %v, sort: %v, desc: %v, loadsub: %v",
		pageIndex, pageSize, sort, desc, loadSub)
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	return specList(session, pageIndex, pageSize, sort, desc, loadSub)
}

func specList(session *xorm.Session,
	pageIndex, pageSize int,
	sort string,
	desc, loadSub bool) ([]*proto.Specification, int64, error) {

	records := []*db.StoresCommoditySpec{}
	offset := 0
	if pageIndex > 1 {
		offset = (pageIndex - 1) * pageSize
	}
	count, err := session.Table(cTableStoresSpecification).Count()
	if err != nil {
		return nil, 0, FormatMysqlError(err)
	}

	session.Table(cTableStoresSpecification).
		Limit(pageSize, offset)
	if len(sort) > 0 {
		if desc {
			session.Desc(sort)
		} else {
			session.Asc(sort)
		}
	}

	err = session.Find(&records)
	// sql, params := session.LastSQL()
	// log.Trace("models.specList: query specification sql: `%v`, parameters: %v, error: %v", sql, params, err)

	if err != nil {
		return nil, 0, FormatMysqlError(err)
	}

	ret := make([]*proto.Specification, 0, len(records))
	for _, record := range records {
		spec := modelSpec2Proto(record, session, loadSub)

		ret = append(ret, spec)
	}

	return ret, count, nil
}

func subSpec(session *xorm.Session, parentID int64, loadSub bool) (*proto.Specification, error) {
	record := &db.StoresCommoditySpec{}
	found, err := session.Table(cTableStoresSpecification).
		Where("parent_id=?", parentID).
		Get(record)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, proto.ErrChildSpecIsNotExist
	}
	spec := modelSpec2Proto(record, session, loadSub)

	return spec, nil
}

func removeSub(session *xorm.Session, parentID int64) error {
	parent, err := getSpec(session, parentID)
	if err != nil {
		return err
	}

	if !hasSub(parent) {
		if parent.Segmentable != 0 || parent.SubAmount > 0 {
			parent.Segmentable = 0
			parent.SubAmount = 0
			_, err = session.Table(cTableStoresSpecification).
				Where("id=?", parent.Id).
				Cols("segmentable,sub_amount").Update(parent)
			return err
		}

		return nil
	}

	sub, err := getSpec(session, parent.SubId)
	if err != nil {
		return err
	}

	parent.SubId = 0
	parent.Segmentable = 0
	parent.SubAmount = 0
	_, err = session.Table(cTableStoresSpecification).
		Where("id=?", parent.Id).
		Cols("sub_id,segmentable,sub_amount").Update(parent)
	if err != nil {
		return err
	}

	sub.ParentId = parent.Id
	_, err = session.Table(cTableStoresSpecification).
		Where("id=?", sub.Id).
		Cols("parent_id").Update(sub)
	if err != nil {
		return err
	}
	return nil
}

// RemoveSub remove the sub of specified spec
func RemoveSub(parentID int64) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return FormatMysqlError(err)
	}

	err = removeSub(session, parentID)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.RemoveSub: rollback failed with %v", rollbackErr)
		}
		return FormatMysqlError(err)
	}

	return FormatMysqlError(session.Commit())
}

func isDeleteableSpec(session *xorm.Session, sid int64) (bool, error) {
	sku := &db.StoresCommoditySku{}
	found, err := session.Table(cTableStoresCommoditySku).Where("spec_id=?", sid).Get(sku)
	if err != nil {
		return false, FormatMysqlError(err)
	}
	if found {
		return false, proto.ErrDeleteSkuBeforeSpec
	}

	return true, nil
}

func deleteSpec(session *xorm.Session, sid int64) error {
	spec, err := getSpec(session, sid)
	if err != nil {
		return err
	}

	// remove paranet first
	if spec.ParentId > 0 {
		err = removeSub(session, spec.ParentId)
		if err != nil {
			return err
		}
	}

	// then remove sub
	if spec.SubId > 0 {
		err = removeSub(session, sid)
		if err != nil {
			return err
		}
	}

	// finally remove spec
	deleted, err := session.Table(cTableStoresSpecification).
		Where("id=?", sid).
		Delete(&db.StoresCommoditySpec{})
	if deleted != 1 {
		log.Error("model.DeleteSpec wants to remove %v recordes for id %v", deleted, sid)
		return proto.ErrDatabase
	}

	return err
}

// DeleteSpec deletes specified spec
func DeleteSpec(sid int64) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	var err error
	if err = session.Begin(); err != nil {
		return FormatMysqlError(err)
	}

	err = deleteSpec(session, sid)
	if err != nil {
		if rollbackErr := session.Rollback(); rollbackErr != nil {
			log.Error("models.DeleteSpec: rollback failed with %v", rollbackErr)
		}
		return FormatMysqlError(err)
	}

	return FormatMysqlError(session.Commit())
}

func updateSpec(session *xorm.Session, spec *proto.SpecEntry) error {
	model, err := getSpec(session, spec.ID)
	if err != nil {
		return err
	}

	if len(spec.Name) == 0 {
		return proto.ErrCommonInvalidParam
	}

	model.SpecName = spec.Name
	model.Detail = spec.Detail

	_, err = session.Table(cTableStoresSpecification).
		Where("id=?", spec.ID).
		Cols("spec_name,detail").Update(model)

	return FormatMysqlError(err)
}

// UpdateSpec updates specified spec's name and detail information
func UpdateSpec(spec *proto.SpecEntry) error {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	return updateSpec(session, spec)
}
