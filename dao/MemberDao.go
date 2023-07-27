package dao

import (
	"cloudrestaurant/model"
	"cloudrestaurant/tool"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type MemberDao struct {
}

// InsertCode 向数据库中写入短信验证码数据
func (md MemberDao) InsertCode(sms model.SmsCode) error {
	if err := DB.Create(&sms).Error; err != nil {
		return err
	}
	return nil
}

// ValidateSmsCode 查询登录时验证码是否正确
func (md MemberDao) ValidateSmsCode(phone, code string) *model.SmsCode {
	smsCode := model.SmsCode{}
	if err := DB.Where("phone=? AND code=?", phone, code).First(&smsCode).Error; err != nil {
		log.Println("验证码输入错误" + err.Error())
	}
	// 当验证码验证成功后需要删除该条记录，以回收存储空间，同时防止验证码重复。
	if err := DB.Delete(&smsCode).Error; err != nil {
		log.Println("删除验证码错误" + err.Error())
	}
	return &smsCode
}

// QueryByPhone 根据手机号member表中查询是否包含该用户
func (md MemberDao) QueryByPhone(phone string) *model.Member {
	member := model.Member{Mobile: phone}
	if err := DB.First(&member).Error; err != nil {
		log.Println("查询用户错误" + err.Error())
	}
	return &member
}

// CreateMember 向member表中创建用户
func (md MemberDao) CreateMember(newMember *model.Member) error {
	if err := DB.Create(newMember).Error; err != nil {
		log.Println("添加用户错误" + err.Error())
		return err
	}
	return nil
}

// IsPasswordRight 检查密码是否正确
func (md MemberDao) IsPasswordRight(c *gin.Context, name, inputPassword string) *model.Member {
	var user model.Member
	if err := DB.Where("user_name=?", name).First(&user).Error; err != nil {
		tool.Fail(c, "用户名输入错误", nil)
		return nil
	}

	// 解析数据库中加密的密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)); err != nil {
		tool.Fail(c, "密码输入错误", nil)
		return nil
	}
	return &user
}

// AddAvatar 向用户表中添加avatar字段信息
func (md MemberDao) AddAvatar(id int, path string) error {
	var meb model.Member
	meb.Id = id
	if err := DB.First(&meb).Error; err != nil {
		log.Println("更新头像时，查询用户id出错")
		return errors.New("更新头像时，查询用户id出错" + err.Error())
	}
	if err := DB.Model(&meb).Update("avatar", path).Error; err != nil {
		return errors.New("更新头像时，写入avatar出错" + err.Error())
	}
	return nil
}

// QueryById 根据id在member表中查询是否包含该用户
func (md MemberDao) QueryById(id int) *model.Member {
	member := model.Member{Id: id}
	if err := DB.First(&member).Error; err != nil {
		log.Println("查询用户错误" + err.Error())
		return nil
	}
	return &member
}
