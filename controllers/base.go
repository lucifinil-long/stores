package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/lucifinil-long/stores/config"
	"github.com/lucifinil-long/stores/models"
	"github.com/lucifinil-long/stores/proto"
	"github.com/lucifinil-long/stores/utils"
	"github.com/mkideal/log"
)

func init() {
	if log.GetLevel() == log.LvTRACE {
		logReq := func(ctx *context.Context) {
			dataBytes, _ := json.Marshal(ctx.Request.Form)
			log.Trace("Get '%v' request from '%v' to '%v' with data %v",
				ctx.Request.Method, ctx.Request.RemoteAddr, ctx.Request.RequestURI, string(dataBytes))
		}

		beego.InsertFilter("/*", beego.BeforeExec, logReq)
	}

	// 验证权限
	AccessRegister()
}

// BaseController is the base of all controllers
type BaseController struct {
	beego.Controller
}

// Response is http response handler
func (bc *BaseController) Response(rsp *proto.Response) {
	bc.Data[cJSONKey] = rsp
	bc.ServeJSON()

	if log.GetLevel() == log.LvTRACE {
		respBytes, _ := json.Marshal(bc.Data[cJSONKey])
		log.Trace("responsed %v for '%v' request from '%v' for path '%v'", string(respBytes),
			bc.Ctx.Request.Method, bc.Ctx.Request.RemoteAddr, bc.Ctx.Request.RequestURI)
	}
}

// CurrentUser get current login user from this session
func (bc BaseController) CurrentUser() *proto.User {
	uinfo := bc.Ctx.Input.Session(cSessionUserInfoKey)
	if uinfo == nil {
		return nil
	}

	user := uinfo.(proto.User)
	return &user
}

// OperationLog adds a operation log entry
// @param session is database session, can be nil; if nil will use default database session
// @param action is action summary
// @param details are more addtional detail
func (bc BaseController) OperationLog(action string, details ...string) {
	uid := 0
	username := ""
	user := bc.GetSession(cSessionUserInfoKey)
	if user != nil {
		uid = user.(proto.User).ID
		username = user.(proto.User).Username
	}

	detail := strings.Join(details, ", ")
	log.Info("%v(UID:%v) performed '%v' action, detail: %v", username, uid, action, detail)

	models.AddOperationLog(uid, username,
		bc.Ctx.Request.RemoteAddr, action, detail)
}

// AccessRegister checks access and register user's nodes
func AccessRegister() {
	var Check = func(ctx *context.Context) {

		configs := config.GetConfigs()
		authType := configs.UserAuthType()
		authGateway := configs.AuthGateway()
		var accesslist map[string]bool

		if authType != config.UserAuthTypeNone {
			params := strings.Split(strings.ToLower(strings.Split(ctx.Request.RequestURI, "?")[0]), "/")
			if CheckPackageNeedAccess(params[1]) {
				uinfo := ctx.Input.Session(cSessionUserInfoKey)
				if uinfo == nil {
					rsp := &proto.Response{
						Status:      proto.ReturnStatusTempRedirect,
						Description: cTipRelogin,
						Protocol:    authGateway,
					}
					ctx.Output.JSON(rsp, true, false)
					log.Info("Don't allowed '%v' request from %v to '%v' for no logged in user information.",
						ctx.Request.Method, ctx.Request.RemoteAddr, ctx.Request.RequestURI)
					return
				}
				// super admin用户不用认证权限
				currentUser := uinfo.(proto.User)
				if models.IsSuperAdmin(&currentUser) {
					log.Trace("super administrator no need to check access list")
					return
				}

				if authType == config.UserAuthTypeLogin {
					listFromSession := ctx.Input.Session(cSessionAccessListKey)
					if listFromSession != nil {
						accesslist = listFromSession.(map[string]bool)
					}
				} else if authType == config.UserAuthTypeRealtime {
					session := configs.OrmEngine.NewSession()
					defer session.Close()
					accesslist, _ = models.GetUserAccessList(currentUser.ID)
				}

				ret := AccessDecision(params, accesslist)
				if !ret {
					log.Info("Don't allowed '%v' request from %v to '%v' for no authority.",
						ctx.Request.Method, ctx.Request.RemoteAddr, ctx.Request.RequestURI)
					ctx.Output.JSON(
						&proto.Response{
							Status:      proto.ReturnStatusNotAuthorize,
							Description: cTipNoAuth,
							Protocol:    "",
						},
						true,
						false)
				}
			}
		}
	}

	beego.InsertFilter("/*", beego.BeforeRouter, Check)
}

// CheckPackageNeedAccess determines whether need to verify
func CheckPackageNeedAccess(pack string) bool {
	if len(pack) == 0 {
		return false
	}

	packages := config.GetConfigs().NotAuthPackages()
	for _, nap := range packages {
		if pack == nap {
			return false
		}
	}
	return true
}

// AccessDecision test whether permissions for node
func AccessDecision(params []string, accesslist map[string]bool) bool {
	if CheckPackageNeedAccess(params[0]) {
		if len(accesslist) < 1 {
			return false
		}

		s := strings.Join(params, "/")
		_, ok := accesslist[s]
		if ok {
			return true
		}
	} else {
		return true
	}
	return false
}

// CheckLogin validates login info
// @param session is database xorm session
// @param username is the username that be checked
// @param password is the correspond password that be checked
// @return (*proto.User, nil) if validate user infomation successfuly, otherwise return (nil, error)
func CheckLogin(username, password string) (*proto.User, error) {
	user, err := models.GetUserInfoByUsername(username)
	if err != nil {
		return nil, err
	}

	passwordMd5 := strings.ToLower(utils.String2MD5(password))
	password = strings.ToLower(password)
	userDBPwd := strings.ToLower(user.Password)

	if userDBPwd != password && userDBPwd != passwordMd5 {
		return nil, models.ErrUserWrongPwd
	}

	return user, nil
}
