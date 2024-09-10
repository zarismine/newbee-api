package api

import (
	"newbee/services/mallservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type MessageController struct {
	Ctx iris.Context
}

func (m *MessageController) GetBy(id int) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	err := mallservice.MessageService.UpdateMessage(id, token)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}