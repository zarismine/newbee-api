package main // 声明 main 包，表明当前是一个可执行程序

import (
	"newbee/initialize"
	"newbee/router"
)

func main(){  // main函数，是程序执行的入口
	initialize.Init()  // 在终端打印 Hello World!
	router.NewServer()
}
