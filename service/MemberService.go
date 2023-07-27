package service

import (
	"cloudrestaurant/dao"
	"cloudrestaurant/model"
	"cloudrestaurant/param"
	"cloudrestaurant/tool"
	"fmt"
	"github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

type MemberService struct {
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
		if err = md.InsertCode(smsCode); err != nil {
			log.Println("将验证码信息写入数据库时发生错误")
			return false
		}
		return true
	}
	return false
}

// SMSLogin 具体实现用户手机号+短信验证码登录功能
func (ms *MemberService) SMSLogin(param param.SmsLoginParam) *model.Member {
	// 1.获取用户输入的手机号和验证码
	phone := param.Phone
	code := param.Code

	// 2.验证手机号+验证码是否正确（去数据库查验证码是否对应）
	md := dao.MemberDao{}
	sms := md.ValidateSmsCode(phone, code)
	if sms.Id == 0 {
		return nil // 手机号+验证码 验证失败
	}

	// 3.根据手机号在数据库member表中查是否包含该用户
	member := md.QueryByPhone(phone)
	if member.Id != 0 {
		return member
	} else {
		// 4.不包含用户记录时，向member表中添加记录
		newMember := model.Member{
			UserName:     phone,
			Mobile:       phone,
			RegisterTime: time.Now().Unix(),
		}
		if err := md.CreateMember(&newMember); err != nil {
			return nil
		}
		return &newMember
	}
}

// Register 验证短信验证码，并进行用户注册
func (ms *MemberService) Register(param param.RegisterParam) *model.Member {
	// 1.获取用户输入的手机号和验证码
	phone := param.Phone
	code := param.Code
	name := param.Name
	password := param.Password

	// 如果没有用户名，则随机生成十位字符串
	if len(name) == 0 {
		name = tool.RandomName(8)
	}

	// 2.验证手机号+验证码是否正确（去数据库查验证码是否对应）
	md := dao.MemberDao{}
	sms := md.ValidateSmsCode(phone, code)
	if sms.Id == 0 {
		return nil // 手机号+验证码 验证失败
	}

	// 3.创建用户，并进行用户密码加密
	// 利用bcrypt包中的GenerateFromPassword方法，将输入的密码进行哈希加密，返回hashPassword
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("密码加密错误" + err.Error())
		return nil
	}
	newMember := model.Member{
		UserName:     name,
		Mobile:       phone,
		Password:     string(hashPassword),
		RegisterTime: time.Now().Unix(),
	}
	if err = md.CreateMember(&newMember); err != nil {
		return nil
	}
	return &newMember
}
