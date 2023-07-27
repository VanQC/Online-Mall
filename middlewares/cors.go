package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Cors 跨域访问，cross origin resource share
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method // 获取请求类型GET/POST.....
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header()

		{ // 花括号里面这部分是多余的
			// 解析前端Request Headers中的其他关键字信息
			var headerKeys []string
			for key := range c.Request.Header {
				headerKeys = append(headerKeys, key)
			}
			// 将headerKeys 进行拼接，并用逗号分隔
			headerStr := strings.Join(headerKeys, ",")
			// 在headStr前面加上"access-control-allow-origin, access-control-allow-headers"前缀
			if headerStr != "" {
				headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
			} else {
				headerStr = "access-control-allow-origin, access-control-allow-headers"
			}
		}

		// 当发生跨域访问时，origin字段不为空，会携带相应值
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")                           // 这个应该和下面一行一样
			c.Header("Access-Control-Allow-Origin", "*")                                        // 设置允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE") // 设置允许访问的方法
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json") // 设置返回格式是json
		}
		// 上面这段，教程里并未过多说明
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		// 继续处理后续请求
		c.Next()
	}
}
