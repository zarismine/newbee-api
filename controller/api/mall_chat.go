package api

import (
	// "fmt"
	"errors"
	"newbee/services/mallservice"
	"newbee/web"

	"github.com/kataras/iris/v12"
)

type ChatController struct {
	Ctx iris.Context
}

func (m *ChatController) Get() {
	mallservice.ChatService.Chat(m.Ctx.ResponseWriter(), m.Ctx.Request())
}

func (m *ChatController) GetBy(id int) *web.JsonResult {
	token := m.Ctx.GetHeader("Token")
	messages, yourname, err := mallservice.ChatService.GetRecord(id, token)
	user := mallservice.MallUserService.GetById(id)
	if user == nil {
		return web.JsonError(errors.New("无效的id"))
	}
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(map[string]interface{} {
		"yourname" : yourname,
		"nickname" : user.NickName,
		"messages" : messages,
	})
}