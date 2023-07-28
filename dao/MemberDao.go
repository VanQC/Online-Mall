package dao

import (
	"cloudrestaurant/model"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

type MemberDao struct {
	*gorm.DB
}

func NewMemberDao(ctx context.Context) *MemberDao {
	return &MemberDao{DB.WithContext(ctx)}
}

func NewMemberDaoByDB(db *gorm.DB) *MemberDao {
	return &MemberDao{db}
}

// InsertSmsCode 向数据库中写入短信验证码数据
func (dao *MemberDao) InsertSmsCode(sms model.SmsCode) error {
	if err := dao.DB.Create(&sms).Error; err != nil {
		return err
	}
	return nil
}

// ValidateSmsCode 查询登录时验证码是否正确
func (dao *MemberDao) ValidateSmsCode(phone, code string) *model.SmsCode {
	smsCode := model.SmsCode{}
	if err := dao.DB.Where("phone=? AND code=?", phone, code).First(&smsCode).Error; err != nil {
		log.Println("验证码输入错误" + err.Error())
	}
	// 当验证码验证成功后需要删除该条记录，以回收存储空间，同时防止验证码重复。
	if err := dao.DB.Delete(&smsCode).Error; err != nil {
		log.Println("删除验证码错误" + err.Error())
	}
	return &smsCode
}

// IsPasswordRight 检查密码是否正确
func (dao *MemberDao) IsPasswordRight(name, inputPassword string) (*model.Member, error) {
	var user model.Member
	if err := dao.DB.Where("user_name=?", name).First(&user).Error; err != nil {
		return nil, err
	}
	// 解析数据库中加密的密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)); err != nil {
		return nil, fmt.Errorf("密码输入错误")
	}
	return &user, nil
}

// CreateMember 向member表中创建用户
func (dao *MemberDao) CreateMember(newMember *model.Member) error {
	if err := dao.DB.Create(newMember).Error; err != nil {
		log.Println("添加用户错误" + err.Error())
		return err
	}
	return nil
}

// QueryByPhone 根据手机号member表中查询是否包含该用户
func (dao *MemberDao) QueryByPhone(phone string) (*model.Member, bool, error) {
	member := &model.Member{Mobile: phone}
	if err := dao.DB.First(member).Error; err != nil {
		return nil, false, err
	}
	if member.Id != 0 {
		return member, true, nil
	}
	return nil, false, nil
}

// QueryById 根据id在member表中查询是否包含该用户
func (dao *MemberDao) QueryById(id uint) (member *model.Member, err error) {
	err = dao.DB.Model(&model.Member{}).Where("id=?", id).First(&member).Error
	if err != nil {
		log.Println("查询用户错误" + err.Error())
		return
	}
	return
}

func (dao *MemberDao) UpdateUserById(id uint, member *model.Member) (err error) {
	err = dao.DB.Model(&model.Member{}).Where("id=?", id).Updates(&member).Error
	return
}

// AddAvatar 向用户表中添加avatar字段信息
func (dao *MemberDao) AddAvatar(id uint, path string) error {
	var meb model.Member
	meb.Id = id
	if err := dao.DB.First(&meb).Error; err != nil {
		log.Println("更新头像时，查询用户id出错")
		return errors.New("更新头像时，查询用户id出错" + err.Error())
	}
	if err := dao.DB.Model(&meb).Update("avatar", path).Error; err != nil {
		return errors.New("更新头像时，写入avatar出错" + err.Error())
	}
	return nil
}
