package mallservice

import (
	"errors"
	"newbee/global"
	"newbee/models/constants"
	"newbee/models/jsontime"
	"newbee/models/mall"
	"newbee/models/mall/response"
	"newbee/models/manage"
	"newbee/pkg/passwd"
	"time"

	"github.com/jinzhu/copier"
)

var MallOrderService = newMallOrderService()

func newMallOrderService() *mallOrderService {
	return &mallOrderService{}
}

type mallOrderService struct {
}

type goodsItem struct {
	SellingPrice   int
	GoodsName      string
	GoodsCoverImg  string
}
func (m *mallOrderService) Save(cartItemIds []int, addressId int, token string) (string, error) {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return "", err
	}
	var cartItems []mall.MallShoppingCartItem
	mapIdGoods := make(map[int]goodsItem)
	global.DB.Where("cart_item_id in ? and is_deleted = 0", cartItemIds).Find(&cartItems)
	if len(cartItems) == 0 {
		return "", errors.New("购物车数据异常！")
	}
	for _, val := range(cartItems) {
		if val.UserId != userToken.UserId {
			return "", errors.New("该用户无权限")
		}
		goods := MallGoodsService.Take("goods_id = ?", val.GoodsId)
		mapIdGoods[val.GoodsId] = goodsItem{SellingPrice :goods.SellingPrice,GoodsName : goods.GoodsName, GoodsCoverImg: goods.GoodsCoverImg}
		if goods.GoodsSellStatus == 1 || goods.StockNum < val.GoodsCount{
			return "", errors.New("商品已下架或库存不足订单失败")
		}

	}
	if err = global.DB.Model(&cartItems).Updates(map[string]interface{} {
		"is_deleted"  : 1,
		"update_time" : jsontime.JSONTime{Time: time.Now()},
	}).Error;err != nil {
		return "", err
	}
	priceTotal := 0
	orderNo := passwd.GenOrderNo()
	var newBeeMallOrder manage.MallOrder
	newBeeMallOrder.OrderNo = orderNo
	newBeeMallOrder.UserId = userToken.UserId
	for _, val := range cartItems {
		priceTotal = priceTotal + val.GoodsCount*mapIdGoods[val.GoodsId].SellingPrice
	}
	if priceTotal < 1 {
		return "", errors.New("订单价格异常！")
	}
	newBeeMallOrder.CreateTime = jsontime.JSONTime{Time: time.Now()}
	newBeeMallOrder.UpdateTime = jsontime.JSONTime{Time: time.Now()}
	newBeeMallOrder.TotalPrice = priceTotal
	newBeeMallOrder.ExtraInfo = ""
	if err = global.DB.Save(&newBeeMallOrder).Error; err != nil {
		return "", errors.New("订单入库失败！")
	}
	var newBeeMallOrderItems []manage.MallOrderItem
	for _, val := range cartItems {
		OrderItem := manage.MallOrderItem{
			OrderId: newBeeMallOrder.OrderId,
			GoodsId: val.GoodsId,
			GoodsName: mapIdGoods[val.GoodsId].GoodsName,
			GoodsCoverImg: mapIdGoods[val.GoodsId].GoodsCoverImg,
			SellingPrice: mapIdGoods[val.GoodsId].SellingPrice,
			GoodsCount: val.GoodsCount,
			AddressId: addressId,
			CreateTime: jsontime.JSONTime{Time: time.Now()},
		}
		newBeeMallOrderItems = append(newBeeMallOrderItems, OrderItem)
	}
	return orderNo, global.DB.Save(&newBeeMallOrderItems).Error
}

