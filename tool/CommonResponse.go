package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	FAIL = iota
	SUCCESS
)

// Response 向前端返回响应的通用函数
func Response(c *gin.Context, httpstatu int, code int, msg string, data interface{}) {
	c.JSON(httpstatu, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// Success 用于成功时，向前端返回数据
func Success(c *gin.Context, msg string, data interface{}) {
	Response(c, http.StatusOK, SUCCESS, msg, data)
}

// Fail 用于失败时，向前端返回数据
func Fail(c *gin.Context, msg string, data interface{}) {
	Response(c, http.StatusBadRequest, FAIL, msg, data)
}

func StatusOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"Msg:": data})
}

func BadRequest(c *gin.Context, data interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{"Error:": data})
}
