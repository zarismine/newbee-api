package api

import (
	"errors"
	"newbee/models/mall"
	"newbee/models/mall/request"
	"newbee/services/mallservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type MallUserController struct {
	Ctx iris.Context
}

func (c *MallUserController) GetBy(id int) *web.JsonResult {
	user := mallservice.MallUserService.GetById(id)
	if user != nil {
		return web.JsonData(c.buildbyUser(user))
	}
	return web.JsonError(errors.New("用户不存在"))
}

func (c *MallUserController) PostRegister() *web.JsonResult {
	var req request.RegisterUserParam
	_ = c.Ctx.ReadJSON(&req)
	user, err := mallservice.MallUserService.SignUp(req.LoginName, "", req.Password)
	if err != nil {
		return web.JsonError(err)
	}
	return &web.JsonResult{
		ErrorCode: 0,
		Message  : "注册成功",
		Data     : c.buildbyUser(user),
		Success  : true,
	}
}

func (c *MallUserController) PostLogin() *web.JsonResult {
	var req request.UserLoginParam
	_ = c.Ctx.ReadJSON(&req)
	user, token, err := mallservice.MallUserService.SignIn(req.LoginName, req.PasswordMd5)
	if err != nil {
		return web.JsonError(err)
	}
	res := c.buildbyUser(user)
	res["token"] = token
	return web.JsonData(res)
}

func (c *MallUserController) GetInfo() *web.JsonResult {
	token := c.Ctx.GetHeader("Token")
	user,err := mallservice.MallUserService.GetUserByToken(token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(c.buildbyUser(user))
}

func (c *MallUserController) PutInfo() *web.JsonResult {
	var req request.UpdateUserInfoParam
	token := c.Ctx.GetHeader("Token")
	user,err := mallservice.MallUserService.GetUserByToken(token)
	if err != nil {
		return web.JsonError(err)
	}
	_ = c.Ctx.ReadJSON(&req)
	err = mallservice.MallUserService.UpdateUser(req.NickName, req.PasswordMd5, req.IntroduceSign, user)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *MallUserController) PostLogout() *web.JsonResult {
	token := c.Ctx.GetHeader("Token")
	err := mallservice.MallUserTokenService.DeleteUserToken(token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *MallUserController) buildbyUser(u *mall.MallUser) map[string]interface{} {
	data := make(map[string]interface{})
	data["loginname"] = u.LoginName
	data["nickname"] = u.NickName
	data["introducesign"] = u.IntroduceSign
	return data
}

