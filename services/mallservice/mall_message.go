package mallservice

import (
	"errors"
	"fmt"
	"newbee/global"
	"newbee/models/mall"
	"strconv"
)

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct {
}

func (m *messageService) UpdateMessage(messageId int, token string) error {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return err
	}
	msg := new(mall.MallMessage)
	global.DB.Model(&mall.MallMessage{}).Where("message_id = ?", messageId).First(msg)
	if msg.RecvId != userToken.UserId {
		return errors.New("无权限更改")
	}
	cacheUserContact := fmt.Sprintf("%s%v", global.CacheUserContactPrefix, msg.RecvId)
	UpDateCache(cacheUserContact, strconv.Itoa(msg.SendId), "", true)
	err = global.DB.Model(msg).Update("message_status", 0).Error
	return err
}
