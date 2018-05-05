package controllers

import "github.com/lucifinil-long/stores/models/db"

// PageController hanldes page request
type PageController struct {
	BaseController
}

// AdminHomepage handles home page request
func (pc *PageController) AdminHomepage() {
	userInfo := pc.GetSession(cSessionUserInfoKey)
	if userInfo == nil {
		userInfo = db.StoresUser{}
	}

	pc.Data[cUserInfo] = userInfo
	pc.TplName = "easyui/admin/index.tpl"
}

// AdminLoginPage handles login page request
func (pc *PageController) AdminLoginPage() {
	pc.TplName = "easyui/admin/login.tpl"
}

// AdminUsersPage handles logs page request
func (pc *PageController) AdminUsersPage() {
	pc.TplName = "easyui/admin/users.tpl"
	pc.OperationLog(cActionViewUserList)
}

// AdminOperationsPage handles logs page request
func (pc *PageController) AdminOperationsPage() {
	pc.TplName = "easyui/admin/operations.tpl"
	pc.OperationLog(cActionViewOperationLogsList)
}
