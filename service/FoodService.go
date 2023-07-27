package service

import (
	"cloudrestaurant/dao"
	"cloudrestaurant/model"
)

type FoodService struct {
}

// GetCategories 服务层获取食品信息
func (fs *FoodService) GetCategories() ([]model.FoodCategory, error) {
	fd := dao.FoodDao{}
	return fd.FindFoodCategory()
}

// AddCategories 添加食品信息
func (fs *FoodService) AddCategories(fcg []model.FoodCategory) error {
	fd := dao.FoodDao{}
	return fd.AddFoodCategory(fcg)
}
