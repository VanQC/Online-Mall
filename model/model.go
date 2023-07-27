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
	Id           int    `json:"id,omitempty"`
	UserName     string `gorm:"varchar(20)" json:"userName,omitempty"`
	Mobile       string `gorm:"varchar(11)" json:"mobile,omitempty"`
	Password     string `gorm:"varchar(255)" json:"password,omitempty"`
	RegisterTime int64  `gorm:"bigint" json:"register_time,omitempty"`
	Avatar       string `gorm:"varchar(255)" json:"avatar,omitempty"` // 用户头像
	Balance      string `gorm:"double" json:"balance,omitempty"`      // 用户余额
	IsActive     string `gorm:"tinyint" json:"is_active,omitempty"`   // 用户是否激活
	City         string `gorm:"varchar(10)" json:"city,omitempty"`
}

// FoodCategory 食品种类
type FoodCategory struct {
	// 类别id
	Id int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	// 食品类别标题
	Title string `json:"title" gorm:"varchar(20)"`
	// 食品描述
	Description string `json:"description" gorm:"varchar(30)"`
	// 食品种类图片链接
	ImageUrl string `json:"image_url" gorm:"varchar(255)"`
	// 食品类别链接
	LinkUrl string `json:"link_url" gorm:"varchar(255)"`
	// 该类别是否在服务状态
	IsInServing bool `json:"is_in_serving"`
}
