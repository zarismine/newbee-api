package admin

import (
	"newbee/models/manage/request"
	"newbee/services/adminservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type UserController struct {
	Ctx iris.Context
}

func (u *UserController) Get() *web.JsonResult {
	token := u.Ctx.GetHeader("Token")
	req := new(request.PageInfo)
	u.Ctx.ReadQuery(req)
	list, total, err := adminservice.MallUserService.GetMallUserInfoList(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"list" : list,
		"totalCount" : total,
		"currPage"   : req.PageNumber,
	})
}

func (u *UserController) PutBy (id int) *web.JsonResult {
	token := u.Ctx.GetHeader("Token")
	req := new(request.IdsReq)
	u.Ctx.ReadJSON(req)
	err := adminservice.MallUserService.LockUser(req.Ids, id, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}