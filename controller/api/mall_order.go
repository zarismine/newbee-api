package api

import (
	"newbee/models/mall/request"
	"newbee/services/mallservice"
	"newbee/web"
	"strconv"

	"github.com/kataras/iris/v12"
)

type MallOrderController struct {
	Ctx iris.Context
}

func (m *MallOrderController) Post () *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.SaveOrder)
	m.Ctx.ReadJSON(req)
	orderNo, err := mallservice.MallOrderService.Save(req.CartItemIds, req.AddressId, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(orderNo)
}

func (m *MallOrderController) Get() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	status := m.Ctx.URLParam("status")
	list, toal, err := mallservice.MallOrderService.GetOrderList(token, status)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"currPage"   : 1,
		"list"       : list,
		"pageSize"   : 5,
		"totalCount" : toal,
		"totalPage"  : 0,
	})
}

func (m *MallOrderController) GetBy(orderNo string) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	order, err := mallservice.MallOrderService.GetOrderItemByorderNo(orderNo, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(order)
}

func (m *MallOrderController) GetPay() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	orderNo := m.Ctx.URLParam("orderNo")
	payType, _ := strconv.Atoi(m.Ctx.URLParam("payType"))
	err := mallservice.MallOrderService.PaySuccess(orderNo, token, payType)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (m *MallOrderController) PutFinishBy(orderNo string) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	err := mallservice.MallOrderService.FinishOrder(orderNo, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (m *MallOrderController) PutCancelBy(orderNo string) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	err := mallservice.MallOrderService.CancelOrder(orderNo, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
