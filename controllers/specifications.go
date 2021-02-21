package controllers

import (
	"encoding/json"

	"github.com/lucifinil-long/stores/models"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

// SpecificationController hanldes specification related request
type SpecificationController struct {
	BaseController
}

// SpecList handles specification list request
func (sc *SpecificationController) SpecList() {
	page, _ := sc.GetInt(cPage)
	pageSize, _ := sc.GetInt(cRows)
	sort := sc.GetString(cSort)
	order := sc.GetString(cOrder)
	desc := false
	if order == cSortDirection {
		desc = true
	}

	records, count, err := models.SpecList(page, pageSize, sort, desc, true)
	specs := []proto.Specification{}
	for _, record := range records {
		specs = append(specs, *record)
	}
	res := &proto.SpecificationListRes{
		Total: count,
		Rows:  specs,
	}
	rsp := &proto.Response{
		Status:      proto.ReturnStatusSuccess,
		Description: cRspSuccess,
		Protocol:    res,
	}
	if err != nil {
		log.Error("SpecificationController.SpecList: get \"%v\" when get user list", err)
		rsp.Status = proto.ReturnStatusFailed
		rsp.Description = err.Error()
	}

	sc.Response(rsp)
}

// AddSpec handles add specification request
func (sc *SpecificationController) AddSpec() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.AddSpecificationRes{},
	}

	insertedJSON := sc.GetString(cInsert)
	spec := &proto.SpecEntry{}
	json.Unmarshal([]byte(insertedJSON), spec)

	if len(spec.Name) == 0 {
		rsp.Description = proto.ErrCommonInvalidParam.Error()
		sc.Response(rsp)
		return
	}

	if _, err := models.AddSpec(spec); err != nil {
		rsp.Description = err.Error()
		sc.Response(rsp)
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	sc.Response(rsp)
}

// UpdateSpec handles specification update request
func (sc *SpecificationController) UpdateSpec() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.UpdateSpecificationRes{},
	}
	updateJSON := sc.GetString(cUpdated)
	spec := &proto.SpecEntry{}
	json.Unmarshal([]byte(updateJSON), spec)

	if len(spec.Name) == 0 {
		rsp.Description = proto.ErrCommonInvalidParam.Error()
		sc.Response(rsp)
		return
	}

	if err := models.UpdateSpec(spec); err != nil {
		rsp.Description = err.Error()
		sc.Response(rsp)
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	sc.Response(rsp)
}
