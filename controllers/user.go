package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/lucifinil-long/stores/models"
	"github.com/lucifinil-long/stores/proto"
	"github.com/lucifinil-long/stores/utils"
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

// AddUser handles add user request
func (uc *UserController) AddUser() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.AddUserRes{},
	}

	// we should always check login authority for admin system operation even config allows anonymity login
	currentUser := uc.CurrentUser()
	if currentUser == nil {
		rsp.Status = proto.ReturnStatusNeedLogin
		rsp.Description = ErrReloginNeed.Error()
		uc.Response(rsp)
		return
	}

	insertedJSON := uc.GetString(cInserts)
	user := &proto.NewUser{}
	json.Unmarshal([]byte(insertedJSON), user)

	colMap := map[string]string{
		"Username": "用户名",
		"Nickname": "昵称",
		"Mobile":   "手机",
		"Remark":   "备注",
		"Password": "密码",
	}

	err := models.AddUser(user)
	if err != nil {
		log.Error("UserController.AddUser: models.AddUser returned '%v'", err)
		uc.OperationLog(cActionAddUser, cAddUserFailed, fmt.Sprintf(cErrorFormat, err))
		err = models.FormatMysqlError(err)
		if err == models.ErrDupKey {
			err = utils.AddErrorPrefix(err, colMap["Username"]+"'"+user.Username+"'")
		}
		rsp.Description = err.Error()
	} else {
		uc.OperationLog(cActionAddUser, cAddUserSuccess, fmt.Sprintf(cUserInfoFormat, 0, user.Username))
		rsp.Description = cAddUserSuccess
		rsp.Status = proto.ReturnStatusSuccess
	}

	uc.Response(rsp)
}
