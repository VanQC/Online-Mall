package controller

import (
	"cloudrestaurant/captcha"
	"cloudrestaurant/dao"
	"cloudrestaurant/middlewares"
	"cloudrestaurant/model"
	"cloudrestaurant/param"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MemberController struct {
	service service.MemberService
	dao     dao.MemberDao
}

// MemberRouter 解析"/api/"路由接口
func (mc *MemberController) MemberRouter(engine *gin.Engine) {
	// url: http://localhost:8090/api/
	api := engine.Group("/api/")
	{
		// 发送短信验证码
		api.GET("sendcode", mc.sendSmsCode)
		// 手机号+密码+短信验证码--注册（自己写的）
		api.POST("register", mc.register)
		// 手机号+短信验证码登录or注册
		api.POST("login_sms", mc.smsLogin)
		// 生成图片验证码
		api.GET("captcha", mc.generateCaptcha)
		// 验证图片验证码
		api.POST("verifycaptcha", mc.verifyCaptcha)

		// 手机号+密码登录
		api.POST("login_pwd", mc.pwdLogin)
		// 头像上传
		api.POST("upload/avatar", middlewares.JwtAuth(), mc.uploadAvatar)
		// 查询用户信息
		api.GET("user_info", middlewares.JwtAuth(), mc.userInfo)
	}
}

// SendSmsCode （handlerFunc）发送短信验证码 /sendcode?phone=13003509877
func (mc *MemberController) sendSmsCode(c *gin.Context) {
	// 获取请求参数
	phone, exist := c.GetQuery("phone")
	if !exist {
		tool.Fail(c, "参数解析失败", nil)
		return
	}

	// 调用“用户服务”下的发送验证码的方法
	if mc.service.SendCode(phone) {
		tool.Success(c, "发送成功", nil)
		return
	} else {
		tool.Fail(c, "发送失败", nil)
	}
}

// register 手机号+密码+短信验证码注册
func (mc *MemberController) register(c *gin.Context) {
	// 1.从前端绑定参数
	rp := param.RegisterParam{}
	if err := c.ShouldBindJSON(&rp); err != nil {
		tool.Fail(c, tool.GetValidMsg(err, &rp), nil) // 通过binding实现参数绑定格式验证，并返回自定义错误
		return
	}

	// 2.数据格式验证————通过binding代替
	//if !tool.CheckPhonePassword(c, rp.Phone, rp.Password) {
	//	c.JSON(http.StatusOK, gin.H{"message": "输入数据格式错误"})
	//	return
	//}

	// 3.根据手机号在数据库member表中查是否包含该用户
	user := mc.dao.QueryByPhone(rp.Phone)
	if user.Id != 0 {
		tool.Fail(c, "用户已存在，无法注册", nil)
		return
	}
	// 4.用户注册
	member := mc.service.Register(rp)
	if member != nil {
		tool.Success(c, "注册成功", member)
		return
	}
	tool.Fail(c, "注册失败", nil)
}

// smsLogin 手机短信登录
func (mc *MemberController) smsLogin(c *gin.Context) {
	// 1.前端参数解析绑定
	smsLoginParam := param.SmsLoginParam{}
	if err := c.ShouldBindJSON(&smsLoginParam); err != nil {
		tool.Fail(c, tool.GetValidMsg(err, &smsLoginParam), nil)
		return
	}
	tool.Success(c, "解析参数成功", smsLoginParam)

	// 2.根据解析数据实现登录功能
	member := mc.service.SMSLogin(smsLoginParam)
	if member != nil {
		// 登录成功后发送token
		token, err := tool.GenerateToken(member)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统错误"})
			log.Println("token generate error: " + err.Error())
		}
		// 返回结果
		tool.Success(c, "登录成功", gin.H{"token": token})
	}
	//if member != nil {
	//	// 登录成功后创建session
	//	sess, _ := json.Marshal(member)
	//	err := tool.SetSession(c, "user_"+strconv.Itoa(member.Id), sess)
	//	if err != nil {
	//		tool.Fail(c, "登录失败-创建session失败", nil)
	//		return
	//	}
	//	tool.Success(c, "登录成功", member)
	//	return
	//}
	//tool.Fail(c, "登录失败", nil)
}

