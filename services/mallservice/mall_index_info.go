package mallservice

import (
	"newbee/global"
	"newbee/models/constants"
	"newbee/models/mall/response"

	"github.com/jinzhu/copier"
)

var MallIndexInfoService = newMallIndexInfoService()

func newMallIndexInfoService() *indexInfoService {
	return &indexInfoService{}
}

type indexInfoService struct {
}

func (m *indexInfoService) GetDetailData() ([]*response.GoodsResponse, []*response.GoodsResponse, []*response.GoodsResponse, error) {
	hotGoodses := getGoodsDetailByConfigType(constants.HotCode)
	newGoodses := getGoodsDetailByConfigType(constants.NewCode)
	recommendGoodses := getGoodsDetailByConfigType(constants.RecommendCode)
	return hotGoodses, newGoodses, recommendGoodses, nil
}

func getGoodsDetailByConfigType(configtype int) []*response.GoodsResponse {
	var Goodses []*response.GoodsResponse
	global.DB.Table("newbee.tb_newbee_mall_index_config").Where("config_type = ? and is_deleted = 0", configtype).Find(&Goodses)
	if len(Goodses) == 0 {
		return []*response.GoodsResponse{}
	}
	for _,val := range(Goodses) {
		goods := MallGoodsService.Take("goods_id = ?",val.GoodsId)
		copier.Copy(val,goods)
	}
	return Goodses
}