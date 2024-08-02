package mallservice

import (
	"fmt"
	"newbee/global"
	"newbee/models/mall"
	"newbee/models/mall/response"
)

var ContactService = newContactService()

func newContactService() *contactService {
	return &contactService{}
}

type contactService struct {
}

func findLastestMessage(userIdA, userIdB int) *mall.MallMessage {
	message := new(mall.MallMessage)
	err := global.DB.Table("tb_newbee_mall_message").Where("(send_id = ? AND recv_id = ?) OR (send_id = ? AND recv_id = ?)",
		userIdA, userIdB, userIdB, userIdA).Order("create_time DESC").Take(message).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return message
}

func (m *contactService) GetUserList(token string) ([]*response.ContactResponse, error) {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil, err
	}
	var users []mall.MallUser
	global.DB.Model(&mall.MallUser{}).Find(&users)
	var contactResponse []*response.ContactResponse
	for _, u := range users {
		if u.UserId == userToken.UserId {
			continue
		}
		lastestMessage := findLastestMessage(userToken.UserId, u.UserId)
		if lastestMessage == nil {
			Res := &response.ContactResponse{
				NickName: u.NickName,
				UserId:   u.UserId,
			}
			contactResponse = append(contactResponse, Res)
		} else {
			Res := &response.ContactResponse{
				NickName:       u.NickName,
				UserId:         u.UserId,
				MessageContent: lastestMessage.Content,
				MessageTime:    lastestMessage.CreateTime.Format("2006-01-02 15:04:05"),
			}
			contactResponse = append(contactResponse, Res)
		}
	}
	return contactResponse, nil
}
