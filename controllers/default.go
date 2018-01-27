package controllers

import (
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models"
	"github.com/lucifinil-long/stores/proto"
	"github.com/mkideal/log"
)

// MainController handles homepage request
type MainController struct {
	BaseController
}

// TreeMenu handles home page tree menu request
func (mc *MainController) TreeMenu() {
	rsp := &proto.Response{}
	userinfo := mc.GetSession(cSessionUserInfoKey)
	if userinfo == nil {
		log.Trace("MainController.TreeMenu: cannot get logged in user information to session [%v]",
			mc.CruSession.SessionID())

		rsp.Status = proto.ReturnStatusTempRedirect
		rsp.Description = cTipRelogin
		rsp.Protocol = config.GetConfigs().AuthGateway()
	} else {
		user := userinfo.(proto.User)
		tree := models.GetTreeMenuForUser(&user)
		rsp.Protocol = tree
		if len(tree) == 0 {
			rsp.Status = proto.ReturnStatusFailed
			rsp.Description = cRspInteralError
		} else {
			rsp.Status = proto.ReturnStatusSuccess
			rsp.Description = cRspSuccess
		}
	}

	mc.Response(rsp)
}
