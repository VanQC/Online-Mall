package param

// SmsLoginParam 手机号+验证码登录时的参数
type SmsLoginParam struct {
	Phone string `json:"phone" binding:"required,len=11" msg:"手机号必须为11位"`
	Code  string `json:"code" binding:"required,len=6" msg:"验证码为六位数字"`
}

type RegisterParam struct {
	SmsLoginParam
	Name       string `json:"name" binding:"required,min=3,max=20" msg:"用户名长度最小3位、最大12位"`
	Password   string `json:"password" binding:"required,min=6" msg:"密码长度至少6位"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password" msg:"两次密码不一致"`
}

// CaptchaInfo 验证码信息
type CaptchaInfo struct {
	Id          string `json:"id"`
	VerifyValue string `json:"verifyValue"` // 验证值
}

// NamePwdParam 用户名密码登录时相关参数
type NamePwdParam struct {
	Name     string `json:"name,omitempty" binding:"required,min=3,max=20" msg:"用户名长度最小3位、最大12位"` // 用户名
	Password string `json:"password,omitempty" binding:"required,min=6" msg:"密码长度至少6位"`           // 密码
	CaptchaInfo
}
