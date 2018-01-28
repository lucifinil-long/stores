package controllers

import (
	"github.com/lucifinil-long/stores/models"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

// UserController hanldes admin system user related request
type UserController struct {
	BaseController
}

// UserList handles user list request
func (uc *UserController) UserList() {
	page, _ := uc.GetInt(cPage)
	pageSize, _ := uc.GetInt(cRows)
	sort := uc.GetString(cSort)
	order := uc.GetString(cOrder)
	desc := false
	if order == cSortDirection {
		desc = true
	}

	users, count, err := models.GetUserList(page, pageSize, sort, desc)
	res := &proto.UserListRes{
		Total: count,
		Rows:  users,
	}
	rsp := &proto.Response{
		Status:      proto.ReturnStatusSuccess,
		Description: cRspSuccess,
		Protocol:    res,
	}
	if err != nil {
		log.Error("UserController.UserList: get \"%v\" when get user list", err)
		rsp.Status = proto.ReturnStatusFailed
		rsp.Description = err.Error()
	}

	uc.Response(rsp)
}
