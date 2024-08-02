package mall

import "newbee/models/jsontime"

type MallMessage struct {
	MessageId  int               `json:"messageId,omitempty" form:"messageId" gorm:"primarykey;AUTO_INCREMENT"`
	SendId     int               `json:"sendId" form:"sendId" gorm:"column:send_id;comment:发送方id;type:int"`
	RecvId     int               `json:"recvId" form:"recvId" gorm:"column:recv_id;comment:接受方id;type:int"`
	Type       int               `json:"type" form:"type" gorm:"column:type;comment:数据类型(-1-心跳 0-文本信息 1-图片 2-语音 3-红包);type:int"`
	Content    string            `json:"content" form:"content" gorm:"column:content;comment:消息内容;type:varchar(1000)"`
	CreateTime jsontime.JSONTime `gorm:"column:create_time;comment:创建时间;type:datetime"`
}

func (MallMessage) TableName() string {
	return "tb_newbee_mall_message"
}