package request

//用户注册
type RegisterUserParam struct {
	LoginName string `json:"loginname"`
	Password  string `json:"password"`
}

//更新用户信息
type UpdateUserInfoParam struct {
	NickName      string `json:"nickname"`
	PasswordMd5   string `json:"passwordmd5"`
	IntroduceSign string `json:"introducesign"`
}

type UserLoginParam struct {
	LoginName   string `json:"loginname"`
	PasswordMd5 string `json:"passwordmd5"`
}
