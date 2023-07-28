package router

import (
	"cloudrestaurant/controller"
	"cloudrestaurant/middlewares"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	// 调用跨域请求中间件
	r.Use(middlewares.Cors())

	// 初始化session存储引擎
	//tool.InitSessionStore(r)

	// 注册路由
	RegisterRouter(r)

	return r
}

// RegisterRouter 路由设置
func RegisterRouter(engine *gin.Engine) {
	new(controller.MemberController).MemberRouter(engine)
	new(controller.ProductController).ProductRouter(engine)
	new(controller.FavoriteController).FavoriteRouter(engine)
	new(controller.OrderController).OrderRouter(engine)
	new(controller.AddressController).AdderssRouter(engine)
	new(controller.MoneyController).PayRouter(engine)
}
