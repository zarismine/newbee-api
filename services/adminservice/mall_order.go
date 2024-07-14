package adminservice

import (
	"errors"
	"newbee/global"
	"newbee/models/constants"
	"newbee/models/jsontime"
	"newbee/models/manage"
	"newbee/models/manage/request"
	"newbee/models/manage/response"
	"time"
)

var MallOrderService = newMallOrderService()

func newMallOrderService() *mallOrderService {
	return &mallOrderService{}
}

type mallOrderService struct {
}

func (m *mallOrderService) GetOrderList(req *request.OrdersReq, token string) ([]*response.OrderRes, int64, error) {
	_, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return nil, 0, err
	}
	var resOrders []*response.OrderRes
	var total int64
	db := global.DB.Table("tb_newbee_mall_order")
	if req.OrderNo != "" {
		db = db.Where("order_no = ?",req.OrderNo)
	}
	if req.OrderStatus != "" {
		db = db.Where("order_status = ?",req.OrderStatus)
	}
	limit := req.PageSize
	offset := limit * (req.PageNumber - 1)
	db.Count(&total)
	db.Limit(limit).Offset(offset).Order("create_time desc").Find(&resOrders)
	return resOrders, total, nil
}

func (m *mallOrderService) OrderCheckdone(req *request.IdsReq, token string) error {
	adminUser, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return err
	}
	var Orders []*manage.MallOrder
	global.DB.Table("tb_newbee_mall_order").Where("order_id in ? and is_deleted = 0 and order_status = 1",req.Ids).Find(&Orders)
	if len(Orders) == 0{
		return errors.New("未找到订单信息")
	}
	for _, order := range(Orders) {
		global.DB.Model(order).UpdateColumns(map[string]interface{} {
			"order_status" : 2,
			"extra_info" : order.ExtraInfo + "管理员" + adminUser.NickName + "确认配货\n",
			"update_time" : jsontime.JSONTime{Time: time.Now()},
		})
	}
	return nil
}


func (m *mallOrderService) OrderCheckout(req *request.IdsReq, token string) error {
	adminUser, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return err
	}
	var Orders []*manage.MallOrder
	global.DB.Table("tb_newbee_mall_order").Where("order_id in ? and is_deleted = 0 and order_status = 2",req.Ids).Find(&Orders)
	if len(Orders) == 0{
		return errors.New("未找到订单信息")
	}
	for _, order := range(Orders) {
		global.DB.Model(order).UpdateColumns(map[string]interface{} {
			"order_status" : 3,
			"extra_info" : order.ExtraInfo + "管理员" + adminUser.NickName + "确认出库\n",
			"update_time" : jsontime.JSONTime{Time: time.Now()},
		})
	}
	return nil
}

func (m *mallOrderService) OrderClose(req *request.IdsReq, token string) error {
	adminUser, err := AdminUser.GetProfileByToken(token)
	if err != nil {
		return err
	}
	var Orders []*manage.MallOrder
	global.DB.Table("tb_newbee_mall_order").Where("order_id in ? and is_deleted = 0",req.Ids).Find(&Orders)
	if len(Orders) == 0{
		return errors.New("未找到订单信息")
	}
	for _, order := range(Orders) {
		global.DB.Model(order).UpdateColumns(map[string]interface{} {
			"order_status" : -3,
			"extra_info" : order.ExtraInfo + "管理员" + adminUser.NickName + "关闭订单\n",
			"update_time" : jsontime.JSONTime{Time: time.Now()},
		})
	}
	return nil
}

func (m *mallOrderService) OrderDetail(id int, token string) (*response.OrderDetailResponse, error) {
	_, err := AdminUserTokenService.GetByToken(token)
	if err != nil {
		return nil, err
	}
	order := new(manage.MallOrder)
	global.DB.Table("tb_newbee_mall_order").Where("order_id = ? and is_deleted = 0", id,).Find(order)
	if (order == &manage.MallOrder{}) {
		return nil, errors.New("订单不存在或已删除")
	}
	var goods []response.GoodsResponse
	global.DB.Table("tb_newbee_mall_order_item").Where("order_id = ?",id).Find(&goods)
	var orderDetail = &response.OrderDetailResponse {
		OrderNo: order.OrderNo,
		OrderStatusString: constants.Status[order.OrderStatus],
		OrderExtraInfo: order.ExtraInfo,
		CreateTime: order.CreateTime,
		NewBeeMallOrderItemVOS: goods,
	}
	return orderDetail, nil
}