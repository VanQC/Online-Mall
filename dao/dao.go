package dao

import (
	"cloudrestaurant/model"
	"cloudrestaurant/tool"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// GetDB 获取DB对象
func GetDB() *gorm.DB {
	return DB
}

// InitDatabase 初始化数据库连接
func InitDatabase() (err error) {
	DBConfig := tool.GetConfig().Database

	dsn := DBConfig.User + ":" + DBConfig.Password + "@tcp(" + DBConfig.Host + ":" + DBConfig.Port + ")/" +
		DBConfig.DbName + "?charset=" + DBConfig.Charset + "&parseTime=" + DBConfig.ParseTime + "&loc=" + DBConfig.Loc
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

//func InitDatabase() (err error) {
//	dsn := "root:123456@tcp(127.0.0.1:3306)/cloudrestaurant?charset=utf8mb4&parseTime=True&loc=Local"
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		return
//	}
//	return
//}

// CreateTable 创建数据表，只用执行一次
func CreateTable() {
	// 创建验证码的数据表
	sms := model.SmsCode{}
	if err := DB.First(&sms).Error; err != nil {
		err = DB.AutoMigrate(&sms)
		if err != nil {
			panic("创建数据表错误")
		}
	}

	// 创建用户信息数据表
	member := model.Member{}
	if err := DB.First(&member).Error; err != nil {
		err = DB.AutoMigrate(&member)
		if err != nil {
			panic("创建数据表错误")
		}
	}

	// 创建食品类别数据表
	food := model.FoodCategory{}
	if err := DB.First(&food).Error; err != nil {
		err = DB.AutoMigrate(&food)
		if err != nil {
			log.Panicln("创建数据表错误: " + err.Error())
		}
	}
}
