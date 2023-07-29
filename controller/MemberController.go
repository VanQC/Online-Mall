package controller

import (
	"cloudrestaurant/captcha"
	"cloudrestaurant/dao"
	"cloudrestaurant/middlewares"
	"cloudrestaurant/model"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type MemberController struct {
}

func (mc *MemberController) MemberRouter(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/")
	{
		v1.POST("register", mc.register)  // 手机号+密码+短信验证码 注册
		v1.POST("login_sms", mc.smsLogin) // 手机号+短信验证码 注册/登录（若已经注册，则登录；若未注册则注册且登录）
		v1.POST("login_pwd", mc.pwdLogin) // 手机号+密码+图片验证码 登录

		v1.GET("sendcode", mc.sendSmsCode)    // 发送短信验证码
		v1.GET("captcha", mc.generateCaptcha) // 生成图片验证码

		auth := v1.Group("", middlewares.JwtAuth())
		{
			auth.GET("user_info", mc.userInfo)          // 查询用户信息
			auth.POST("upload/avatar", mc.uploadAvatar) // 头像上传
		}
	}
}

// register 手机号+密码+短信验证码 注册/登录
func (mc *MemberController) register(c *gin.Context) {
	// 从前端绑定参数
	ms := service.MemberService{}
	if err := c.ShouldBind(&ms); err == nil {
		res := ms.Register(c.Request.Context())
		tool.StatusOk(c, res)
	} else {
		tool.BadRequest(c, tool.GetValidMsg(err, ms)) // 通过binding实现参数绑定格式验证，并返回自定义错误
	}
}

// SendSmsCode （handlerFunc）发送短信验证码 /sendcode?phone=13003509877
func (mc *MemberController) sendSmsCode(c *gin.Context) {
	// 获取请求参数
	phone, exist := c.GetQuery("phone")
	if !exist {
		tool.BadRequest(c, "参数解析失败")
		return
	}
	// 调用“用户服务”下的发送验证码的方法
	ms := service.MemberService{}
	if ms.SendCode(phone) {
		tool.Success(c, "发送成功", nil)
		return
	} else {
		tool.Fail(c, "发送失败", nil)
	}
}

// smsLogin 手机短信登录
func (mc *MemberController) smsLogin(c *gin.Context) {
	ms := service.MemberService{}
	if err := c.ShouldBind(&ms.SmsLoginParam); err == nil {
		res := ms.SMSLogin(c.Request.Context())
		tool.StatusOk(c, res)
	} else {
		tool.BadRequest(c, tool.GetValidMsg(err, ms)) // 返回binding验证器自定义错误
	}
}

// generateCaptcha 生成图片验证码
func (mc *MemberController) generateCaptcha(c *gin.Context) {
	id, b64s, err := captcha.CreateCaptcha()
	if err != nil {
		tool.Fail(c, err.Error(), nil)
		return
	}
	tool.StatusOk(c, gin.H{"code": 1, "data": b64s, "captchaId": id, "msg": "success"})
}

// pwdLogin 用户名密码 + 图片验证码登录
func (mc *MemberController) pwdLogin(c *gin.Context) {
	nps := service.NamePwdService{}
	if err := c.ShouldBind(&nps); err == nil {
		res := nps.PWDLogin(c.Request.Context())
		tool.StatusOk(c, res)
	} else {
		tool.BadRequest(c, tool.GetValidMsg(err, nps)) // 返回binding验证器自定义错误
	}
}

// userInfo 获取用户信息
func (mc *MemberController) userInfo(c *gin.Context) {
	// 从 jwt中间件 设置的"MemberInfo"字段中获取信息
	v, exist := c.Get("MemberInfo")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "msg": "查找用户信息失败"})
		return
	}

	// 进行类型断言，将v断言为User结构体类型，然后输出其中name和telephone字段即可。
	member, ok := v.(model.Member)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "类型转换失败"})
		return
	}
	tool.Success(c, "查询成功", member)
}

// uploadAvatar 更新用户头像
func (mc *MemberController) uploadAvatar(c *gin.Context) {
	fileRead, err := c.FormFile("avatar") //对应前端代码：<input type="file" name="avatar">
	if err != nil {
		tool.Fail(c, "参数解析失败", nil)
		return
	}
	v, exist := c.Get("MemberInfo")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "msg": "查找用户信息失败"})
		return
	}

	// 进行类型断言，将v断言为User结构体类型，然后输出其中name和telephone字段即可。
	member, ok := v.(model.Member)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "类型转换失败"})
		return
	}

	// 3.将file保存到本地
	// 将读取到的文件保存在本地（服务器端）
	savePath := "./uploadfile/" + "[" + strconv.FormatInt(time.Now().Unix(), 10) + "]--" + fileRead.Filename //设置本地保存路径名
	err = c.SaveUploadedFile(fileRead, savePath)
	if err != nil {
		tool.Fail(c, "头像更新失败", nil)
		return
	}

	// 4.将保存的文件路径存储在 用户表 的avatar字段中（数据库）
	md := dao.NewMemberDao(c.Request.Context())
	if err = md.AddAvatar(member.Id, savePath[1:]); err != nil {
		tool.Fail(c, "头像更新失败", err)
		return
	}

	// 5.返回结果
	tool.Success(c, "头像上传成功", nil)
	//给前端返回文件
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintln("attachment; filename=%s", fileRead.Filename))
	c.File(savePath)
}
