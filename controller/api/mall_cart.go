package api

import (
	"errors"
	"newbee/models/mall/request"
	"newbee/services/mallservice"
	"newbee/web"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
)

type MallCartController struct {
	Ctx iris.Context
}

func (m *MallCartController) Get() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	list, err := mallservice.MallCartService.GetShoppingCartItemByToken(token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(list)
}

func (m *MallCartController) Put() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.EditCart)
	m.Ctx.ReadJSON(req)
	err := mallservice.MallCartService.EditShoppingCartItem(req.GoodsCount, req.CartItemId, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (m *MallCartController) Post() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.AddCart)
	m.Ctx.ReadJSON(req)
	err := mallservice.MallCartService.AddShoppingCartItem(req.GoodsId, req.GoodsCount, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (m *MallCartController) DeleteBy(id int) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	req := new(request.EditCart)
	m.Ctx.ReadJSON(req)
	err := mallservice.MallCartService.DeleteCartItemById(id, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (m *MallCartController) GetSettle() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	cartItemIds := m.Ctx.URLParam("cartItemIds")
	cartItemId_Str := strings.Split(cartItemIds, ",")
    cartItemIds_Int := make([]int, len(cartItemId_Str))
    for i, s := range cartItemId_Str {
        num, err := strconv.Atoi(s)
        if err != nil {
            return web.JsonError(errors.New("参数错误"))
        }
        cartItemIds_Int[i] = num
    }
	list, err := mallservice.MallCartService.GetShoppingCartItemById(cartItemIds_Int, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(list)
}