package models

import "errors"

var (
	// ErrUserNotExist is the error for user is not exists
	ErrUserNotExist = errors.New("用户不存在")
	// ErrUserDisabled is the error for user is disable
	ErrUserDisabled = errors.New("用户处于被禁用状态")
	// ErrUserWrongPwd is the error for wrong password
	ErrUserWrongPwd = errors.New("密码错误")
)
