package api

import (
	"errors"
	"newbee/models/mall/request"
	"newbee/services/mallservice"
	"newbee/web"
	"github.com/kataras/iris/v12"
)

type MallUserAddressController struct {
	Ctx iris.Context
}

func (m *MallUserAddressController) Get() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	if userAddressList,err := mallservice.MallUserAddressService.GetAddressByToken(token);err != nil {
		return web.JsonError(err)
	}else if len(userAddressList) == 0 {
		return web.JsonSuccess()
	}else {
		return web.JsonData(userAddressList)
	}
}

func (m *MallUserAddressController) Post() *web.JsonResult {
	req := new(request.AddAddressParam)
	token := m.Ctx.GetHeader("Token")
	_ = m.Ctx.ReadJSON(req)
	err := mallservice.MallUserAddressService.AddUserAddress(token,req)
	if err != nil {
		return web.JsonError(err)
	}else {
		return web.JsonSuccess()
	}
}

func (m *MallUserAddressController) Put() *web.JsonResult {
	req := new(request.UpdateAddressParam)
	token := m.Ctx.GetHeader("Token")
	_ = m.Ctx.ReadJSON(req)
	err := mallservice.MallUserAddressService.EditUserAddress(token,req)
	if err != nil {
		return web.JsonError(err)
	}else {
		return web.JsonSuccess()
	}
}

func (m *MallUserAddressController) GetBy (addressid int) *web.JsonResult {
	userAddress := mallservice.MallUserAddressService.GetAddressByAddressId(addressid)
	if userAddress != nil {
		return web.JsonData(userAddress)
	}else {
		return web.JsonError(errors.New("参数错误"))
	}
}

func (m *MallUserAddressController) GetDefault () *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	defaultAddress, err := mallservice.MallUserAddressService.GetDefaultAddressByToken(token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(defaultAddress)
}

func (m *MallUserAddressController) DeleteBy (addressid int) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	err := mallservice.MallUserAddressService.DeleteByAddressId(token,addressid)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}