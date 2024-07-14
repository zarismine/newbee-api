package mall

type ThirdLevelCategoryVOS struct {
	CategoryId    int       `json:"categoryId"`
	CategoryLevel int       `json:"categoryLevel"`
	CategoryName  string    `json:"categoryName"`
	ParentId      int       `json:"parentId"`
}

type SecondLevelCategoryVOS struct {
	ThirdLevelCategoryVOS
	ThirdLevelCategory    []ThirdLevelCategoryVOS `json:"thirdLevelCategoryVOS"`
}

type FirstLevelCategoryVOS struct {
	ThirdLevelCategoryVOS
	SecondLevelCategory    []SecondLevelCategoryVOS `json:"secondLevelCategoryVOS"`
}