package adminservice

import (
	"errors"
	"newbee/global"
	"newbee/models/manage"
	"newbee/models/manage/request"
	"newbee/pkg/dates"
	"strconv"

	"gorm.io/gorm"
)

var GoodsService = newGoodsService()

func newGoodsService() *goodsService {
	return &goodsService{}
}

type goodsService struct {
}

func (g *goodsService) Create(goods *manage.MallGoodsInfo) error {
	return global.DB.Create(goods).Error
}

func (g *goodsService) Save(goods *manage.MallGoodsInfo) error {
	return global.DB.Save(goods).Error
}

func (g *goodsService) Take(where ...interface{}) *manage.MallGoodsInfo {
	ret := new(manage.MallGoodsInfo)
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (g *goodsService) AddGoods(req *request.GoodsInfoAddParam, token string) error {
	adminUser, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	goodsCategory := new(manage.MallGoodsCategory)
	err = global.DB.Where("category_id = ? AND is_deleted = 0", req.GoodsCategoryId).First(goodsCategory).Error
	if err != nil {
		return err
	}
	if goodsCategory.CategoryLevel != 3 {
		return errors.New("分类数据异常")
	}
	if !errors.Is(global.DB.Where("goods_name = ? AND goods_category_id = ?", req.GoodsName, req.GoodsCategoryId).First(&manage.MallGoodsInfo{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("已存在相同的商品信息")
	}
	originalPrice, _ := strconv.Atoi(req.OriginalPrice)
	sellingPrice, _ := strconv.Atoi(req.SellingPrice)
	stockNum, _ := strconv.Atoi(req.StockNum)
	goodsSellStatus, _ := strconv.Atoi(req.GoodsSellStatus)
	goodsRank, _ := strconv.Atoi(req.GoodsRank)
	goodsInfo := &manage.MallGoodsInfo{
		GoodsName           :        req.GoodsName,
		GoodsIntro          :        req.GoodsIntro,
		GoodsCategoryId     :        req.GoodsCategoryId,
		GoodsCoverImg       :        req.GoodsCoverImg,
		GoodsDetailContent  :        req.GoodsDetailContent,
		OriginalPrice       :        originalPrice,
		SellingPrice        :        sellingPrice,
		StockNum            :        stockNum,
		Tag                 :        req.Tag,
		GoodsRank           :        goodsRank,
		GoodsSellStatus     :        goodsSellStatus,
		UpdateUser          :        adminUser.AdminUserId,
		CreateTime          :        dates.NowTimestamp(),
		UpdateTime          :        dates.NowTimestamp(),
	}
	return g.Create(goodsInfo)
}

func (g *goodsService) UpdateGoods (req *request.GoodsInfoUpdateParam,token string) error {
	adminUser, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	goodsCategory := new(manage.MallGoodsCategory)
	err = global.DB.Where("category_id = ? AND is_deleted = 0", req.GoodsCategoryId).First(goodsCategory).Error
	if err != nil {
		return err
	}
	if goodsCategory.CategoryLevel != 3 {
		return errors.New("分类数据异常")
	}
	originalPrice, _ := strconv.Atoi(req.OriginalPrice)
	sellingPrice, _ := strconv.Atoi(req.SellingPrice)
	goodsId, _ := strconv.Atoi(req.GoodsId)
	stockNum, _ := strconv.Atoi(req.StockNum)
	goodsSellStatus, _ := strconv.Atoi(req.GoodsSellStatus)
	goodsRank, _ := strconv.Atoi(req.GoodsRank)
	goodsInfo := &manage.MallGoodsInfo{
		GoodsId             :        goodsId,
		GoodsName           :        req.GoodsName,
		GoodsIntro          :        req.GoodsIntro,
		GoodsCategoryId     :        req.GoodsCategoryId,
		GoodsCoverImg       :        req.GoodsCoverImg,
		GoodsDetailContent  :        req.GoodsDetailContent,
		OriginalPrice       :        originalPrice,
		SellingPrice        :        sellingPrice,
		StockNum            :        stockNum,
		Tag                 :        req.Tag,
		GoodsRank           :        goodsRank,
		GoodsSellStatus     :        goodsSellStatus,
		UpdateUser          :        adminUser.AdminUserId,
		// CreateTime          :        dates.NowTimestamp(),
		UpdateTime          :        dates.NowTimestamp(),
	}
	return g.Save(goodsInfo)
}
func (g *goodsService) GetGoodsList (info *request.PageInfo,token string,where ...interface{}) (interface{}, int64, error){
	_, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return nil, 0, errors.New("无效token")
	}
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	db := global.DB.Model(&manage.MallGoodsInfo{})
	var mallGoodsInfos []manage.MallGoodsInfo
	var total int64
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// if where != nil {
	// 	db = db.Where(where)
	// }
	err = db.Limit(limit).Offset(offset).Order("goods_rank desc").Find(&mallGoodsInfos).Error
	return mallGoodsInfos, total, err
}

func (g *goodsService) GetGoodsById (id int, token string) (*manage.MallGoodsInfo, error) {
	_, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return nil, errors.New("无效token")
	}
	goods := g.Take("goods_id = ?", id)
	if goods == nil {
		return nil, errors.New("无效id")
	}
	return goods, nil
}

func (g *goodsService) UpdateStatusByIds (status int,ids []int, token string) error {
	adminUser, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return errors.New("无效token")
	}
	return global.DB.Table("tb_newbee_mall_goods_info").Where("goods_id in ?",ids).UpdateColumns(map[string]interface{} {
		"goods_sell_status" : status,
		"update_time"       : dates.NowTimestamp(),
		"update_user"       : adminUser.AdminUserId,
	}).Error
}

func (g *goodsService) SearchGoodsList (info *request.PageInfo,token string) (interface{}, int64, error){
	_, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return nil, 0, errors.New("无效token")
	}
	limit := info.PageSize
	offset := info.PageSize * (info.PageNumber - 1)
	db := global.DB.Model(&manage.MallGoodsInfo{})
	db = db.Where("goods_name like ? COLLATE utf8_general_ci OR tag like ? COLLATE utf8_general_ci", "%"+info.SearchMsg+"%","%"+info.SearchMsg+"%")
	var mallGoodsInfos []manage.MallGoodsInfo
	var total int64
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Order("goods_rank desc").Find(&mallGoodsInfos).Error
	return mallGoodsInfos, total, err
}