package model

import (
	"cloudrestaurant/cache"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

// 商品模型
type Product struct {
	gorm.Model
	Name          string `gorm:"size:255;index"`
	CategoryID    uint   `gorm:"not null"`
	Title         string
	Info          string `gorm:"size:1000"`
	ImgPath       string
	Price         string
	DiscountPrice string
	OnSale        bool `gorm:"default:false"`
	Num           int
	BossID        uint
	BossName      string
	BossAvatar    string
}

// View 获取点击数
func (product *Product) View() uint64 {
	key := fmt.Sprintf("view:product:%s", strconv.Itoa(int(product.ID)))
	countStr, _ := cache.RedisClient.Get(key).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView 商品游览
func (product *Product) AddView() {
	id := strconv.Itoa(int(product.ID))
	// 增加视频点击数
	key := fmt.Sprintf("view:product:%s", id)
	cache.RedisClient.Incr(key)
	// 增加排行点击数
	cache.RedisClient.ZIncrBy("rank", 1, id)
}

type ProductImg struct {
	gorm.Model
	ProductID uint `gorm:"not null"`
	ImgPath   string
}

// Category 种类
type Category struct {
	gorm.Model
	CategoryName string
}

type BasePage struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}
