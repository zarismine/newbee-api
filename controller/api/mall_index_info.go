package api

import (
	"newbee/services/mallservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type MallIndexInfoController struct {
	Ctx iris.Context
}

func (m MallIndexInfoController) Get() *web.JsonResult {
	hotGoodses, newGoodses, recommendGoodses, err := mallservice.MallIndexInfoService.GetDetailData()
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"carousels" : nil,
		"hotGoodses" : hotGoodses,
		"newGoodses" : newGoodses,
		"recommendGoodses" : recommendGoodses,
	})
}