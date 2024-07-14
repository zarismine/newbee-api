package mallservice

import (
	"newbee/global"
	"newbee/models/mall/response"
	"newbee/models/manage"
	"newbee/pkg/stringopt"
)

var MallGoodsService = newMallGoodsService()

func newMallGoodsService() *mallGoodsService {
	return &mallGoodsService{}
}

type mallGoodsService struct {
}

func (m *mallGoodsService) Take (where ...interface{}) *manage.MallGoodsInfo {
	ret := &manage.MallGoodsInfo{}
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (m *mallGoodsService) SearchByCategory(pageNumber, pageSize, goodsCategory int, keyWord, orderBy, token string) ([]response.GoodsResponse, int64, error) {
	_, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil, 0, err
	}
	db := global.DB.Model(&manage.MallGoodsInfo{})
	if keyWord != "" {
		db = db.Where("goods_name like ? or goods_intro like ? or tag like ?", "%" + keyWord + "%", "%" + keyWord + "%",  "%" + keyWord + "%")
	}
	if goodsCategory != 0 {
		db = db.Where("goods_category_id = ?", goodsCategory)
	}
	var total int64
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	switch orderBy {
	case "new":
		db = db.Order("goods_id desc")
	case "price":
		db = db.Order("selling_price asc")
	default:
		db = db.Order("goods_rank desc")
	}
	goodsList := []manage.MallGoodsInfo{}
	searchGoodsList := []response.GoodsResponse{}
	limit := pageSize
	offset := limit * (pageNumber - 1)
	err = db.Limit(limit).Offset(offset).Find(&goodsList).Error
	for _, goods := range goodsList {
		searchGoods := response.GoodsResponse{
			GoodsId:       goods.GoodsId,
			GoodsName:     stringopt.SubStrLen(goods.GoodsName, 28),
			GoodsIntro:    stringopt.SubStrLen(goods.GoodsIntro, 28),
			GoodsCoverImg: goods.GoodsCoverImg,
			SellingPrice:  goods.SellingPrice,
		}
		searchGoodsList = append(searchGoodsList, searchGoods)
	}
	return searchGoodsList, total, err
}

func (m *mallGoodsService) DetailByGoodsId (id int, token string) (*response.GoodsDetailResponse, error) {
	_, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil, err
	}
	goods := new(manage.MallGoodsInfo)
	err = global.DB.Where("goods_id = ?",id).First(goods).Error
	if err != nil {
		return nil,err
	}
	goodsDetailResponse := &response.GoodsDetailResponse{
		GoodsCarouselList : []string{goods.GoodsCoverImg},
		GoodsCoverImg     :goods.GoodsCoverImg ,
		GoodsDetailContent: goods.GoodsDetailContent,
		GoodsId: goods.GoodsId,
		GoodsName: goods.GoodsName,
		GoodsIntro: goods.GoodsIntro,
		SellingPrice: goods.SellingPrice,
		Tag: goods.Tag,
	}
	return goodsDetailResponse,nil
}