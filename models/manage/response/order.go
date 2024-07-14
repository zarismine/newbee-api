package response

import "newbee/models/jsontime"

type OrderRes struct {
	OrderId     string            `json:"orderId" form:"orderId"`
	OrderNo     string            `json:"orderNo" form:"orderNo"`
	TotalPrice  int               `json:"totalPrice" form:"totalPrice"`
	PayStatus   int               `json:"payStatus" form:"payStatus"`
	PayType     int               `json:"payType" form:"payType"`
	OrderStatus int               `json:"orderStatus" form:"orderStatus"`
	IsDeleted   int               `json:"isDeleted" form:"isDeleted"`
	CreateTime  jsontime.JSONTime `json:"createTime" form:"createTime"`
}


type OrderDetailResponse struct {
	OrderNo                string             `json:"orderNo"`
	OrderStatusString      string             `json:"orderStatusString"`
	OrderExtraInfo         string             `json:"extraInfo"`
	CreateTime             jsontime.JSONTime  `json:"createTime"`
	NewBeeMallOrderItemVOS []GoodsResponse    `json:"newBeeMallOrderItemVOS"`
}