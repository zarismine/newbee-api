package mallservice

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"newbee/global"
	"newbee/models/mall"
	"newbee/models/mall/response"
	"strconv"
	"time"
)

var ContactService = newContactService()

func newContactService() *contactService {
	return &contactService{}
}

type contactService struct {
}

func findLastestMessage(userIdA, userIdB int) (*mall.MallMessage, int64) {
	message := new(mall.MallMessage)
	db := global.DB.Table("tb_newbee_mall_message").Where("(send_id = ? AND recv_id = ?) OR (send_id = ? AND recv_id = ?)", userIdA, userIdB, userIdB, userIdA)
	var count int64
	err := db.Order("create_time DESC").Take(message).Error
	global.DB.Table("tb_newbee_mall_message").Where("send_id = ? AND recv_id = ? AND message_status = 1", userIdB, userIdA).Count(&count)
	if err != nil {
		fmt.Println(err)
		return nil, 0
	}
	return message, count
}

func (m *contactService) GetUserList(token string) ([]*response.ContactResponse, error) {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil, err
	}
	var contactResponse []*response.ContactResponse
	ctx := context.Background()
	cacheUserContact := fmt.Sprintf("%s%v", global.CacheUserContactPrefix, userToken.UserId)
	res, _ := global.Redis.HVals(ctx, cacheUserContact).Result()
	if len(res) != 0 {
		for _, v := range res {
			resp := new(response.ContactResponse)
			err = jsoniter.Unmarshal([]byte(v), resp)
			if err != nil {
				return nil, err
			}
			contactResponse = append(contactResponse, resp)
		}
		return contactResponse, nil
	}

	var users []mall.MallUser
	global.DB.Model(&mall.MallUser{}).Find(&users)

	for _, u := range users {
		if u.UserId == userToken.UserId {
			continue
		}
		lastestMessage, noReadMsgCount := findLastestMessage(userToken.UserId, u.UserId)
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
				Count:          noReadMsgCount,
			}
			contactResponse = append(contactResponse, Res)
		}
	}
	data := make(map[string]interface{})
	for _, v := range contactResponse {
		temp, _ := jsoniter.Marshal(*v)
		data[strconv.Itoa(v.UserId)] = temp
	}
	global.Redis.HMSet(ctx, cacheUserContact, data)
	global.Redis.Expire(ctx, cacheUserContact, 2*time.Hour)
	return contactResponse, nil
}
