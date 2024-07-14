package request

type MallGoodsCategoryReq struct {
	CategoryId    int             `json:"categoryId"`
	CategoryLevel int             `json:"categoryLevel" `
	ParentId      int             `json:"parentId"`
	CategoryName  string          `json:"categoryName" `
	CategoryRank  string          `json:"categoryRank" `
}

type SearchCategoryParams struct {
	CategoryLevel int   `json:"categoryLevel" form:"categoryLevel"`
	ParentId      int   `json:"parentId" form:"parentId"`
	PageNumber    int   `json:"pageNumber" form:"pageNumber"` // 页码
	PageSize      int   `json:"pageSize" form:"pageSize"`     // 每页大小
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}