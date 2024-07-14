package adminservice

import (
	"errors"
	"newbee/global"
	"newbee/models/mall"
	"newbee/models/manage/request"
)

var MallUserService = newMallUserService()

func newMallUserService() *mallUserService {
	return &mallUserService{}
}

type mallUserService struct {
}

func (m *mallUserService) LockUser(ids []int, lockStatus int, token string) (err error) {
	_, err = AdminUser.GetProfileByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	if lockStatus != 0 && lockStatus != 1 {
		return errors.New("操作非法！")
	}
	err = global.DB.Table("tb_newbee_mall_user").Where("user_id in ?", ids).Update("locked_flag", lockStatus).Error
	return err
}

func (m *mallUserService) GetMallUserInfoList(info *request.PageInfo,token string) (list interface{}, total int64, err error) {
	_, err = AdminUser.GetProfileByToken(token)
	if err != nil {
		return nil, 0, errors.New("无效token")
	}
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	db := global.DB.Table("tb_newbee_mall_user")
	var mallUsers []mall.MallUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("create_time desc").Find(&mallUsers).Error
	return mallUsers, total, err
}