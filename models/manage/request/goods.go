package request

import (
	// "newbee/models/manage"
)

type PageInfo struct {
	PageNumber    int      `json:"pageNumber" form:"pageNumber"`
	PageSize      int      `json:"pageSize" form:"pageSize"`     
	SearchMsg     string   `json:"searchMsg" form:"searchMsg"`
}

type GoodsInfoAddParam struct {
	GoodsName          string     `json:"goodsName"`
	GoodsIntro         string     `json:"goodsIntro"`
	GoodsCategoryId    int        `json:"goodsCategoryId"`
	GoodsCoverImg      string     `json:"goodsCoverImg"`
	GoodsCarousel      string     `json:"goodsCarousel"`
	GoodsDetailContent string     `json:"goodsDetailContent"`
	OriginalPrice      string     `json:"originalPrice"`
	SellingPrice       string     `json:"sellingPrice"`
	StockNum           string     `json:"stockNum"`
	Tag                string     `json:"tag"`
	GoodsRank          string     `json:"goodsRank"`
	GoodsSellStatus    string     `json:"goodsSellStatus"`
}

// GoodsInfoUpdateParam 更新商品信息的入参
type GoodsInfoUpdateParam struct {
	GoodsId            string          `json:"goodsId"`
	GoodsName          string          `json:"goodsName"`
	GoodsIntro         string          `json:"goodsIntro"`
	GoodsCategoryId    int             `json:"goodsCategoryId"`
	GoodsCoverImg      string          `json:"goodsCoverImg"`
	GoodsCarousel      string          `json:"goodsCarousel"`
	GoodsDetailContent string          `json:"goodsDetailContent"`
	OriginalPrice      string          `json:"originalPrice"`
	SellingPrice       string          `json:"sellingPrice"`
	StockNum           string          `json:"stockNum"`
	Tag                string          `json:"tag"`
	GoodsRank          string          `json:"goodsRank"`
	GoodsSellStatus    string          `json:"goodsSellStatus"`
}

type StockNumDTO struct {
	GoodsId    int     `json:"goodsId"`
	GoodsCount int     `json:"goodsCount"`
}