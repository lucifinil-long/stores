package routers

import (
	"github.com/lucifinil-long/stores/controllers"
	"github.com/lucifinil-long/stores/utils"

	"github.com/astaxie/beego"
)

func init() {
	beego.AddFuncMap("stringsToJson", utils.Strings2JSON)

	// pages
	beego.Router("/", &controllers.PageController{}, "*:Homepage")
	beego.Router("/pages/index", &controllers.PageController{}, "*:Homepage")
	beego.Router("/pages/login", &controllers.PageController{}, "*:LoginPage")

	// hompage related APIs
	beego.Router("/public/treemenu", &controllers.MainController{}, "*:TreeMenu")
	beego.Router("/public/isloggedin", &controllers.MainController{}, "*:IsLoggedIn")
	beego.Router("/public/login", &controllers.MainController{}, "*:Login")
	beego.Router("/public/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/public/changepwd", &controllers.MainController{}, "*:ChangePassword")
}
