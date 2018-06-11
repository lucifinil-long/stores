package controllers

// session keys
const (
	cSessionUserInfoKey   = "stores_ui"
	cSessionAccessListKey = "stores_al"
)

const (
	cUpdates = "updates"
	cInserts = "inserts"
	cDeletes = "deletes"
	cUpdated = "updated"
	cInsert  = "insert"
	cReq     = "req"
)

const (
	cJSONKey  = "json"
	cUserInfo = "userinfo"

	cUID           = "uid"
	cPassword      = "password"
	cOld           = "old"
	cNew           = "new"
	cRepeat        = "repeat"
	cPage          = "page"
	cRows          = "rows"
	cSort          = "sort"
	cOrder         = "order"
	cSortDirection = "desc"
	cUserID        = "uid"
	cDepotID       = "depot_id"
)

const (
	cActionLogin                 = "用户登录"
	cActionLoginOut              = "用户登出"
	cActionModifyPwd             = "修改登录密码"
	cActionAddUser               = "添加用户"
	cActionUpdateUser            = "更新用户"
	cActionDeleteUser            = "删除用户"
	cActionViewUserList          = "查看仓储管理系统用户列表"
	cActionViewOperationLogsList = "查看仓储管理系统操作日志"
	cActionViewCommoditiesList   = "查看商品定义列表"
	cActionViewLocationList      = "查看库房列表"
)

const (
	cTipNoAuth                = "权限不足"
	cTipWrongPwd              = "密码有误"
	cTipModifyPwdSuccessfully = "密码修改成功"
	cTipDifferentPwd          = "两次输入密码不一致"
	cTipRequestFailed         = "请求失败"
	cTipRelogin               = "会话已过期，需要刷新页面重新登录!"
	cTipLoggedin              = "用户已登录"
)

const (
	cRspSuccess      = "Success"
	cRspInteralError = "服务器内部错误"
	cRspLoginSuccess = "登录成功"
)

const (
	cFailed            = "失败"
	cSuccess           = "成功"
	cAddUserFailed     = "添加用户失败"
	cAddUserSuccess    = "添加用户成功"
	cUpdateUserFailed  = "更新用户失败"
	cUpdateUserSuccess = "更新用户成功"
	cDeleteUserFailed  = "删除用户失败"
	cDeleteUserSuccess = "删除用户成功"
)

const (
	cErrorFormat    = "错误：%v"
	cUserInfoFormat = "用户ID：%v 手机： %v"
)