func (m *mallOrderService) GetOrderList(token, status string) ([]*response.MallOrderResponse, int64, error) {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return []*response.MallOrderResponse{}, 0, errors.New("无效的token")
	}
	db := global.DB.Table("tb_newbee_mall_order")
	if status != "" {
		db = db.Where("order_status = ?", status)
	}
	db = db.Where("user_id = ? and is_deleted = 0", userToken.UserId)
	var total int64
	err = db.Count(&total).Error
	if total == 0{
		return []*response.MallOrderResponse{}, total, err
	}
	var orders []*manage.MallOrder
	limit := 5
	pageNumber := 1
	offset := limit * (pageNumber - 1)
	err = db.Offset(offset).Order("update_time desc").Find(&orders).Error
	if len(orders) == 0 {
		return []*response.MallOrderResponse{}, total, err
	}
	var resp []*response.MallOrderResponse
	copier.Copy(&resp, &orders)
	for _, val := range(resp) {
		val.OrderStatusString = constants.Status[val.OrderStatus]
		var orderitems []manage.MallOrderItem
		global.DB.Table("newbee.tb_newbee_mall_order_item").Where("order_id = ?", val.OrderId).Find(&orderitems)
		if len(orderitems) == 0 {
			return []*response.MallOrderResponse{}, total, errors.New("订单异常！")
		}
		var cartItem []response.CartItemResponse
		copier.Copy(&cartItem, &orderitems)
		val.NewBeeMallOrderItemVOS = cartItem
	}
	return resp, total, nil

}

func (m *mallOrderService) GetOrderItemByorderNo (orderNo, token string) (*response.MallOrderResponse, error) {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return nil, errors.New("无效的token")
	}
	order := new(manage.MallOrder)
	db := global.DB.Table("tb_newbee_mall_order")
	db.Where("order_no = ? and is_deleted = 0 ", orderNo).Find(order)
	if order.UserId != userToken.UserId {
		return nil, errors.New("无权限")
	}
	resp := new(response.MallOrderResponse)
	copier.Copy(resp, order)
	resp.OrderStatusString = constants.Status[resp.OrderStatus]
	var orderitems []*manage.MallOrderItem
	global.DB.Table("newbee.tb_newbee_mall_order_item").Where("order_id = ?", resp.OrderId).Find(&orderitems)
	if len(orderitems) == 0 {
		return nil, errors.New("订单异常！")
	}
	var cartItem []response.CartItemResponse
	copier.Copy(&cartItem, &orderitems)
	resp.NewBeeMallOrderItemVOS = cartItem
	return resp, nil
}

func (m *mallOrderService) PaySuccess(orderNo, token string, payType int) error {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return errors.New("无效的token")
	}
	order := new(manage.MallOrder)
	db := global.DB.Table("tb_newbee_mall_order")
	db.Where("order_no = ? and is_deleted = 0 ", orderNo).Find(order)
	if order.UserId != userToken.UserId {
		return errors.New("无权限")
	}
	return global.DB.Model(order).Updates(map[string]interface{} {
		"pay_status"   : 1,
		"update_time"  : jsontime.JSONTime{Time: time.Now()},
		"pay_type"     : payType,
		"order_status" : 1,
		"pay_time"     : jsontime.JSONTime{Time: time.Now()},
	}).Error
}

func (m *mallOrderService) FinishOrder(orderNo, token string) error {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return errors.New("无效的token")
	}
	order := new(manage.MallOrder)
	db := global.DB.Table("tb_newbee_mall_order")
	db.Where("order_no = ? and is_deleted = 0 ", orderNo).Find(order)
	if order.UserId != userToken.UserId {
		return errors.New("无权限")
	}
	return global.DB.Model(order).Updates(map[string]interface{} {
		"update_time"  : jsontime.JSONTime{Time: time.Now()},
		"order_status" : 4,
	}).Error
}

func (m *mallOrderService) CancelOrder(orderNo, token string) error {
	userToken, err := MallUserTokenService.GetUserTokenByToken(token)
	if err != nil {
		return errors.New("无效的token")
	}
	order := new(manage.MallOrder)
	db := global.DB.Table("tb_newbee_mall_order")
	db.Where("order_no = ? and is_deleted = 0", orderNo).Find(order)
	if order.UserId != userToken.UserId {
		return errors.New("无权限")
	}
	return global.DB.Model(order).Updates(map[string]interface{} {
		"update_time"  : jsontime.JSONTime{Time: time.Now()},
		"order_status" : -1,
 		// "is_deleted"   : 1,
	}).Error
}