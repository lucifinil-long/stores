package models

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models/db"
	"github.com/mkideal/log"
)

// AddOperationLog adds a operation log entry to db
// @param uid is the id of user that makes operation
// @param from is the address that user is from
// @param nickname is the nickname of user that makes operation
// @param action is action summary
// @param detail is action detail
func AddOperationLog(uid int64, nickname, from, action, detail string) {
	session := config.GetConfigs().OrmEngine.NewSession()
	defer session.Close()

	addOperationLog(session, uid, nickname, from, action, detail)
}

// addOperationLog adds a operation log entry to db
// @param session is database session, can be nil; if nil will use default database session
// @param uid is the id of user that makes operation
// @param from is the address that user is from
// @param nickname is the nickname of user that makes operation
// @param action is action summary
// @param detail is action detail
func addOperationLog(session *xorm.Session, uid int64, nickname, from, action, detail string) {
	if session == nil {
		session = config.GetConfigs().OrmEngine.NewSession()
		defer session.Close()
	}

	entry := db.StoresOpLog{
		UserId:      uid,
		Nickname:    nickname,
		From:        from,
		Action:      action,
		Detail:      detail,
		CreatedTime: time.Now()}

	if inserted, err := session.Table(entry).Insert(entry); err != nil || inserted != 1 {
		log.Error("models.AddOperationLog: failed to save operation log!! error: %v, inserted: %v", err, inserted)
	}

}
