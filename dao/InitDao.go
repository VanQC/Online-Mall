package dao

import (
	"cloudrestaurant/model"
	"cloudrestaurant/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"os"
	"strings"
	"time"
)

var DB *gorm.DB

// GetDB 获取DB对象
func GetDB() *gorm.DB {
	return DB
}

// InitDatabase 初始化数据库连接
func InitDatabase() (err error) {
	DBConfig := tool.GetConfig().Database

	dsn := strings.Join([]string{DBConfig.User, ":", DBConfig.Password, "@tcp(", DBConfig.Host, ":" + DBConfig.Port, ")/",
		DBConfig.DbName, "?charset=", DBConfig.Charset, "&parseTime=", DBConfig.ParseTime, "&loc=", DBConfig.Loc}, "")
	pathRead := dsn
	pathWrite := dsn
	// todo 尝试换为logurs
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: ormLogger})
	if err != nil {
		return err
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(20)  // 设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) // 打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	err = DB.Use(dbresolver.
		Register(dbresolver.Config{
			// `db2` 作为 sources，`db3`、`db4` 作为 replicas
			Sources:  []gorm.Dialector{mysql.Open(pathRead)},                         // 写操作
			Replicas: []gorm.Dialector{mysql.Open(pathWrite), mysql.Open(pathWrite)}, // 读操作
			Policy:   dbresolver.RandomPolicy{},                                      // sources/replicas 负载均衡策略
		}))
	if err != nil {
		return err
	}
	Migration()
	return nil
}

// Migration 数据迁移
func Migration() {
	err := DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&model.Member{},
		&model.Product{},
		&model.Category{},
		&model.Favorite{},
		&model.ProductImg{},
		&model.Order{},
		&model.Address{},
	)
	if err != nil {
		fmt.Println("register table fail")
		os.Exit(0)
	}
	fmt.Println("register table success")
}
