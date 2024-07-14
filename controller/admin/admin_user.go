package admin

import (
	"newbee/models/manage"
	"newbee/models/manage/request"
	"newbee/services/adminservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type AdminUserController struct {
	Ctx iris.Context
}

func (a *AdminUserController) PostLogin() *web.JsonResult {
	req := new(request.MallAdminLoginParam)
	a.Ctx.ReadJSON(req)
	token, err := adminservice.AdminUser.AdminSignIn(req.UserName,req.PasswordMd5)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]string{
		"token" : token,
	})
}

func (a *AdminUserController) PostRegister() *web.JsonResult {
	req := new(request.MallAdminRegisterParam)
	a.Ctx.ReadJSON(req)
	adminUser, err := adminservice.AdminUser.AdminSignUp(req.LoginUserName, req.LoginPassword, req.NickName, req.KEY)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(a.buildByUser(adminUser))
}

func (a *AdminUserController) GetProfile() *web.JsonResult {
	token := a.Ctx.GetHeader("Token")
	adminUser, err := adminservice.AdminUser.GetProfileByToken(token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(a.buildByUser(adminUser))
}

func (a *AdminUserController) buildByUser (user *manage.MallAdminUser) map[string]interface{} {
	return map[string]interface{}{
		"loginUserName" : user.LoginUserName,
		"nickName"  : user.NickName,
		"adminUserId"  : user.AdminUserId,
	}
}