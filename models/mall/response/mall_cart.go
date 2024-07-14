package response

type CartItemResponse struct {
	CartItemId     int     `json:"cartItemId" form:"cartItemId"`
	GoodsCount     int     `json:"goodsCount" form:"goodsCount"`
	GoodsCoverImg  string  `json:"goodsCoverImg" form:"goodsCoverImg"`
	GoodsId        int     `json:"goodsId" form:"goodsId"`
	GoodsName      string  `json:"goodsName" form:"goodsName"`
	SellingPrice   int     `json:"sellingPrice" form:"sellingPrice"`
}