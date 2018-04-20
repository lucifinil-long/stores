package proto

import "errors"

var (
	// ErrReloginNeed is the error when user session is expired
	ErrReloginNeed = errors.New("用户会话已过期，需要重新登陆才能执行该操作")
	// ErrNotEnoughPasswordStrength indicates the password strength is not enough
	ErrNotEnoughPasswordStrength = errors.New("密码强度不够，请使用至少5位长度以上，并且包含特殊字符和数字及大小写混合的密码")

	// ErrUserNotExist is the error for user is not exists
	ErrUserNotExist = errors.New("用户不存在")
	// ErrUserDisabled is the error for user is disable
	ErrUserDisabled = errors.New("用户处于被禁用状态")
	// ErrUserWrongPwd is the error for wrong password
	ErrUserWrongPwd = errors.New("密码错误")
	// ErrCanNotDeleteSelf is the error for in case user wants to delete self
	ErrCanNotDeleteSelf = errors.New("用户不能自己删除自己")
	// ErrCanNotUpdateHighLevelUser is the error for low level user wants to update high level user's information
	ErrCanNotUpdateHighLevelUser = errors.New("低级用户无权更新高级用户信息")
)

var (
	// ErrCommonInvalidParam is common invalid parameter error
	ErrCommonInvalidParam = errors.New("错误参数")
	// ErrCommonInternalError is common internal error
	ErrCommonInternalError = errors.New("内部错误")
)

var (
	// ErrDatabase is common database error
	ErrDatabase = errors.New("数据库操作错误。")
	// ErrDupKey is duplicate key record error
	ErrDupKey = errors.New("已经存在相同关键字的记录。")
	// ErrDataTooLong is data too long error
	ErrDataTooLong = errors.New("内容长度超过字段长度限制。")
	// ErrNotStructType is error for object is not struct type
	ErrNotStructType = errors.New("非结构类型。")
	// ErrInvalidValueForNutNullField is the error for there is nil for not null field
	ErrInvalidValueForNutNullField = errors.New("非空字段必须填写有效值。")
)
