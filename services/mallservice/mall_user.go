package mallservice

import (
	"errors"
	"newbee/global"
	"newbee/models/constants"
	"newbee/models/jsontime"
	"newbee/models/mall"
	"newbee/pkg/passwd"
	"newbee/pkg/verfiy"
	"strings"
	"time"
)

var MallUserService = newMallUserService()

func newMallUserService() *mallUserService {
	return &mallUserService{}
}

type mallUserService struct {
}

func (m *mallUserService) Get (id int) *mall.MallUser {
	ret := &mall.MallUser{}
	if err := global.DB.First(ret, "user_id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (m *mallUserService) Take (where ...interface{}) *mall.MallUser {
	ret := &mall.MallUser{}
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (m *mallUserService) GetById (Id int) *mall.MallUser {
	return m.Take("user_id = ?",Id)
}

func (m *mallUserService) GetByLoginName (LoginName string) *mall.MallUser {
	return m.Take("login_name = ?",LoginName)
}

func (m *mallUserService) Create(user *mall.MallUser) (err error) {
	err = global.DB.Create(user).Error
	return
}

func (m *mallUserService) Update(user *mall.MallUser) (err error) {
	err = global.DB.Updates(user).Error
	return
}

func (m *mallUserService) Delete(user *mall.MallUser) (err error) {
	err = global.DB.Delete(user).Error
	return
}

func (m *mallUserService) GetUserByToken(token string) (*mall.MallUser,error) {
	usertoken,err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil,err
	}
	return m.GetById(usertoken.UserId),nil
}

func (m *mallUserService) UpdateUser(nickname,passwordmd5,introducesign string,user *mall.MallUser) error {
	if nickname == "" {
		return errors.New("请输入昵称")
	}
	if passwordmd5 != "" {
		user.PasswordMd5 = passwordmd5
	}
	user.IntroduceSign = introducesign
	user.NickName = nickname
	return m.Update(user)
}

func (m *mallUserService) SignUp(loginname, nickname, password string) (*mall.MallUser, error) {
	loginname = strings.TrimSpace(loginname)
	nickname = strings.TrimSpace(nickname)
	if len(nickname) == 0 {
		nickname = "user_" + loginname
	}
	err := verfiy.IsPassword(password)
	if err != nil {
		return nil, err
	}
	if len(loginname) > 0 {
		if err := verfiy.IsUsername(loginname); err != nil {
			return nil, err
		}
		if m.GetByLoginName(loginname) != nil {
			return nil, errors.New("账号：" + loginname + " 已被占用")
		}
	}
	user := &mall.MallUser{
		NickName      :   nickname,
		LoginName     :   loginname,
		PasswordMd5   :   passwd.Hash([]byte(password)),
		IntroduceSign :   "随新所欲，蜂富多彩",
		CreateTime    :   jsontime.JSONTime{Time: time.Now()},
	}

	err = m.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *mallUserService) SignIn(username, password string) (*mall.MallUser, string, error) {
	if len(username) == 0 {
		return nil, "", errors.New("账号不能为空")
	}
	if err := verfiy.IsPassword(password); err != nil {
		return nil, "", err
	}
	var user *mall.MallUser = nil
	user = m.GetByLoginName(username)
	if user == nil || user.IsDeleted == constants.StatusOk {
		return nil, "", errors.New("账号不存在")
	}
	if password != user.PasswordMd5 {
		return nil, "", errors.New("账号或密码错误")
	}
	token := MallUserTokenService.GetById(user.UserId)
	if token == nil {
		token = MallUserTokenService.GenerateToken(user.UserId)
	}else {
		MallUserTokenService.UpdateToken(token)
	}
	return user, token.Token, nil
}