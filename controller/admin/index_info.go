package admin

import (
	"newbee/models/manage/request"
	"newbee/models/manage/response"
	"newbee/services/adminservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type IndexInfoController struct {
	Ctx iris.Context
}

func (m *IndexInfoController) Post() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	var req request.MallIndexConfigAddParams
	_ = m.Ctx.ReadJSON(&req)
	if err := adminservice.IndexInfoService.CreateMallIndexConfig(req, token); err != nil {
		return web.JsonError(err)
	} 
	return web.JsonSuccess()
}

func (m *IndexInfoController) Delete() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	var req request.IdsReq
	_ = m.Ctx.ReadJSON(&req)
	if err := adminservice.IndexInfoService.DeleteMallIndexConfig(req, token); err != nil {
		return web.JsonError(err)
	} 
	return web.JsonSuccess()
}

func (m *IndexInfoController) Put() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	var req request.MallIndexConfigUpdateParams
	_ = m.Ctx.ReadJSON(&req)
	if err := adminservice.IndexInfoService.UpdateMallIndexConfig(req, token); err != nil {
		return web.JsonError(err)
	} 
	return web.JsonSuccess()
}

func (m IndexInfoController) GetBy(id int) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	mallIndexConfig, err := adminservice.IndexInfoService.GetMallIndexConfig(uint(id),token)
	if err != nil {
		return web.JsonError(err)
	} 
	return web.JsonData(mallIndexConfig)
}

func (m IndexInfoController) Get() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	var pageInfo request.MallIndexConfigSearch
	_ = m.Ctx.ReadQuery(&pageInfo)
	if list, total, err := adminservice.IndexInfoService.GetMallIndexConfigInfoList(pageInfo,token); err != nil {
		return web.JsonError(err)
	} else {
		return web.JsonData(response.PageResult{
			List:       list,
			TotalCount: total,
			CurrPage:   pageInfo.PageNumber,
			PageSize:   pageInfo.PageSize,
		})
	}
}