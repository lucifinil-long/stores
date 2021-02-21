package main

import (
	"github.com/astaxie/beego"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/controllers"
	_ "github.com/lucifinil-long/stores/routers"
	"github.com/mkideal/log"
)

func main() {
	// 强制使用session，后继需要保存登录信息
	beego.BConfig.WebConfig.Session.SessionOn = true

	// 强制使用https
	beego.BConfig.Listen.EnableHTTP = false
	beego.BConfig.Listen.EnableHTTPS = true

	// initializes configs and log environment
	config.InitConfigs()
	// uninitialize og environment
	defer log.Uninit(nil)

	beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}
