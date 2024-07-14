package constants

const (
	StatusOk = 1
	StatusDeleted = 0
)

var Status = map[int]string {
	0 : "待支付",
	1 : "已支付",
	2 :"配货完成",
	3 :"出库成功",
	4 :"交易成功",
	-1:"手动关闭",
	-2:"超时关闭",
	-3:"商家关闭",
}

const (
	HotCode = 3
	NewCode = 4
	RecommendCode = 5
)
