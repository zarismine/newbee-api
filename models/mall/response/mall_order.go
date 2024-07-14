package response

import "newbee/models/jsontime"

type MallOrderResponse struct {
	OrderId                int                `json:"orderId"`
	OrderNo                string             `json:"orderNo"`
	TotalPrice             int                `json:"totalPrice"`
	PayType                int                `json:"payType"`
	OrderStatus            int                `json:"orderStatus"`
	OrderStatusString      string             `json:"orderStatusString"`
	CreateTime             jsontime.JSONTime  `json:"createTime"`
	NewBeeMallOrderItemVOS []CartItemResponse `json:"newBeeMallOrderItemVOS"`
}

