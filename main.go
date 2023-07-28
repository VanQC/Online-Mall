package main

import (
	"cloudrestaurant/cache"
	"cloudrestaurant/dao"
	"cloudrestaurant/router"
	"cloudrestaurant/tool"
	"log"
)

func main() {
	// 解析app.json文件
	cfg, err := tool.ParseConfig("./config/app_config.json")
	if err != nil {
		panic(err)
	}

	// 初始化mysql数据库连接
	if err = dao.InitDatabase(); err != nil {
		log.Println("初始化数据库错误，连接数据库失败。")
		return
	}

	// 初始化redis数据库连接
	cache.InitRedisStore()
	tool.InitLog()

	// 创建路由
	r := router.SetRouter()

	// 启动路由，端口号来自app.json文件
	if err = r.Run(cfg.AppHost + ":" + cfg.AppPort); err != nil {
		panic(err)
	}
}
