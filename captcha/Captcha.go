package captcha

import (
	"github.com/mojocn/base64Captcha"
)

// 开辟一个验证码使用的存储空间
//var store = base64Captcha.DefaultMemStore // 使用base64captcha自带的store进行存储（内部实现是通过map和list）
// 使用Redis实现store

// CreateCaptcha 生成验证码
func CreateCaptcha() (id, b64s string, err error) {
	// 将驱动设置为包内置的默认数字验证码驱动
	// 这也可以选择其他类型的验证码驱动，或自己进行配置
	driver := base64Captcha.DefaultDriverDigit

	// 生成数字验证码，同时保存至store
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 生成base64图像验证及id
	id, b64s, err = captcha.Generate()
	return
}

// CheckCaptcha 校验验证码
func CheckCaptcha(id, verifyValue string) bool {
	//排除前端提交的空请求
	if id == "" || verifyValue == "" {
		return false
	}
	//调用验证验证码的方法
	if store.Verify(id, verifyValue, true) {
		return true
	}
	return false
}
