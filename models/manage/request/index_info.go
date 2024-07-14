package request

import "newbee/models/manage"

type MallIndexConfigSearch struct {
	manage.MallIndexConfig
	PageNumber int `json:"pageNumber" form:"pageNumber"` // 页码
	PageSize   int `json:"pageSize" form:"pageSize"`     // 每页大小
}

type MallIndexConfigAddParams struct {
	ConfigName  string `json:"configName"`
	ConfigType  int    `json:"configType"`
	GoodsId     int    `json:"goodsId"`
	RedirectUrl string `json:"redirectUrl"`
	ConfigRank  int    `json:"configRank"`
}

type MallIndexConfigUpdateParams struct {
	ConfigId    int    `json:"configId"`
	ConfigName  string `json:"configName"`
	RedirectUrl string `json:"redirectUrl"`
	ConfigType  int    `json:"configType"`
	GoodsId     int    `json:"goodsId"`
	ConfigRank  int    `json:"configRank"`
}