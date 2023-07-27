package dao

import (
	"cloudrestaurant/model"
	"errors"
)

type FoodDao struct {
}

// FindFoodCategory 在数据库中查询全部食品类别信息
func (fd *FoodDao) FindFoodCategory() ([]model.FoodCategory, error) {
	var categories []model.FoodCategory
	if err := DB.Find(&categories).Error; err != nil {
		return nil, errors.New("在数据库中查询食品类别失败：" + err.Error())
	}
	return categories, nil
}

// AddFoodCategory 利用事务机制批量添加食品类别信息
func (fd *FoodDao) AddFoodCategory(categories []model.FoodCategory) error {
	// 开启事务
	tx := DB.Begin()
	// 捕获任何可能发生的panic，并在回滚事务后重新抛出异常
	// 这是一种常用的做法，可以确保即使函数内部出现了异常，也可以保证事务的完整性
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 先检查事务对象是否有错误，如果有则直接返回
	if err := tx.Error; err != nil {
		return err
	}
	// 在事务中执行一系列写入操作
	for _, category := range categories {
		if err := tx.Create(&category).Error; err != nil {
			tx.Rollback() // 遇到错误时进行回滚操作
			return err
		}
	}
	// 提交事务
	return tx.Commit().Error
}
