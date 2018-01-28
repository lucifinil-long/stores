package controllers

import "github.com/lucifinil-long/stores/models/db"

// PageController hanldes page request
type PageController struct {
	BaseController
}

// Homepage handles home page request
func (pc *PageController) Homepage() {
	userInfo := pc.GetSession(cSessionUserInfoKey)
	if userInfo == nil {
		userInfo = db.StoresUser{}
	}

	pc.Data[cUserInfo] = userInfo
	pc.TplName = "easyui/public/index.tpl"
}

// LoginPage handles login page request
func (pc *PageController) LoginPage() {
	pc.TplName = "easyui/public/login.tpl"
}

// AdminUsersPage handles logs page request
func (pc *PageController) AdminUsersPage() {
	pc.TplName = "easyui/admin/user.tpl"
	pc.OperationLog(cActionViewUserList)
}

// AdminOperationsPage handles logs page request
func (pc *PageController) AdminOperationsPage() {
	pc.TplName = "easyui/admin/operations.tpl"
	pc.OperationLog(cActionViewOperationLogsList)
}
