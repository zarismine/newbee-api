package admin

import (
	"errors"
	"newbee/models/manage/request"
	"newbee/services/adminservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type GoodsController struct {
	Ctx iris.Context
}

func (g *GoodsController) Post() *web.JsonResult {
	token := g.Ctx.GetHeader("Token")
	req := new(request.GoodsInfoAddParam)
	g.Ctx.ReadJSON(req)
	err := adminservice.GoodsService.AddGoods(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (g *GoodsController) Put() *web.JsonResult {
	token := g.Ctx.GetHeader("Token")
	req := new(request.GoodsInfoUpdateParam)
	g.Ctx.ReadJSON(req)
	err := adminservice.GoodsService.UpdateGoods(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (g *GoodsController) GetList() *web.JsonResult {
	token := g.Ctx.GetHeader("Token")
	req := new(request.PageInfo)
	g.Ctx.ReadQuery(req)
	goodslist, total, err := adminservice.GoodsService.GetGoodsList(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"list"       : goodslist,
		"totalCount" : total,
		"currPage"   : req.PageNumber,
	})
}


func (g *GoodsController) GetBy(id int) *web.JsonResult {
	token := g.Ctx.GetHeader("Token")
	req := new(request.PageInfo)
	g.Ctx.ReadQuery(req)
	goods,  err := adminservice.GoodsService.GetGoodsById(id, token)
	if err != nil {
		return web.JsonError(err)
	}
	categoryLink, err := adminservice.GoodsCategoryService.SearchCategoryLinkById(goods.GoodsCategoryId,token)
	if err != nil {
		return web.JsonError(err)
	}
	categoryLink["goods"] = goods
	return web.JsonData(categoryLink)
}

func (g *GoodsController) PutStatusBy(id int) *web.JsonResult {
	token := g.Ctx.GetHeader("Token")
	req := new(request.IdsReq)
	g.Ctx.ReadJSON(req)
	if id == 1  || id ==0{
		err := adminservice.GoodsService.UpdateStatusByIds(id, req.Ids, token)
		if err != nil {
			return web.JsonError(err)
		}
		return web.JsonSuccess()
	}
	return web.JsonError(errors.New("参数错误"))
}

func (g *GoodsController) GetSearch() *web.JsonResult {
	token := g.Ctx.GetHeader("Token")
	req := new(request.PageInfo)
	g.Ctx.ReadQuery(req)
	goodslist, total, err := adminservice.GoodsService.SearchGoodsList(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"list"       : goodslist,
		"totalCount" : total,
		"currPage"   : req.PageNumber,
	})
}