package adminservice

import (
	"errors"
	"newbee/global"
	"newbee/models/manage"
	"newbee/pkg/passwd"
)

var AdminUser = newAdminUser()

func newAdminUser() *adminUser {
	return &adminUser{}
}

type adminUser struct {
}

func (a *adminUser) GetByLoginName (name string) *manage.MallAdminUser {
	user := new(manage.MallAdminUser)
	if err := global.DB.Where("login_user_name = ?", name).First(user).Error;err != nil {
		return nil
	}
	return user
}

func (a *adminUser) GetById (id int) *manage.MallAdminUser {
	user := new(manage.MallAdminUser)
	if err := global.DB.Where("admin_user_id = ?", id).First(user).Error;err != nil {
		return nil
	}
	return user
}

func (a *adminUser) Create (adminUser *manage.MallAdminUser) error {
	return global.DB.Create(adminUser).Error
}

func (a *adminUser) AdminSignUp (LoginName,PasswordMd5,NickName,KEY string) (*manage.MallAdminUser,error) {
	if passwd.Hash([]byte(KEY)) != "1ff75a6a8513c212" {
		return nil,errors.New("关键KEY错误")
	}
	if len(LoginName) == 0 {
		return nil, errors.New("账号不能为空")
	}
	if len(PasswordMd5) == 0 {
		return nil, errors.New("密码不能为空")
	}
	if len(NickName) == 0 {
		NickName = "admin_user_" + LoginName
	}
	adminUser := a.GetByLoginName(LoginName)
	if adminUser != nil {
		return nil, errors.New("账号已经存在")
	}
	adminUser = &manage.MallAdminUser{
		LoginUserName : LoginName,
		LoginPassword : passwd.Hash([]byte(PasswordMd5)),
		NickName      : NickName,
	}
	a.Create(adminUser)
	return adminUser, nil
}

func (a *adminUser) AdminSignIn (LoginName,PasswordMd5 string) (string,error) {
	if len(LoginName) == 0 {
		return "", errors.New("账号不能为空")
	}
	adminUser := a.GetByLoginName(LoginName)
	if adminUser == nil {
		return "", errors.New("账号不存在")
	}
	if adminUser.LoginPassword != PasswordMd5 {
		return "", errors.New("账号或密码错误")
	}
	adminUserToken := AdminUserTokenService.GetById(adminUser.AdminUserId)
	if adminUserToken == nil {
		adminUserToken = AdminUserTokenService.GenerateToken(adminUser.AdminUserId)
	}else {
		AdminUserTokenService.UpdateToken(adminUserToken)
	}
	return adminUserToken.Token, nil
}

func (a *adminUser) GetProfileByToken (token string) (*manage.MallAdminUser, error) {
	adminToken, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return nil, err
	}
	adminUser := a.GetById(adminToken.AdminUserId)
	if adminUser == nil {
		return nil, errors.New("账号已删除")
	}
	return adminUser, nil
}