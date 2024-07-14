package request

type OrdersReq struct {
	PageNumber    int      `json:"pageNumber" form:"pageNumber"`
	PageSize      int      `json:"pageSize" form:"pageSize"`     
	OrderNo       string   `json:"orderNo" form:"orderNo"`
	OrderStatus   string   `json:"orderStatus" form:"orderStatus"`
}