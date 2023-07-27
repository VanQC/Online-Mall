package middlewares

import (
	"cloudrestaurant/dao"
	"cloudrestaurant/model"
	"cloudrestaurant/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JwtAuth 自定义中间件实现用户认证，解析用户发送的token
// 尝试从用户请求中获取authorization header，若请求头部包含token信息，则进行token的解析，并将解析的结果写入上下文中
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1.获取authorization header（授权头部）
		tokenString := c.GetHeader("Authorization")

		// 2.validate token format 验证token格式
		// token为空，或不是以"Bearer "开，则都是错误的
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort()
			return
		}
		// 前缀占用七个字符，所以从第七位索引开始
		tokenString = tokenString[7:]

		// 3.token验证通过后，解析tokenString
		token, claims, err := tool.ParseToken(tokenString)

		// 解析失败时的操作：
		if err != nil || !token.Valid { // token.Valid 判断token是否有效
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort()
			return
		}

		// 解析成功后的操作：
		userId := claims.UserId
		var member model.Member
		dao.DB.First(&member, userId)
		// 用户信息不存在：
		if member.Id == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort()
			return
		}
		// 用户信息存在：将user信息写入上下文，写入后便可通过c.Get()获取user信息
		c.Set("MemberInfo", member)
		c.Next()
	}
}
