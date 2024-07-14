package adminservice

import (
	"errors"
	"newbee/global"
	"newbee/models/jsontime"
	"newbee/models/manage"
	"newbee/models/manage/request"
	"newbee/pkg/verfiy"
	"time"
)

var GoodsCategoryService = newGoodsCategoryService()

func newGoodsCategoryService() *goodsCategoryService {
	return &goodsCategoryService{}
}

type goodsCategoryService struct {
}

func (gc *goodsCategoryService) Create (goodsCategory *manage.MallGoodsCategory) error {
	return global.DB.Create(goodsCategory).Error
}

func (gc *goodsCategoryService) Save (goodsCategory *manage.MallGoodsCategory) error {
	return global.DB.Save(goodsCategory).Error
}

func (gc *goodsCategoryService) Take (where ...interface{}) *manage.MallGoodsCategory {
	ret := new(manage.MallGoodsCategory)
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (gc *goodsCategoryService) AddGoodsCategory (parentId, level, rank int, name, token string) error {
	adminUser, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	if gc.Take("parent_id = ? AND category_name = ? AND is_deleted = 0",parentId, name) != nil {
		return errors.New("存在相同分类")
	}
	category := &manage.MallGoodsCategory{
		CategoryLevel:  level,
		CategoryName :  name,
		CategoryRank :  rank,
		UserId       :  adminUser.AdminUserId,
		IsDeleted    :  0,
		ParentId     :  parentId,
		CreateTime   :  jsontime.JSONTime{Time: time.Now()},
		UpdateTime   :  jsontime.JSONTime{Time: time.Now()},
	}
	if !verfiy.Contains([]int{1,2,3,4},level) {
		return errors.New("参数有误")
	}
	return gc.Create(category)
}

func (gc *goodsCategoryService) UpdateGoodsCategory (categoryId, level, rank int, name, token string) error {
	_, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	goodsCategory := gc.Take("category_id = ?", categoryId)
	if goodsCategory.CategoryName == name && goodsCategory.CategoryRank == rank {
		return errors.New("未进行修改")
	}
	if goodsCategory.CategoryName != name && gc.Take("category_level = ? AND category_name = ? AND is_deleted = 0",level, name) != nil {
		return errors.New("存在相同分类")
	}
	return global.DB.Model(&manage.MallGoodsCategory{}).Where("category_id = ?",categoryId).UpdateColumns(map[string]interface{}{
		"category_name" : name,
		"category_rank" : rank,
		"update_time"   : jsontime.JSONTime{Time: time.Now()},
	}).Error
}

func (gc *goodsCategoryService) SearchCategory (req *request.SearchCategoryParams, token string) (int64, interface{}, error) {
	_, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return 0, nil, errors.New("无效token")
	}
	limit := req.PageSize
	if limit > 1000 {
		limit = 1000
	}
	offset := req.PageSize * (req.PageNumber - 1)
	db := global.DB.Model(&manage.MallGoodsCategory{})
	var categoryList []manage.MallGoodsCategory
	var total int64
	db = db.Where("category_level = ? AND parent_id = ? AND is_deleted = 0", req.CategoryLevel, req.ParentId)
	err = db.Count(&total).Error
	if err != nil {
		return total, categoryList, err
	} 
	db = db.Order("category_rank desc").Limit(limit).Offset(offset)
	err = db.Find(&categoryList).Error
	return total, categoryList, err
}

func (gc *goodsCategoryService) SearchCategoryById (id int, token string) (*manage.MallGoodsCategory, error) {
	_, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return nil, errors.New("无效token")
	}
	goodsCategory := gc.Take("category_id = ?", id)
	if goodsCategory == nil {
		return nil, errors.New("无效id")
	}
	return goodsCategory, nil
}

func (gc *goodsCategoryService) SearchCategoryLinkById (categoryId int, token string) (map[string]interface{}, error) {
	_, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return nil, errors.New("无效token")
	}
	categoryLink := make(map[string]interface{})
	if thirdCategory := gc.Take("category_id = ?", categoryId);thirdCategory != nil {
		categoryLink["thirdCategory"] = thirdCategory.CategoryName
		if secondCategory := gc.Take("category_id = ?", thirdCategory.ParentId);secondCategory != nil {
			categoryLink["secondCategory"] = secondCategory.CategoryName
			if firstCategory := gc.Take("category_id = ?", secondCategory.ParentId);firstCategory != nil {
				categoryLink["firstCategory"] = firstCategory.CategoryName
				return categoryLink, nil
			}
		}
	}
	return nil, nil
}

func (gc *goodsCategoryService) DeleteById (ids []int,token string) error {
	adminUser, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	return global.DB.Where("category_id in ?",ids).UpdateColumns(manage.MallGoodsCategory{IsDeleted: 1,UpdateTime: jsontime.JSONTime{Time: time.Now()}, UserId: adminUser.AdminUserId}).Error
}