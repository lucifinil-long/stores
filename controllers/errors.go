package controllers

import "errors"

var (
	// ErrReloginNeed is the error when user session is expired
	ErrReloginNeed = errors.New("用户会话已过期，需要重新登陆才能执行该操作")
)
