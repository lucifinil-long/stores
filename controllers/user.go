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
		rsp.Description = proto.ErrReloginNeed.Error()
		uc.Response(rsp)
		return
	}

	insertedJSON := uc.GetString(cInserts)
	user := &proto.NewUser{}
	json.Unmarshal([]byte(insertedJSON), user)

	log.Info("AddUser: new user: %v, json: [%v]", user, insertedJSON)
	colMap := map[string]string{
		"Nickname": "用户名",
		"Mobile":   "手机",
		"Password": "密码",
		"Remark":   "备注",
	}
	field, err := utils.CheckDataValid(user)
	if err != nil {
		rsp.Description = utils.AddErrorPrefix(err, colMap[field]).Error()
		uc.Response(rsp)
		return
	}

	if user.Mobile <= 0 {
		rsp.Description = proto.ErrCommonInvalidParam.Error()
		uc.Response(rsp)
		return
	}
	if len(user.Password) < 5 {
		rsp.Description = proto.ErrNotEnoughPasswordStrength.Error()
		uc.Response(rsp)
		return
	}

	err = models.AddUser(user)
	if err != nil {
		log.Error("UserController.AddUser: models.AddUser returned '%v'", err)
		uc.OperationLog(cActionAddUser, cAddUserFailed, fmt.Sprintf(cErrorFormat, err))
		err = models.FormatMysqlError(err)
		if err == proto.ErrDupKey {
			err = utils.AddErrorPrefix(err, colMap["Nickname"]+"'"+user.Nickname+"'")
		}
		rsp.Description = err.Error()
	} else {
		uc.OperationLog(cActionAddUser, cAddUserSuccess, fmt.Sprintf(cUserInfoFormat, user.ID, user.Mobile))
		rsp.Description = cAddUserSuccess
		rsp.Status = proto.ReturnStatusSuccess
	}

	uc.Response(rsp)
}

// DeleteUser handles delete user request
func (uc *UserController) DeleteUser() {
	rsp := &proto.Response{
		Status:   proto.ReturnStatusFailed,
		Protocol: &proto.DeleteUserRes{},
	}

	// we should always check login authority for admin system operation even config allows anonymity login
	currentUser := uc.CurrentUser()
	if currentUser == nil {
		rsp.Status = proto.ReturnStatusNeedLogin
		rsp.Description = proto.ErrReloginNeed.Error()
		uc.Response(rsp)
		return
	}

	uid, _ := uc.GetInt(cUserID, 0)

	if currentUser.ID == uid {
		rsp.Description = proto.ErrCanNotDeleteSelf.Error()
		uc.Response(rsp)
		return
	}

	user, err := models.GetUser(uid)
	if err != nil && err != proto.ErrUserNotExist {
		rsp.Description = proto.ErrCommonInternalError.Error()
		uc.Response(rsp)
		log.Error("UserController.DeleteUser: models.GetUser returned '%v'", err)
		uc.OperationLog(cActionDeleteUser, cDeleteUserFailed, fmt.Sprintf(cErrorFormat, err))
		return
	} else if err == proto.ErrUserNotExist {
		rsp.Status = proto.ReturnStatusSuccess
		rsp.Description = cRspSuccess
		uc.Response(rsp)
		uc.OperationLog(cActionDeleteUser, cDeleteUserFailed, fmt.Sprintf(cErrorFormat, err))
		return
	}

	// if currentUser.Level > user.Level {
	// 	rsp.Description = proto.ErrCanNotUpdateHighLevelUser.Error()
	// 	uc.Response(rsp)
	// 	return
	// }

	if err = models.DeleteUser(uid); err != nil {
		rsp.Description = proto.ErrCommonInternalError.Error()
		uc.Response(rsp)
		log.Error("UserController.DeleteUser: models.DeleteUser returned '%v'", err)
		uc.OperationLog(cActionDeleteUser, cDeleteUserFailed, fmt.Sprintf(cErrorFormat, err))
		return
	}

	rsp.Status = proto.ReturnStatusSuccess
	rsp.Description = cRspSuccess
	uc.Response(rsp)

	uc.OperationLog(cActionDeleteUser, cDeleteUserSuccess,
		fmt.Sprintf(
			cUserInfoFormat,
			user.Id,
			user.Nickname,
		))
}
