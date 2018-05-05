package routers

import (
	"github.com/lucifinil-long/stores/controllers"
	"github.com/lucifinil-long/stores/utils"

	"github.com/astaxie/beego"
)

func init() {
	beego.AddFuncMap("stringsToJson", utils.Strings2JSON)

	// pages
	beego.Router("/", &controllers.PageController{}, "*:AdminHomepage")
	beego.Router("/pages/admin/index", &controllers.PageController{}, "*:AdminHomepage")
	beego.Router("/pages/admin/login", &controllers.PageController{}, "*:AdminLoginPage")
	beego.Router("/pages/admin/users", &controllers.PageController{}, "*:AdminUsersPage")
	beego.Router("/pages/admin/operations", &controllers.PageController{}, "*:AdminOperationsPage")

	beego.Router("/public/isloggedin", &controllers.MainController{}, "*:IsLoggedIn")
	// hompage related APIs
	beego.Router("/public/admin/treemenu", &controllers.MainController{}, "*:TreeMenu")
	beego.Router("/public/admin/login", &controllers.MainController{}, "*:Login")
	beego.Router("/public/admin/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/public/admin/changepwd", &controllers.MainController{}, "*:ChangePassword")
	beego.Router("/public/admin/accesslist", &controllers.MainController{}, "*:AccessList")

	// system user related
	beego.Router("/admin/user/list", &controllers.UserController{}, "*:UserList")
	beego.Router("/admin/user/add", &controllers.UserController{}, "*:AddUser")
	beego.Router("/admin/user/delete", &controllers.UserController{}, "*:DeleteUser")
}
