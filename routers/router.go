package routers

import (
	"github.com/lucifinil-long/stores/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
