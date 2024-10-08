package request

type MallAdminLoginParam struct {
	UserName    string    `json:"userName"`
	PasswordMd5 string    `json:"passwordMd5"`
}
type MallAdminRegisterParam struct {
	LoginUserName string     `json:"loginUserName"`
	LoginPassword string     `json:"loginPassword"`
	NickName      string     `json:"nickName"`
	KEY           string     `json:"KEY"`
}

type MallUpdateNameParam struct {
	LoginUserName string     `json:"loginUserName"`
	NickName      string     `json:"nickName"`
}

type MallUpdatePasswordParam struct {
	OriginalPassword string     `json:"originalPassword"`
	NewPassword      string     `json:"newPassword"`
}