// generateCaptcha 生成图片验证码
func (mc *MemberController) generateCaptcha(c *gin.Context) {
	id, b64s, err := captcha.CreateCaptcha()
	if err != nil {
		tool.Fail(c, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "data": b64s, "captchaId": id, "msg": "success"})
}

// verifyCaptcha 验证验证码是否正确
func (mc *MemberController) verifyCaptcha(c *gin.Context) {
	cbj := param.CaptchaInfo{}
	if err := c.ShouldBindJSON(&cbj); err != nil {
		tool.Fail(c, "用户参数解析失败", nil)
		return
	}
	//verify the captcha
	if captcha.CheckCaptcha(cbj.Id, cbj.VerifyValue) {
		tool.Success(c, "验证成功", nil)
	} else {
		tool.Fail(c, "验证失败", nil)
		return
	}
}

// pwdLogin 用户名密码 + 图片验证码登录
func (mc *MemberController) pwdLogin(c *gin.Context) {
	// 1.获取前端的用户名和密码
	npp := param.NamePwdParam{}
	if err := c.ShouldBindJSON(&npp); err != nil {
		tool.Fail(c, tool.GetValidMsg(err, &npp), nil)
		return
	}
	// 2.验证图片验证码
	if !captcha.CheckCaptcha(npp.Id, npp.VerifyValue) {
		tool.Fail(c, "验证码输入错误", nil)
		return
	}
	// 3.验证用户名和密码是否正确
	meb := mc.dao.IsPasswordRight(c, npp.Name, npp.Password)
	if meb != nil {
		// 登录成功后发送token
		token, err := tool.GenerateToken(meb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统错误"})
			log.Println("token generate error: " + err.Error())
		}
		// 返回结果
		tool.Success(c, "登录成功", gin.H{"token": token})
	}
	// 改为 jwt 认证
	//if meb != nil {
	//	// 登录成功后创建session
	//	sessionValue, _ := json.Marshal(meb)
	//	err := tool.SetSession(c, "user_"+strconv.Itoa(meb.Id), sessionValue)
	//	if err != nil {
	//		tool.Fail(c, "登录失败-创建session失败", nil)
	//		return
	//	}
	//	// 设置cookie
	//	c.SetCookie("cookie_user", strconv.Itoa(meb.Id), 10*60, "/", "localhost", false, true)
	//	tool.Success(c, "登录成功", meb)
	//	return
	//}
	//tool.Fail(c, "登录失败", nil)
}

// uploadAvatar 更新用户头像
func (mc *MemberController) uploadAvatar(c *gin.Context) {
	// 1.解析上传的参数：user_id 和 file文件
	//userId := c.PostForm("user_id")
	//log.Println("userid:" + userId)
	// 从请求中读取单个文件
	fileRead, err := c.FormFile("avatar") //对应前端代码：<input type="file" name="avatar">
	if err != nil {
		tool.Fail(c, "参数解析失败", nil)
		return
	}

	// 2.判断userid对应的用户是否登录-->更改为 jwt中间件 验证(v1.0.2)
	//sessionValue := tool.GetSession(c, "user_"+userId)
	//if sessionValue == nil {
	//	tool.Fail(c, "参数不合法，用户尚未登录", nil)
	//	return
	//}
	//// 将session value解析为Member结构体
	//var member model.Member
	//err = json.Unmarshal(sessionValue.([]byte), &member)
	//if err != nil { // 先断言，再转换
	//	log.Println("session 解析失败")
	//	return
	//}
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
	if err = mc.dao.AddAvatar(member.Id, savePath[1:]); err != nil {
		tool.Fail(c, "头像更新失败", err)
		return
	}

	// 5.返回结果
	tool.Success(c, "头像上传成功", nil)
	//给前端返回文件
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintln("attachment; filename=%s", fileRead.Filename))
	c.File(savePath)
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

// userInfo 获取用户信息
//func (mc *MemberController) userInfo(c *gin.Context) {
//	v, exist := c.Get("cookieAuth")
//	if !exist {
//		tool.Fail(c, "查询cookie失败", nil)
//		return
//	}
//	cookie, ok := v.(*http.Cookie)
//	if !ok {
//		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "类型转换失败"})
//		return
//	}
//	id, _ := strconv.Atoi(cookie.Value)
//	if meb := mc.dao.QueryById(id); meb != nil {
//		tool.Success(c, "查询成功", meb)
//		return
//	}
//	tool.Fail(c, "查询失败", nil)
//}
