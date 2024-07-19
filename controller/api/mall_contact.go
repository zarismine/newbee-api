package api

import (
	// "fmt"
	"newbee/services/mallservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type ContactController struct {
	Ctx iris.Context
}

func (m *ContactController) Get() *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	list, err := mallservice.ContactService.GetUserList(token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(list)
}