package api

import (
	"newbee/services/mallservice"
	"newbee/web"
	"github.com/kataras/iris/v12"
)

type MallCategoryController struct {
	Ctx iris.Context
}

func (m *MallCategoryController) Get() *web.JsonResult {
	list, err := mallservice.MallCategoryService.GetList()
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(list)
}