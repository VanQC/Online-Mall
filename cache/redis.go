package cache

import (
	"cloudrestaurant/tool"
	"github.com/go-redis/redis"
	"log"
)

// 声明Redis客户端连接
var RedisClient *redis.Client

// InitRedisStore 建立Redis数据库连接
func InitRedisStore() {
	RedisConfig := tool.GetConfig().RedisConfig
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisConfig.Host + ":" + RedisConfig.Port,
		Password: RedisConfig.Password,
		DB:       RedisConfig.Db,
	})
	//RedisClient = redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "",
	//	DB:       0,
	//})
	// 测试redis是否连接成功
	if result, err := RedisClient.Ping().Result(); err != nil {
		log.Println("ping err :", err)
		return
	} else {
		log.Println("redis连接成功：" + result)
	}
}
