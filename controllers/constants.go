package controllers

// session keys
const (
	cSessionUserInfoKey   = "stores_ui"
	cSessionAccessListKey = "stores_al"
)

const (
	cJSONKey  = "json"
	cUserInfo = "userinfo"

	cUsername      = "username"
	cPassword      = "password"
	cOld           = "old"
	cNew           = "new"
	cRepeat        = "repeat"
	cPage          = "page"
	cRows          = "rows"
	cSort          = "sort"
	cOrder         = "order"
	cSortDirection = "desc"
)

const (
	cActionLogin                 = "用户登录"
	cActionLoginOut              = "用户登出"
	cActionModifyPwd             = "修改登录密码"
	cActionViewUserList          = "查看仓储管理系统用户列表"
	cActionViewOperationLogsList = "查看仓储管理系统操作日志"
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
