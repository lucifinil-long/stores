package models

import (
	"github.com/go-sql-driver/mysql"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

var (
	mysqlFormatedError = map[uint16]error{
		1022: proto.ErrDupKey,
		1062: proto.ErrDupKey,
		1406: proto.ErrDataTooLong,
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
		return proto.ErrDatabase
	}

	return err
}
