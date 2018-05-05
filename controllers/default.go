package controllers

import (
	"fmt"

	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models"
	"github.com/lucifinil-long/stores/proto"
	"github.com/lucifinil-long/stores/utils"
	"github.com/mkideal/log"
)

// MainController handles homepage request
type MainController struct {
	BaseController
}

// TreeMenu handles home page tree menu request
func (mc *MainController) TreeMenu() {
	rsp := &proto.Response{}
	userinfo := mc.GetSession(cSessionUserInfoKey)
	if userinfo == nil {
		log.Trace("MainController.TreeMenu: cannot get logged in user information to session [%v]",
			mc.CruSession.SessionID())

		rsp.Status = proto.ReturnStatusTempRedirect
		rsp.Description = cTipRelogin
		rsp.Protocol = config.GetConfigs().AuthGateway()
	} else {
		user := userinfo.(proto.User)
		tree := models.GetTreeMenuForUser(&user)
		rsp.Protocol = tree
		if len(tree) == 0 {
			rsp.Status = proto.ReturnStatusFailed
			rsp.Description = cRspInteralError
		} else {
			rsp.Status = proto.ReturnStatusSuccess
			rsp.Description = cRspSuccess
		}
	}

	mc.Response(rsp)
}

// IsLoggedIn handles user login status request
func (mc *MainController) IsLoggedIn() {
	userinfo := mc.GetSession(cSessionUserInfoKey)

	configs := config.GetConfigs()
	res := &proto.IsLoggedInRes{
		Status:   false,
		Redirect: configs.AuthGateway(),
	}
	rsp := &proto.Response{
		Status:   proto.ReturnStatusSuccess,
		Protocol: res,
	}

	if userinfo != nil {
		rsp.Description = cTipLoggedin
		res.Status = true
		res.Redirect = configs.AdminHomepage()
	} else {
		rsp.Description = cTipRelogin
	}

	mc.Response(rsp)
}

// Logout handles logout request
func (mc *MainController) Logout() {
	mc.OperationLog(cActionLoginOut, cRspSuccess)

	mc.DelSession(cSessionUserInfoKey)

	mc.Response(&proto.Response{
		Status:      proto.ReturnStatusTempRedirect,
		Description: cTipRelogin,
		Protocol:    config.GetConfigs().AuthGateway(),
	})
}

// Login handles login request
func (mc *MainController) Login() {
	uid, _ := mc.GetInt64(cUID)
	password := mc.GetString(cPassword)
	user, err := CheckLogin(uid, password)
	res := &proto.LoginRes{}
	rsp := &proto.Response{
		Protocol: res,
	}
	if err == nil {
		mc.SetSession(cSessionUserInfoKey, *user)
		accesslist, _ := models.GetUserAccessList(user.ID)
		mc.SetSession(cSessionAccessListKey, accesslist)

		// update last logintime
		models.UpdateUserLoginTime(user.ID)
		mc.OperationLog(cActionLogin, cRspLoginSuccess)

		res.Redirect = config.GetConfigs().AdminHomepage()
		res.User = *user

		rsp.Status = proto.ReturnStatusSuccess
		rsp.Description = cRspLoginSuccess
		rsp.Protocol = res
	} else {
		log.Error("MainController.Login: login failed with error: %v", err)
		mc.OperationLog(cActionLogin, fmt.Sprintf("%v login failed with %v", uid, err))
		rsp.Status = proto.ReturnStatusFailed
		rsp.Description = err.Error()
	}

	mc.Response(rsp)
}

// ChangePassword handles change password request
func (mc *MainController) ChangePassword() {
	userinfo := mc.GetSession(cSessionUserInfoKey)
	if userinfo == nil {
		mc.Response(&proto.Response{
			Status:      proto.ReturnStatusTempRedirect,
			Description: cTipRelogin,
			Protocol:    config.GetConfigs().AuthGateway(),
		})
		return
	}
	oldPwd := mc.GetString(cOld)
	newPwd := mc.GetString(cNew)
	repeatPwd := mc.GetString(cRepeat)
	if newPwd != repeatPwd {
		mc.Response(&proto.Response{
			Status:      proto.ReturnStatusFailed,
			Description: cTipDifferentPwd,
			Protocol:    "",
		})
		return
	}

	user, err := CheckLogin(int64(userinfo.(proto.User).ID), oldPwd)
	if err == nil {
		err = models.UpdateUserPassword(user.ID, utils.String2MD5(newPwd))

		if err == nil {
			mc.OperationLog(cActionModifyPwd, cTipModifyPwdSuccessfully)
			mc.Response(&proto.Response{
				Status:      proto.ReturnStatusSuccess,
				Description: cTipModifyPwdSuccessfully,
				Protocol:    "",
			})
		} else {
			mc.OperationLog(cActionModifyPwd, fmt.Sprintf("failed with %v", err))
			mc.Response(&proto.Response{
				Status:      proto.ReturnStatusFailed,
				Description: err.Error(),
				Protocol:    "",
			})
		}

		return
	}

	log.Info("MainController.ChangePassword: user '%v' modified password failed with %v",
		userinfo.(proto.User).Mobile, err)
	mc.OperationLog(cActionModifyPwd, fmt.Sprintf("failed with %v", err))
	mc.Response(&proto.Response{
		Status:      proto.ReturnStatusFailed,
		Description: cTipWrongPwd,
		Protocol:    "",
	})
}

// AccessList handles access list request
func (mc *MainController) AccessList() {
	nodes, err := models.GetAccessTree(0)
	if err != nil {
		log.Error("MainController.AccessList: models.GetAccessTree returned error: %v", err)
		mc.Response(&proto.Response{
			Status:      proto.ReturnStatusFailed,
			Description: proto.ErrCommonInternalError.Error(),
			Protocol:    "",
		})
		return
	}

	mc.Response(&proto.Response{
		Status:      proto.ReturnStatusSuccess,
		Description: cRspSuccess,
		Protocol:    nodes,
	})
}
