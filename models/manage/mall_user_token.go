package manage

type MallAdminUserToken struct {
	AdminUserId int       `json:"adminUserId" form:"adminUserId" gorm:"primarykey;AUTO_INCREMENT"`
	Token       string    `json:"token" form:"token" gorm:"column:token;comment:token值(32位字符串);type:varchar(32);"`
	UpdateTime  int64     `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:修改时间;"`
	ExpireTime  int64     `json:"expireTime" form:"expireTime" gorm:"column:expire_time;comment:token过期时间;"`
}

func (MallAdminUserToken) TableName() string {
	return "tb_newbee_mall_admin_user_token"
}