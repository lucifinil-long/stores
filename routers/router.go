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
	beego.Router("/pages/admin/commodities", &controllers.PageController{}, "*:CommoditiesPage")
	beego.Router("/pages/admin/locations", &controllers.PageController{}, "*:AdminLocationsPage")
	beego.Router("/pages/admin/specifications", &controllers.PageController{}, "*:AdminSpecificationsPage")

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

	// locations related
	beego.Router("/admin/depot/list", &controllers.LocationController{}, "*:DepotList")
	beego.Router("/admin/depot/delete", &controllers.LocationController{}, "*:DeleteDepot")
	beego.Router("/admin/depot/add", &controllers.LocationController{}, "*:AddDepot")
	beego.Router("/admin/depot/update", &controllers.LocationController{}, "*:UpdateDepot")
	beego.Router("/admin/depot/shelf/list", &controllers.LocationController{}, "*:ShelfList")
	beego.Router("/admin/depot/shelf/add", &controllers.LocationController{}, "*:AddShelfs")
	beego.Router("/admin/depot/shelf/delete", &controllers.LocationController{}, "*:DeleteShelf")
	beego.Router("/admin/depot/shelf/update", &controllers.LocationController{}, "*:UpdateShelf")

	// specifications related
	beego.Router("/admin/specifications/list", &controllers.SpecificationController{}, "*:SpecList")
	beego.Router("/admin/specifications/add", &controllers.SpecificationController{}, "*:AddSpec")
	beego.Router("/admin/specifications/update", &controllers.SpecificationController{}, "*:UpdateSpec")
}
