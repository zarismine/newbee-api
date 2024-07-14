package request

type EditCart struct {
	CartItemId int  `json:"cartItemId"  form:"cartItemId"`
	GoodsCount int  `json:"goodsCount"  form:"goodsCount"`
}

type AddCart struct {
	GoodsId    int  `json:"goodsId"  form:"goodsId"`
	GoodsCount int  `json:"goodsCount"  form:"goodsCount"`
}