package request

type PageSearch struct {
	PageNumber           int      `json:"pageNumber" form:"pageNumber"`
	PageSize             int      `json:"pageSize" form:"pageSize"`     
	GoodsCategoryId      int      `json:"goodsCategoryId" form:"goodsCategoryId"`     
	Keyword              string   `json:"keyword" form:"keyword"`
	OrderBy              string   `json:"orderBy" form:"orderBy"`
}