package controllers

import "github.com/astaxie/beego"

// ErrorController is error handler
type ErrorController struct {
	beego.Controller
}

// Error404 handles 404 error
func (ec *ErrorController) Error404() {
	ec.Data["content"] = "Hey，您访问的页面不存在，开发人员可能把它当点心吃掉了。This is not the page you're looking for. Developer might eat the page yet."
	ec.TplName = "errors/404.tpl"
}
