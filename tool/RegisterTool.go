package tool

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

// RandomName 随机生成用户名
func RandomName(n int) string {
	text := []byte("SFALDSFAOISFDASDDADOQOWEQMXFASLFOAWEWGQEQWEDSMNZJIFJOAI")
	res := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range res {
		res[i] = text[rand.Intn(len(text))]
	}
	return string(res)
}

// CheckPhonePassword 验证手机号和密码是否符合格式要求
// 这个后来被弃用，改用 binding 验证器进行参数格式验证
func CheckPhonePassword(c *gin.Context, phone, password string) bool {
	if len(phone) != 11 {
		log.Println(phone)
		log.Println(len(phone))
		// 415-->服务器无法处理请求附带的媒体格式
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"code": 422,
			"msg":  "手机号必须为11位",
		})
		return false
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return false
	}
	return true
}

// GetValidMsg 该函数会接收一个错误对象和一个任意类型的结构体指针作为参数，然后使用反射来获取结构体中指定字段的 msg 标签的值。
func GetValidMsg(err error, obj any) string {
	// 使用的时候，需要传obj的指针。通过反射获取
	getObj := reflect.TypeOf(obj)
	// 将err接口断言为validator包下的ValidationErrors类型，ok=true则断言成功，否则断言失败。
	if errs, ok := err.(validator.ValidationErrors); ok {
		// 断言成功。errs为切片，可能包含多种错误信息
		for _, e := range errs {
			// 循环每一个错误信息
			// 根据报错字段名，获取结构体的具体字段  【字段，英文：field】
			if field, exist := getObj.Elem().FieldByName(e.Field()); exist {
				msg := field.Tag.Get("msg")
				return msg
			}
		}
	} else {
		log.Println("断言失败")
	}
	return err.Error()
}
