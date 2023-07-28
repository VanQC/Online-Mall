package service

import (
	"cloudrestaurant/captcha"
	"cloudrestaurant/dao"
	"cloudrestaurant/e"
	"cloudrestaurant/model"
	"cloudrestaurant/serializer"
	"cloudrestaurant/tool"
	"context"
	"fmt"
	"github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

// SmsLoginParam 手机号+验证码登录时的参数
type SmsLoginParam struct {
	Phone string `json:"phone" binding:"required,len=11" msg:"手机号必须为11位"`
	Code  string `json:"code" binding:"required,len=6" msg:"验证码为六位数字"`
}

// CaptchaInfo 验证码信息
type CaptchaInfo struct {
	CaptchaId   string `json:"captcha_id"`
	VerifyValue string `json:"verifyValue"` // 验证值
}

type MemberService struct {
	SmsLoginParam
	CaptchaInfo
	Name       string `json:"name" binding:"required,min=3,max=20" msg:"用户名长度最小3位、最大12位"` // 用户名
	Password   string `json:"password" binding:"required,min=6" msg:"密码长度至少6位"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password" msg:"两次密码不一致"` // 密码
}

// SendCode 调用阿里云sdk发送短信验证码
func (ms *MemberService) SendCode(phone string) bool {
	// 1.利用随机数产生六位验证码
	codeNum := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().Unix())).Int31n(1000000))
	code := "{\"code\":" + "'" + codeNum + "'" + "}"

	// 2.调用阿里云sdk，完成验证码的发送
	// 从配置文件获取参数，并实例化client.Config结构体
	cfgSms := tool.GetConfig().Sms
	c := &client.Config{AccessKeyId: &cfgSms.AppKey, AccessKeySecret: &cfgSms.AppSecret, Endpoint: &cfgSms.EndPoint}

	// 根据配置信息，新建一个client对象
	newClient, err := dysmsapi20170525.NewClient(c)
	if err != nil {
		log.Println("newclient error" + err.Error())
		return false
	}

	// 实例化SendSmsRequest结构体
	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &phone,
		TemplateCode:  &cfgSms.TemplateCode,
		SignName:      &cfgSms.SignName,
		TemplateParam: &code,
	}

	// 3.发送短信服务请求，接收返回结果，并判断发送状态
	response, err := newClient.SendSms(request)
	fmt.Println(response)
	if err != nil {
		log.Println("newclient error" + err.Error())
		return false
	}

	// 根据Body.Code字段，判断发送状态
	if *response.Body.Code == "OK" {
		// 向数据库中写入记录
		smsCode := model.SmsCode{Phone: phone, Code: codeNum, BizId: *response.Body.BizId}
		md := dao.MemberDao{}
		if err = md.InsertSmsCode(smsCode); err != nil {
			log.Println("将验证码信息写入数据库时发生错误")
			return false
		}
		return true
	}
	return false
}

// Register 验证短信验证码，并进行用户注册
func (ms *MemberService) Register(ctx context.Context) serializer.Response {
	code := e.SUCCESS
	md := dao.NewMemberDao(ctx)

	// 根据手机号在数据库member表中查是否包含该用户
	_, exist, err := md.QueryByPhone(ms.Phone)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户已存在，无法注册",
		}
	}

	// 如果没有用户名，则随机生成十位字符串
	if len(ms.Name) == 0 {
		ms.Name = tool.RandomName(8)
	}

	// 验证手机号+短信验证码是否正确（去数据库查验证码是否对应）
	sms := md.ValidateSmsCode(ms.Phone, ms.Code)
	if sms.Id == 0 { // 手机号+验证码 验证失败
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "验证码错误",
		}
	}

	// 创建用户，并进行用户密码加密
	// 利用bcrypt包中的GenerateFromPassword方法，将输入的密码进行哈希加密，返回hashPassword
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(ms.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("密码加密错误" + err.Error()) // todo can be remove
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	newMember := model.Member{
		UserName:     ms.Name,
		Mobile:       ms.Phone,
		Password:     string(hashPassword),
		RegisterTime: time.Now().Unix(),
	}
	if err = md.CreateMember(&newMember); err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   newMember,
	}
}

// SMSLogin 具体实现用户手机号+短信验证码 注册/登录
func (ms *MemberService) SMSLogin(ctx context.Context) serializer.Response {

	code := e.SUCCESS
	// 验证手机号+验证码是否正确（去数据库查验证码是否对应）
	md := dao.NewMemberDao(ctx)
	sms := md.ValidateSmsCode(ms.Phone, ms.Code)
	if sms.Id == 0 { // 手机号+验证码 验证失败
		return serializer.Response{
			Msg: "验证码错误",
		}
	}

	// 根据手机号在数据库member表中查是否包含该用户
	member, exist, err := md.QueryByPhone(ms.Phone)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if !exist {
		// 不包含用户记录时，向member表中添加记录
		newMember := model.Member{
			UserName:     ms.Phone,
			Mobile:       ms.Phone,
			RegisterTime: time.Now().Unix(),
		}
		if err = md.CreateMember(&newMember); err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}

	// 发放token，以表示成功登录
	token, err := tool.GenerateToken(member)
	if err != nil {
		logging.Info(err)
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   gin.H{"token": token},
		Msg:    e.GetMsg(code),
	}
}

// PWDLogin 用户名密码 + 图片验证码登录
func (ms *MemberService) PWDLogin(ctx context.Context) serializer.Response {
	code := e.SUCCESS
	// 验证图片验证码
	if !captcha.CheckCaptcha(ms.CaptchaId, ms.VerifyValue) {
		return serializer.Response{Msg: "验证码输入错误"}
	}

	// 验证用户名和密码是否正确
	md := dao.NewMemberDao(ctx)
	meb, err := md.IsPasswordRight(ms.Name, ms.Password)
	if err != nil {
		return serializer.Response{Msg: err.Error()}
	}

	// 发放token，以表示成功登录
	token, err := tool.GenerateToken(meb)
	if err != nil {
		logging.Info(err)
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   gin.H{"token": token},
		Msg:    e.GetMsg(code),
	}
}
