package admin

import (
	// "fmt"
	"newbee/models/manage/request"
	"newbee/services/adminservice"
	"newbee/web"
	"github.com/kataras/iris/v12"
)

type OrderController struct {
	Ctx iris.Context
}

func (m *OrderController) Get() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.OrdersReq)
	m.Ctx.ReadQuery(req)
	orderlist, total, err := adminservice.MallOrderService.GetOrderList(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"list"       : orderlist,
		"totalCount" : total,
		"currPage"   : req.PageNumber,
	})
}

func (m *OrderController) PutCheckdone() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.IdsReq)
	m.Ctx.ReadJSON(req)
	err := adminservice.MallOrderService.OrderCheckdone(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (m *OrderController) PutCheckout() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.IdsReq)
	m.Ctx.ReadJSON(req)
	err := adminservice.MallOrderService.OrderCheckout(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (m *OrderController) PutClose() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.IdsReq)
	m.Ctx.ReadJSON(req)
	err := adminservice.MallOrderService.OrderClose(req, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (m *OrderController) GetBy(id int) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	orderDetail, err := adminservice.MallOrderService.OrderDetail(id, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(orderDetail)
}