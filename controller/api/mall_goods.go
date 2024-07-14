package api

import (
	"newbee/models/mall/request"
	"newbee/services/mallservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type MallGoodsController struct {
	Ctx iris.Context
}

func (m *MallGoodsController) GetSearch() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.PageSearch)
	m.Ctx.ReadQuery(req)
	list, total, err := mallservice.MallGoodsService.SearchByCategory(req.PageNumber, req.PageSize, req.GoodsCategoryId, req.Keyword, req.OrderBy, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"currPage"   : req.PageNumber,
		"pageSize"   : req.PageSize,
		"list"       : list,
		"totalCount" : total,
	})
}

func (m *MallGoodsController) GetDetailBy(id int) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	goodsDetail, err := mallservice.MallGoodsService.DetailByGoodsId(id, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(goodsDetail)
}