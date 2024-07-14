package mallservice

import (
	"errors"
	"newbee/global"
	"newbee/models/mall"
	"newbee/pkg/dates"
	"newbee/pkg/passwd"
	"time"
)

var MallUserTokenService = newMallUserTokenService()

func newMallUserTokenService() *mallUserTokenService {
	return &mallUserTokenService{}
}

type mallUserTokenService struct {
}

func (m *mallUserTokenService) Create(token *mall.MallUserToken) (err error) {
	err = global.DB.Create(token).Error
	return
}

func (m *mallUserTokenService) Take(where ...interface{}) *mall.MallUserToken {
	ret := &mall.MallUserToken{}
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (m *mallUserTokenService) Save (token *mall.MallUserToken) (err error) {
	return global.DB.Save(token).Error
}
func (m *mallUserTokenService) GetById (Id int) *mall.MallUserToken {
	return m.Take("user_id = ?",Id)
}

func (m *mallUserTokenService) GetUserTokenByToken (token string) (*mall.MallUserToken,error) {
	usertoken := m.Take("token = ?",token)
	if usertoken == nil {
		return nil,errors.New("无效token")
	}
	if m.IsVaildToken(usertoken){
		m.UpdateToken(usertoken)
		return usertoken,nil
	}else{
		return nil,errors.New("无效token")
	}
}

func (m *mallUserTokenService) DeleteUserToken (token string) error {
	ret := m.Take("token = ?",token)
	if ret == nil {
		return errors.New("invaild token")
	}
	return global.DB.Model(ret).Update("expire_time",dates.NowTimestamp()).Error
}

func (m *mallUserTokenService) UpdateToken(token *mall.MallUserToken) error {
	return global.DB.Model(token).Update("expire_time",dates.NowTimestamp() + 172800000).Error
}

func (m *mallUserTokenService) GenerateToken(id int) *mall.MallUserToken {
	token := passwd.UUID()
	nowDate := time.Now()
	expireTime, _ := time.ParseDuration("48h")
	expireDate := nowDate.Add(expireTime)
	MallUserToken := &mall.MallUserToken{
		Token     :      token,
		UserId    :      id,
		UpdateTime:      dates.Timestamp(nowDate),
		ExpireTime:      dates.Timestamp(expireDate),
	}
	err := m.Create(MallUserToken)
	if err != nil {
		return nil
	}
	return MallUserToken
}

func (m *mallUserTokenService) IsVaildToken(token *mall.MallUserToken) bool {
	user := MallUserService.GetById(token.UserId)
	if user == nil {
		return false
	}
	if token.ExpireTime > dates.NowTimestamp() && token != nil && user.LockedFlag == 0{
		return true
	}
	return false
}