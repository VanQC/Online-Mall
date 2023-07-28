package model

import "gorm.io/gorm"

// SmsCode 验证码信息
type SmsCode struct {
	Id    int    `gorm:"primary_Key" json:"id"`
	Phone string `gorm:"varchar(11);not null" json:"phone"`
	BizId string `gorm:"varchar(30)" json:"biz_id"` //发送回执ID。可根据发送回执ID在接口QuerySendDetails中查询具体的发送状态。
	Code  string `gorm:"varchar(6)" json:"code"`
	gorm.Model
}

// Member 用户数据结构
type Member struct {
	Id           uint   `json:"id,omitempty"`
	UserName     string `gorm:"varchar(20)" json:"userName,omitempty"`
	Mobile       string `gorm:"varchar(11)" json:"mobile,omitempty"`
	Password     string `gorm:"varchar(255)" json:"password,omitempty"`
	RegisterTime int64  `gorm:"bigint" json:"register_time,omitempty"`
	Avatar       string `gorm:"varchar(255)" json:"avatar,omitempty"` // 用户头像
	Money        string `gorm:"double" json:"balance,omitempty"`      // 用户余额
	IsActive     string `gorm:"tinyint" json:"is_active,omitempty"`   // 用户是否激活
	City         string `gorm:"varchar(10)" json:"city,omitempty"`
}

// AvatarURL 头像地址
func (mem *Member) AvatarURL() string {
	return mem.Avatar
}
