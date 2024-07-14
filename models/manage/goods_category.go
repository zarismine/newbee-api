package manage

import "newbee/models/jsontime"

type MallGoodsCategory struct {
	CategoryId    int                `json:"categoryId" gorm:"primarykey;AUTO_INCREMENT"`
	UserId        int                `json:"UserId" gorm:"comment:管理员ID"`
	CategoryLevel int                `json:"categoryLevel" gorm:"comment:分类等级"`
	ParentId      int                `json:"parentId" gorm:"comment:父类id"`
	CategoryName  string             `json:"categoryName" gorm:"comment:分类名称"`
	CategoryRank  int                `json:"categoryRank" gorm:"comment:排序比重"`
	IsDeleted     int                `json:"isDeleted" gorm:"comment:是否删除"`
	CreateTime    jsontime.JSONTime  `json:"createTime" gorm:"column:create_time;comment:创建时间;type:datetime"` // 创建时间
	UpdateTime    jsontime.JSONTime  `json:"updateTime" gorm:"column:update_time;comment:修改时间;type:datetime"` // 更新时间
}

func (MallGoodsCategory) TableName() string {
	return "tb_newbee_mall_goods_category"
}
