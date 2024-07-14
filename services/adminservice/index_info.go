package adminservice

import (
	"errors"
	"newbee/global"
	"newbee/models/jsontime"
	"newbee/models/manage"
	"newbee/models/manage/request"
	"time"

	"gorm.io/gorm"
)

var IndexInfoService = newIndexInfoService()

func newIndexInfoService() *indexInfoService {
	return &indexInfoService{}
}

type indexInfoService struct {
}

func (m indexInfoService) Take(where ...interface{}) *manage.MallIndexConfig {
	ret := new(manage.MallIndexConfig)
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (m *indexInfoService) CreateMallIndexConfig(req request.MallIndexConfigAddParams, token string) (err error) {
	adminUserToken, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	if goods := GoodsService.Take("goods_id = ?", req.GoodsId); goods == nil {
		return errors.New("商品不存在")
	}
	if configItem := m.Take("config_type = ? and goods_id = ? and is_deleted = 0", req.ConfigType, req.GoodsId); configItem != nil {
		return errors.New("已存在相同的首页配置项")
	}

	mallIndexConfig := manage.MallIndexConfig{
		ConfigName:  req.ConfigName,
		ConfigType:  req.ConfigType,
		GoodsId:     req.GoodsId,
		RedirectUrl: req.RedirectUrl,
		ConfigRank:  req.ConfigRank,
		CreateTime:  jsontime.JSONTime{Time: time.Now()},
		UpdateTime:  jsontime.JSONTime{Time: time.Now()},
		UpdateUser: adminUserToken.AdminUserId,
	}
	err = global.DB.Create(&mallIndexConfig).Error
	return err
}

// DeleteMallIndexConfig 删除MallIndexConfig记录
func (m *indexInfoService) DeleteMallIndexConfig(ids request.IdsReq, token string) error {
	adminUserToken, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	err = global.DB.Where("config_id in ?", ids.Ids).UpdateColumns(manage.MallIndexConfig{IsDeleted: 1,UpdateTime: jsontime.JSONTime{Time: time.Now()}, UpdateUser: adminUserToken.AdminUserId}).Error
	return err
}

// UpdateMallIndexConfig 更新MallIndexConfig记录
func (m *indexInfoService) UpdateMallIndexConfig(req request.MallIndexConfigUpdateParams, token string) (err error) {
	adminUserToken, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	if errors.Is(global.DB.Where("goods_id = ?", req.GoodsId).First(&manage.MallGoodsInfo{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("商品不存在！")
	}
	if errors.Is(global.DB.Where("config_id=?", req.ConfigId).First(&manage.MallIndexConfig{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("未查询到记录！")
	}
	mallIndexConfig := manage.MallIndexConfig{
		ConfigId:    req.ConfigId,
		ConfigType:  req.ConfigType,
		ConfigName:  req.ConfigName,
		RedirectUrl: req.RedirectUrl,
		GoodsId:     req.GoodsId,
		ConfigRank:  req.ConfigRank,
		UpdateTime:  jsontime.JSONTime{Time: time.Now()},
		UpdateUser:  adminUserToken.AdminUserId,
	}
	var newIndexConfig manage.MallIndexConfig
	err = global.DB.Where("config_type=? and goods_id=?", mallIndexConfig.ConfigType, mallIndexConfig.GoodsId).First(&newIndexConfig).Error
	if err != nil && newIndexConfig.ConfigId == mallIndexConfig.ConfigId {
		return errors.New("已存在相同的首页配置项")
	}
	err = global.DB.Where("config_id = ?", mallIndexConfig.ConfigId).Updates(&mallIndexConfig).Error
	return err
}

// GetMallIndexConfig 根据id获取MallIndexConfig记录
func (m *indexInfoService) GetMallIndexConfig(id uint,token string) (mallIndexConfig manage.MallIndexConfig, err error) {
	_, err = AdminUserTokenService.GetByToken(token)
	if err != nil {
		return manage.MallIndexConfig{}, errors.New("无效token")
	}
	err = global.DB.Where("config_id = ?", id).First(&mallIndexConfig).Error
	return
}

// GetMallIndexConfigInfoList 分页获取MallIndexConfig记录
func (m *indexInfoService) GetMallIndexConfigInfoList(info request.MallIndexConfigSearch, token string) (list interface{}, total int64, err error) {
	_, err = AdminUserTokenService.GetByToken(token)
	if err != nil {
		return nil, 0, errors.New("无效token")
	}
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	db := global.DB.Model(&manage.MallIndexConfig{})
	db = db.Where("is_deleted = 0")
	if info.ConfigType != 0 {
		db = db.Where("config_type=?", info.ConfigType)
	}
	var mallIndexConfigs []manage.MallIndexConfig
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("config_rank desc").Find(&mallIndexConfigs).Error
	return mallIndexConfigs, total, err
}