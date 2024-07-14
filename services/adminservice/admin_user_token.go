package adminservice

import (
	"errors"
	"newbee/global"
	"newbee/models/manage"
	"newbee/pkg/dates"
	"newbee/pkg/passwd"
	"time"
)

var AdminUserTokenService = newAdminUserTokenService()

func newAdminUserTokenService() *adminUserTokenService {
	return &adminUserTokenService{}
}

type adminUserTokenService struct {
}

func (a *adminUserTokenService) Create (adminToken *manage.MallAdminUserToken) error {
	return global.DB.Create(adminToken).Error
}

func (a *adminUserTokenService) Take (where ...interface{}) *manage.MallAdminUserToken {
	ret := new(manage.MallAdminUserToken)
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (m *adminUserTokenService) GetById (Id int) *manage.MallAdminUserToken {
	return m.Take("admin_user_id = ?",Id)
}

func (m *adminUserTokenService) IsVaildToken(token *manage.MallAdminUserToken) bool {
	if token.ExpireTime > dates.NowTimestamp() && token != nil{
		return true
	}
	return false
}

func (m *adminUserTokenService) GetByToken (token string) (*manage.MallAdminUserToken, error) {
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

func (a *adminUserTokenService) UpdateToken (token *manage.MallAdminUserToken) error {
	return global.DB.Model(token).Update("expire_time",dates.NowTimestamp() + 172800000).Error
}

func (a *adminUserTokenService) GenerateToken(id int) *manage.MallAdminUserToken {
	token := passwd.UUID()
	nowDate := time.Now()
	expireTime, _ := time.ParseDuration("96h")
	expireDate := nowDate.Add(expireTime)
	MallUserToken := &manage.MallAdminUserToken{
		Token       :      token,
		AdminUserId :      id,
		UpdateTime  :      dates.Timestamp(nowDate),
		ExpireTime  :      dates.Timestamp(expireDate),
	}
	err := a.Create(MallUserToken)
	if err != nil {
		return nil
	}
	return MallUserToken
}