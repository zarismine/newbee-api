package response

type GoodsResponse struct {
	GoodsCoverImg string `json:"goodsCoverImg" form:"goodsCoverImg"`
	GoodsId       int    `json:"goodsId" form:"goodsId"`
	GoodsIntro    string `json:"goodsIntro" form:"goodsIntro"`
	GoodsName     string `json:"goodsName" form:"goodsName"`
	SellingPrice  int    `json:"sellingPrice" form:"sellingPrice"`
}

type GoodsDetailResponse struct {
	GoodsCarouselList   []string  `json:"goodsCarouselList" form:"goodsCarouselList"`
	GoodsCoverImg       string    `json:"goodsCoverImg" form:"goodsCoverImg"`
	GoodsDetailContent  string    `json:"goodsDetailContent" form:"goodsDetailContent"`
	GoodsId             int       `json:"goodsId" form:"goodsId"`
	GoodsName           string    `json:"goodsName" form:"goodsName"`
	GoodsIntro          string    `json:"goodsIntro" form:"goodsIntro"`
	SellingPrice        int       `json:"sellingPrice" form:"sellingPrice"`
	Tag                 string    `json:"tag"`
}
