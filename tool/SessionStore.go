package tool

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
)

// InitSessionStore 初始化session操作
func InitSessionStore(engine *gin.Engine) {
	config := GetConfig().RedisConfig
	store, err := redis.NewStore(10, "tcp", config.Host+":"+config.Port, "", []byte("secret"))
	if err != nil {
		log.Println(err.Error())
	}
	engine.Use(sessions.Sessions("mysession", store))
}

// SetSession 设置session
func SetSession(ctx *gin.Context, key, value interface{}) error {
	session := sessions.Default(ctx)
	if session == nil {
		return errors.New("get nil session")
	}
	session.Set(key, value)
	return session.Save()
}

// GetSession 获取session
func GetSession(ctx *gin.Context, key interface{}) interface{} {
	session := sessions.Default(ctx)
	return session.Get(key)
}
