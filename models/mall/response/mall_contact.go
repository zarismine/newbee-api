package response

// import "newbee/models/jsontime"

type ContactResponse struct {
	// HeadImg        string `json:"headImg" form:"headImg"`
	UserId         int    `json:"userId" form:"userId"`
	MessageContent string `json:"messageContent" form:"messageContent"`
	NickName       string `json:"nickName" form:"nickName"`
	MessageTime    string `json:"messageTime" form:"messageTime"`
	Count          int64  `json:"count" form:"count"`
}
