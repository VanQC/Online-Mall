package middlewares

import (
	"cloudrestaurant/tool"
	"github.com/gin-gonic/gin"
)

const (
	CookieName       = "cookie_user"
	CookieTimeLength = 10 * 60
)

func CookieAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 查找cookie_user是否存在
		// c.Cookie(CookieName)只能返回value，而下面这种可以返回*Cookie
		ck, err := c.Request.Cookie(CookieName)
		if err != nil {
			tool.Fail(c, "cookie验证失败，请登录", nil)
			c.Abort()
			return
		}
		// cookie验证成功后，更新cookie
		c.Set("cookieAuth", ck)
		c.SetCookie(ck.Name, ck.Value, CookieTimeLength, ck.Path, ck.Domain, ck.Secure, ck.HttpOnly)
		c.Next()
	}
}
