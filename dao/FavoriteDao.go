package dao

import (
	"cloudrestaurant/model"
	"context"
	"gorm.io/gorm"
)

type FavoritesDao struct {
	*gorm.DB
}

func NewFavoritesDao(ctx context.Context) *FavoritesDao {
	return &FavoritesDao{DB.WithContext(ctx)}
}

func NewFavoritesDaoByDB(db *gorm.DB) *FavoritesDao {
	return &FavoritesDao{db}
}

// ListFavoriteByUserId 通过 user_id 获取收藏夹列表
func (dao *FavoritesDao) ListFavoriteByUserId(uId uint, pageSize, pageNum int) (favorites []*model.Favorite, total int64, err error) {
	// 总数
	// Preload("Member") 方法用于预加载 Favorite 模型关联的 Member 模型的数据
	err = dao.DB.Model(&model.Favorite{}).Preload("Member").Where("member_id=?", uId).Count(&total).Error
	if err != nil {
		return
	}
	// 分页
	err = dao.DB.Model(model.Favorite{}).Preload("Member").Where("member_id=?", uId).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).Find(&favorites).Error
	return
}

// CreateFavorite 创建收藏夹
func (dao *FavoritesDao) CreateFavorite(favorite *model.Favorite) (err error) {
	err = dao.DB.Create(&favorite).Error
	return
}

// FavoriteExistOrNot 判断是否存在
func (dao *FavoritesDao) FavoriteExistOrNot(pId, uId uint) (exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.Favorite{}).Where("product_id=? AND user_id=?", pId, uId).Count(&count).Error
	if count == 0 || err != nil {
		return false, err
	}
	return true, err

}

// DeleteFavoriteById 删除收藏夹
func (dao *FavoritesDao) DeleteFavoriteById(productId, userId uint) error {
	return dao.DB.Where("member_id=? AND product_id=?", userId, productId).Delete(&model.Favorite{}).Error
}
