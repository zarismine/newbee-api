package admin

import (
	"newbee/models/manage/request"
	"newbee/services/adminservice"
	"newbee/web"
	"strconv"

	"github.com/kataras/iris/v12"
)

type GoodCategotyController struct {
	Ctx iris.Context
}

func (gc *GoodCategotyController) Post() *web.JsonResult {
	token := gc.Ctx.GetHeader("Token")
	req := new(request.MallGoodsCategoryReq)
	gc.Ctx.ReadJSON(req)
	rank, _ := strconv.Atoi(req.CategoryRank)
	err := adminservice.GoodsCategoryService.AddGoodsCategory(req.ParentId,req.CategoryLevel, rank, req.CategoryName, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (gc *GoodCategotyController) Put() *web.JsonResult {
	token := gc.Ctx.GetHeader("Token")
	req := new(request.MallGoodsCategoryReq)
	gc.Ctx.ReadJSON(req)
	rank, _ := strconv.Atoi(req.CategoryRank)
	err := adminservice.GoodsCategoryService.UpdateGoodsCategory(req.CategoryId, req.CategoryLevel, rank, req.CategoryName, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (gc *GoodCategotyController) Get() *web.JsonResult {
	token := gc.Ctx.GetHeader("Token")
	req := new(request.SearchCategoryParams)
	gc.Ctx.ReadQuery(req)
	total, goodsCategory, err := adminservice.GoodsCategoryService.SearchCategory(req,token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"list"        :   goodsCategory,
		"totalCount"  :   total,
		"currPage"    :   req.PageNumber,
	})
}

func (gc *GoodCategotyController) GetBy(id int) *web.JsonResult {
	token := gc.Ctx.GetHeader("Token")
	goodsCategory, err := adminservice.GoodsCategoryService.SearchCategoryById(id,token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"categoryName"   :   goodsCategory.CategoryName,
		"categoryRank"   :   goodsCategory.CategoryRank,
		"parentId"       :   goodsCategory.ParentId,
		"categoryLevel"  :   goodsCategory.CategoryLevel,
	})
}

func (gc *GoodCategotyController) Delete() *web.JsonResult {
	token := gc.Ctx.GetHeader("Token")
	req := new(request.IdsReq)
	gc.Ctx.ReadJSON(req)
	err := adminservice.GoodsCategoryService.DeleteById(req.Ids,token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}