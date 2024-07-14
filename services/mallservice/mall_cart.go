package mallservice

import (
	"errors"
	"newbee/global"
	"newbee/models/jsontime"
	"newbee/models/mall"
	"newbee/models/mall/response"
	"time"
)

var MallCartService = newMallCartService()

func newMallCartService() *mallCartService {
	return &mallCartService{}
}

type mallCartService struct {
}

func (m *mallCartService) Take (where ...interface{}) *mall.MallShoppingCartItem {
	ret := &mall.MallShoppingCartItem{}
	if err := global.DB.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (m *mallCartService) GetShoppingCartItemByToken(token string) ([]response.CartItemResponse, error) {
	usertoken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil, errors.New("不存在的用户")
	}
	var cartItems []mall.MallShoppingCartItem
	global.DB.Where("user_id = ? and is_deleted = 0", usertoken.UserId).Find(&cartItems)
	if len(cartItems) == 0 {
		return []response.CartItemResponse{}, nil
	}
	var res []response.CartItemResponse
	for _, v := range(cartItems) {
		goods := MallGoodsService.Take("goods_id = ?", v.GoodsId)
		res_temp := response.CartItemResponse{
			CartItemId: v.CartItemId,
			GoodsCount: v.GoodsCount,
			GoodsCoverImg: goods.GoodsCoverImg,
			GoodsId: v.GoodsId,
			GoodsName: goods.GoodsName,
			SellingPrice: goods.SellingPrice,
		}
		res = append(res, res_temp)
	}
	return res, nil
}

func (m *mallCartService) EditShoppingCartItem(goodsCount, cartItemId int, token string) error {
	usertoken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return errors.New("不存在的用户")
	}
	cartItem := m.Take("cart_item_id = ? and is_deleted = 0", cartItemId)
	if usertoken.UserId != cartItem.UserId {
		return errors.New("禁止的操作")
	}
	if cartItem == nil {
		return errors.New("无效的订单")
	}
	return global.DB.Model(cartItem).UpdateColumns(map[string]interface{} {
		"goods_count" : goodsCount,
		"update_time" : jsontime.JSONTime{Time: time.Now()},
	}).Error
}

func (m *mallCartService) DeleteCartItemById(cartItemId int, token string) error {
	usertoken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return errors.New("不存在的用户")
	}
	cartItem := m.Take("cart_item_id = ? and is_deleted = 0", cartItemId)
	if cartItem == nil {
		return errors.New("无效的订单")
	}
	if usertoken.UserId != cartItem.UserId {
		return errors.New("禁止的操作")
	}
	return global.DB.Model(cartItem).UpdateColumns(map[string]interface{} {
		"is_deleted"  : 1,
		"update_time" : jsontime.JSONTime{Time: time.Now()},
	}).Error
}

func (m *mallCartService) AddShoppingCartItem(goodsId, goodsCount int, token string) error {
	usertoken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return errors.New("不存在的用户")
	}
	if m.Take("user_id = ? and goods_id = ? and is_deleted = 0", usertoken.UserId, goodsId) != nil {
		return errors.New("订单已存在")
	}
	if MallGoodsService.Take("goods_id = ?", goodsId) == nil {
		return errors.New("商品不存在")
	}
	var total int64
	global.DB.Where("user_id = ?  and is_deleted = 0", usertoken.UserId).Count(&total)
	if total > 20 {
		return errors.New("超出购物车最大容量！")
	}
	shopCartItem := &mall.MallShoppingCartItem{
		UserId     :   usertoken.UserId,
		GoodsId    :   goodsId,
		GoodsCount :   goodsCount,
		CreateTime :   jsontime.JSONTime{Time: time.Now()},
		UpdateTime :   jsontime.JSONTime{Time: time.Now()},
	}
	return global.DB.Save(shopCartItem).Error
}

func (m *mallCartService) GetShoppingCartItemById(cartItemIds []int, token string) ([]response.CartItemResponse, error) {
	usertoken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil, errors.New("不存在的用户")
	}
	var cartItems []mall.MallShoppingCartItem
	global.DB.Where("cart_item_id in ? and is_deleted = 0 and user_id = ?", cartItemIds, usertoken.UserId).Find(&cartItems)
	if len(cartItems) == 0 {
		return []response.CartItemResponse{}, nil
	}
	var res []response.CartItemResponse
	for _, v := range(cartItems) {
		goods := MallGoodsService.Take("goods_id = ?", v.GoodsId)
		res_temp := response.CartItemResponse{
			CartItemId: v.CartItemId,
			GoodsCount: v.GoodsCount,
			GoodsCoverImg: goods.GoodsCoverImg,
			GoodsId: v.GoodsId,
			GoodsName: goods.GoodsName,
			SellingPrice: goods.SellingPrice,
		}
		res = append(res, res_temp)
	}
	return res, nil
}