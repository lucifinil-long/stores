package controllers

// session keys
const (
	cSessionUserInfoKey   = "stores_ui"
	cSessionAccessListKey = "stores_al"
)

const (
	cJSONKey  = "json"
	cUserInfo = "userinfo"
)

const (
	cActionViewUserList          = "查看仓储管理系统用户列表"
	cActionViewOperationLogsList = "查看仓储管理系统操作日志"
)

const (
	cTipNoAuth                = "权限不足"
	cTipWrongPwd              = "密码有误"
	cTipModifyPwdSuccessfully = "密码修改成功"
	cTipDifferentPwd          = "两次输入密码不一致"
	cTipRelogin               = "会话已过期，需要刷新页面重新登录!"
	cTipLoggedin              = "用户已登录"
)

const (
	cRspSuccess      = "Success"
	cRspInteralError = "服务器内部错误"
)
