package models

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/mkideal/log"
)

var (
	// ErrCommonInvalidParam is common invalid parameter error
	ErrCommonInvalidParam = errors.New("错误参数")
	// ErrCommonInternalError is common internal error
	ErrCommonInternalError = errors.New("内部错误")
)

var (
	// ErrUserNotExist is the error for user is not exists
	ErrUserNotExist = errors.New("用户不存在")
	// ErrUserDisabled is the error for user is disable
	ErrUserDisabled = errors.New("用户处于被禁用状态")
	// ErrUserWrongPwd is the error for wrong password
	ErrUserWrongPwd = errors.New("密码错误")
	// ErrDatabase is common database error
	ErrDatabase = errors.New("数据库操作错误。")
	// ErrDupKey is duplicate key record error
	ErrDupKey = errors.New("已经存在相同关键字的记录。")
	// ErrDataTooLong is data too long error
	ErrDataTooLong = errors.New("内容长度超过字段长度限制。")
)

var (
	mysqlFormatedError = map[uint16]error{
		1022: ErrDupKey,
		1062: ErrDupKey,
		1406: ErrDataTooLong,
	}
)

// FormatMysqlError formarts mysql error
func FormatMysqlError(err error) error {
	if err == nil {
		return nil
	}

	if sqlErr, ok := err.(*mysql.MySQLError); ok {
		if err, ok = mysqlFormatedError[sqlErr.Number]; ok {
			return err
		}

		log.Info("models.formatMysqlError get '%v', return common database error", sqlErr)
		return ErrDatabase
	}

	return err
}
