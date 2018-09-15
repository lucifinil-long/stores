package controllers

import (
	"encoding/json"

	"github.com/lucifinil-long/stores/models"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

// LocationController hanldes localtion related request
type LocationController struct {
	BaseController
}

// DepotList handles depot list request
func (lc *LocationController) DepotList() {
	page, _ := lc.GetInt(cPage)
	pageSize, _ := lc.GetInt(cRows)
	sort := lc.GetString(cSort)
	order := lc.GetString(cOrder)
	desc := false
	if order == cSortDirection {
		desc = true
	}

	depots, count, err := models.GetDepots(page, pageSize, sort, desc, true)
	res := &proto.DepotsListRes{
		Total: count,
		Rows:  depots,
	}
	rsp := &proto.Response{
		Status:      proto.ReturnStatusSuccess,
		Description: cRspSuccess,
		Protocol:    res,
	}
	if err != nil {
		log.Error("LocaltionController.DepotList: get \"%v\" when get user list", err)
		rsp.Status = proto.ReturnStatusFailed
		rsp.Description = err.Error()
	}

	lc.Response(rsp)
}

// DeleteDepot handles depot list request
func (lc *LocationController) DeleteDepot() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.DeleteDepotRes{},
	}

	deletedJSON := lc.Input().Get(cDeletes)
	deletes := []int64{}
	json.Unmarshal([]byte(deletedJSON), &deletes)

	if err := models.DeleteDepots(deletes); err != nil {
		rsp.Description = err.Error()
		lc.Response(rsp)
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	lc.Response(rsp)
}

// AddDepot handles depot list request
func (lc *LocationController) AddDepot() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.AddDepotRes{},
	}
	insertedJSON := lc.GetString(cInsert)
	depot := &proto.Depot{}
	json.Unmarshal([]byte(insertedJSON), depot)

	if len(depot.Name) == 0 {
		rsp.Description = proto.ErrCommonInvalidParam.Error()
		lc.Response(rsp)
		return
	}

	if err := models.AddDepot(depot); err != nil {
		rsp.Description = err.Error()
		lc.Response(rsp)
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	lc.Response(rsp)
}

// UpdateDepot handles depot list request
func (lc *LocationController) UpdateDepot() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.UpdateDepotRes{},
	}
	updateJSON := lc.GetString(cUpdated)
	depot := &proto.Depot{}
	json.Unmarshal([]byte(updateJSON), depot)

	if len(depot.Name) == 0 {
		rsp.Description = proto.ErrCommonInvalidParam.Error()
		lc.Response(rsp)
		return
	}

	if err := models.UpdateDepotProperties(depot); err != nil {
		rsp.Description = err.Error()
		lc.Response(rsp)
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	lc.Response(rsp)
}

// ShelfList handles shelf list request
func (lc *LocationController) ShelfList() {
	depotID, _ := lc.GetInt64(cDepotID)
	page, _ := lc.GetInt(cPage)
	pageSize, _ := lc.GetInt(cRows)
	sort := lc.GetString(cSort)
	order := lc.GetString(cOrder)
	desc := false
	if order == cSortDirection {
		desc = true
	}

	shelfs, count, err := models.GetShelfsOfDepot(depotID, page, pageSize, sort, desc)
	res := &proto.ShelfsListRes{
		Total: count,
		Rows:  shelfs,
	}
	rsp := &proto.Response{
		Status:      proto.ReturnStatusSuccess,
		Description: cRspSuccess,
		Protocol:    res,
	}
	if err != nil {
		log.Error("LocaltionController.DepotList: get \"%v\" when get user list", err)
		rsp.Status = proto.ReturnStatusFailed
		rsp.Description = err.Error()
	}

	lc.Response(rsp)
}

// DeleteShelf handles delete shelfs request
func (lc *LocationController) DeleteShelf() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.DeleteDepotRes{},
	}

	deletedJSON := lc.Input().Get(cDeletes)
	deletes := []int64{}
	json.Unmarshal([]byte(deletedJSON), &deletes)

	if err := models.DeleteShelfs(deletes); err != nil {
		rsp.Description = err.Error()
		lc.Response(rsp)
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	lc.Response(rsp)
}

// AddShelfs handles shelf add request
func (lc *LocationController) AddShelfs() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.AddShelfsRes{},
	}
	insertedJSON := lc.GetString(cReq)
	req := &proto.AddShelfsReq{}
	json.Unmarshal([]byte(insertedJSON), req)

	if err := models.AddShelfs(req.DepotID, req.Shelfs); err != nil {
		rsp.Description = err.Error()
		lc.Response(rsp)
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	lc.Response(rsp)
}

// UpdateShelf handles shelf update request
func (lc *LocationController) UpdateShelf() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.UpdateShelfRes{},
	}
	insertedJSON := lc.GetString(cUpdated)
	shelf := &proto.Shelf{}
	json.Unmarshal([]byte(insertedJSON), shelf)

	if len(shelf.Name) == 0 {
		rsp.Description = proto.ErrCommonInvalidParam.Error()
		lc.Response(rsp)
		return
	}

	if err := models.UpdateShelf(*shelf); err != nil {
		rsp.Description = err.Error()
		lc.Response(rsp)
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	lc.Response(rsp)
}
