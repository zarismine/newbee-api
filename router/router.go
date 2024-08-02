package router

import (
	// "fmt"
	"net"
	"newbee/controller/admin"
	"newbee/controller/api"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
)

func isAllowedOrigin(origin string, allowedCidrs []string) bool {
    // 解析请求的IP地址
    ip, _, err := net.SplitHostPort(strings.TrimPrefix(origin, "http://"))
    if err != nil {
        return false
    }
    // 检查IP是否在允许的CIDR范围内
    for _, cidr := range allowedCidrs {
        _, subnet, err := net.ParseCIDR(cidr)
        if err != nil {
            continue
        }
        if subnet.Contains(net.ParseIP(ip)) {
            return true
        }
    }
    return false
}

// 方法二
func crs(ctx iris.Context) {
	// allowed_origins := []string{"192.168.0.0/16","127.0.0.1/32","121.37.0.0/16"}
	origin := ctx.GetHeader("Origin")
	// fmt.Println(origin)
	// if isAllowedOrigin(origin, allowed_origins) {
	// 	fmt.Println("yes")

	// }
	ctx.Header("Access-Control-Allow-Origin", origin)
	ctx.Header("Access-Control-Allow-Credentials", "true")
	if ctx.Method() == iris.MethodOptions {
		ctx.Header("Access-Control-Allow-Methods", "POST, PUT, PATCH, DELETE")
		ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type,X-Requested-With,Token")
		ctx.Header("Access-Control-Max-Age", "86400")
		ctx.StatusCode(iris.StatusNoContent)
		return
	}
	ctx.Next()
}

func NewServer() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.UseRouter(crs)

	mvc.Configure(app.Party("/api"), func(context *mvc.Application) {
		context.Party("/user").Handle(new(api.MallUserController))
		context.Party("/address").Handle(new(api.MallUserAddressController))
		context.Party("/categories").Handle(new(api.MallCategoryController))
		context.Party("/goods").Handle(new(api.MallGoodsController))
		context.Party("/shop-cart").Handle(new(api.MallCartController))
		context.Party("/order").Handle(new(api.MallOrderController))
		context.Party("/index-infos").Handle(new(api.MallIndexInfoController))
		context.Party("/chat").Handle(new(api.ChatController))
		context.Party("/contact").Handle(new(api.ContactController))
	})

	mvc.Configure(app.Party("/api/admin"), func(context *mvc.Application) {
		context.Party("/user").Handle(new(admin.AdminUserController))
		context.Party("/categories").Handle(new(admin.GoodCategotyController))
		context.Party("/goods").Handle(new(admin.GoodsController))
		context.Party("/users").Handle(new(admin.UserController))
		context.Party("/order").Handle(new(admin.OrderController))
		context.Party("/indexConfigs").Handle(new(admin.IndexInfoController))
	})
	app.Run(iris.Addr(":8081"))
}